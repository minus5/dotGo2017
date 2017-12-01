package logger

import (
	"fmt"
	"io"
	"sync"
)

type Logger struct {
	ch chan []byte
	w  *io.Writer
	wg sync.WaitGroup
}

func New(w io.Writer, cap int) *Logger {
	l := &Logger{
		w:  &w,
		ch: make(chan []byte, cap),
	}
	l.wg.Add(1)
	go func() {
		for p := range l.ch {
			w.Write(p)
		}
		l.wg.Done()
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

func (l *Logger) Close() {
	close(l.ch)
	l.wg.Wait()
}
