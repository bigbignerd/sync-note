package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	noteDir       = "/home/bignerd/notes"
	untracked     = "Untracked files"
	modified      = "modified:"
	autoCommit    = "auto commit note"
	operationPush = "push"
	operationPull = "pull"
	syncTick      = 1800 // 30 minutes
)

var (
	GitAdd    []string = []string{"git", "add", "--all"}
	GitPull   []string = []string{"git", "pull", "origin", "master"}
	GitCommit []string = []string{"git", "commit", "-m"}
	GitPush   []string = []string{"git", "push", "-u", "origin", "master"}
	GitStatus []string = []string{"git", "status"}
)

var commitMsg string
var daemon bool
var logger *log.Logger
var env []string

func init() {
	flag.StringVar(&commitMsg, "m", autoCommit, "input commit message")
	flag.BoolVar(&daemon, "d", false, "run note as daemon process")
}

func main() {
	flag.Parse()
	logger = GetLogger()
	env = os.Environ()

	var operation = operationPull
	for _, arg := range os.Args {
		if arg == operationPull || arg == operationPush {
			operation = arg
		}
	}
	if daemon {
		ticker := time.Tick(time.Second * syncTick)
		for {
			<-ticker
			go syncNote(operationPush)
		}
	} else {
		syncNote(operation)
	}
}

//同步笔记
func syncNote(operation string) {

	// 检查文件是否有变化
	fileChange := haveUntrackedFile()
	if fileChange {
		gitCommit(commitMsg)
	}

	//执行push || pull
	switch operation {
	case operationPush:
		gitPull()
		gitPush()
	case operationPull:
		gitPull()
	}
	fmt.Printf("complete!\n")
}

//push to github
func gitPush() bool {
	cmd := NewGitCmd(GitPush)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Printf("push error:%v\n", err)
		logger.Printf("push err msg: %s\n", output)
		return false
	}
	fmt.Println("git push success")
	return true
}

//提交笔记
func gitCommit(msg string) bool {
	cmdAdd := NewGitCmd(GitAdd)
	if err := cmdAdd.Run(); err != nil {
		logger.Printf("git add error:%v\n", err)
		return false
	}

	fmt.Printf("git add success\n")

	gitCommit := append(GitCommit, msg)
	cmdCommit := NewGitCmd(gitCommit)
	if output, err := cmdCommit.CombinedOutput(); err != nil {
		logger.Printf("git commit error:%v\n", err)
		logger.Printf("comm error msg:%s\n", output)
		return false
	}
	fmt.Printf("git commit successful.\n")
	return true
}

//拉取最新的笔记
func gitPull() bool {
	var output []byte
	var err error

	fmt.Printf("git pulling origin master...\n%s", output)

	cmd := NewGitCmd(GitPull)
	if output, err = cmd.CombinedOutput(); err != nil {
		logger.Printf("git pull execute fail. error:%v\n", err)
		logger.Printf("git pull output: %s\n", output)
		return false
	}
	fmt.Println("git pull success")
	return true
}

//是否有未提交的笔记
func haveUntrackedFile() bool {
	var cmd *exec.Cmd
	var output []byte
	var err error

	cmd = NewGitCmd(GitStatus)
	if output, err = cmd.CombinedOutput(); err != nil {
		logger.Printf("git status execute fail. error:%s\n", err)
		logger.Printf("git status output: %s\n", output)
		return false
	}
	if ok := strings.Index(string(output), untracked); ok != -1 {
		return true
	}
	if ok := strings.Index(string(output), modified); ok != -1 {
		return true
	}
	return false
}

//创建一个git command
func NewGitCmd(command []string) *exec.Cmd {
	var cmd *exec.Cmd
	name := command[0]
	args := command[1:]
	cmd = exec.Command(name, args...)
	cmd.Dir = noteDir
	cmd.Env = append(env, "LANG=en_GB")
	return cmd
}
