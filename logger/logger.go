package logger

import (
	"io"
	"log"
	"os"
)

var output = os.Stdout

func init() {
	log.SetOutput(output)
}

//GetOutput returns default logger output
func GetOutput() io.Writer {
	return output
}

//Error log fatal error
func Error(v ...interface{}) {
	log.Panic(v)
}

//Errorf log fatal error (formatted)
func Errorf(s string, v ...interface{}) {
	log.Panicf(s, v)
}

//Critical log critical error
func Critical(v ...interface{}) {
	log.Fatal(v)
}

//Criticalf log critical error (formatted)
func Criticalf(s string, v ...interface{}) {
	log.Fatalf(s, v)
}

//Info logs info message
func Info(v ...interface{}) {
	log.Print(v)
}

//Infof logs info message (formatted)
func Infof(s string, v ...interface{}) {
	log.Printf(s, v)
}
