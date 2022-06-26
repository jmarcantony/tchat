package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger(filename string) Logger {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	w := io.MultiWriter(os.Stdout, f)
	l := log.Default()
	l.SetOutput(w)
	return Logger{l}
}
