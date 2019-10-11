package main

import (
	"io"
	"log"
	"os"
)

const (
	logFile = "/var/log/note-sync.log"
)

func GetLogger() *log.Logger {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(io.Writer(f), "[Note]", log.LstdFlags)
	return logger
}
