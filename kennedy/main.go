package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mprcela/dotGo2017/kennedy/logger"
)

type device struct {
	problem bool
}

func (d *device) Write(p []byte) (l int, err error) {
	if d.problem {
		time.Sleep(3 * time.Second)
	}
	return os.Stdout.Write(p)
}

func main() {
	d := &device{}
	// w := d									// logger that blocks
	w := logger.New(d, 10) // logger that drops lines

	for i := 0; i < 10; i++ {
		go func(id int) {
			for {
				out := fmt.Sprintf("%d: log data\n", id)
				w.Write([]byte(out))
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	r := bufio.NewReader(os.Stdin)
	for {
		r.ReadByte()
		fmt.Println("SWITCHED!")
		d.problem = !d.problem
	}
}
