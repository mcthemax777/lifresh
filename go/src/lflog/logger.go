package lflog

import (
	"log"
	"os"
)

var goLogger *log.Logger

func init() {
	goLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
}

func Logging(level int, msg string) {
	goLogger.Println(msg)
}

// type Logger struct {	
// }
// func (l *Logger) Logging(level int, msg string) {
// 	goLogger.Println(msg)
// }