package logger

import (
	"log"
	"os"
)

var (
	info *log.Logger
	erro *log.Logger
	warn *log.Logger
)

func init() {
	info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Lmicroseconds)
	warn = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Lmicroseconds)
	erro = log.New(os.Stderr, "[ERRO] ", log.Ldate|log.Lmicroseconds)
}

func Info(format string, args ...interface{}) {
	info.Printf(format, args...)
}

func Warn(format string, args ...interface{}) {
	warn.Printf(format, args...)
}

func Error(format string, args ...interface{}) {
	erro.Printf(format, args...)
}
