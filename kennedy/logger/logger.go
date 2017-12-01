package logger

import (
	"fmt"
	"io"
)

type Logger struct {
	ch chan []byte
	w  *io.Writer
}

func New(w io.Writer, cap int) *Logger {
	l := &Logger{
		w:  &w,
		ch: make(chan []byte, cap),
	}
	go func() {
		for p := range l.ch {
			w.Write(p)
		}
	}()
	return l
}

func (l *Logger) Write(p []byte) (length int, err error) {
	select {
	case l.ch <- p:
	default:
		fmt.Println("DROP!")
		return 0, fmt.Errorf("Write package droppped")
	}
	return len(p), nil
}
