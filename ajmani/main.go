package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var testSeconds = 1
var actionNanoseconds = 50

type safeCount struct {
	count int
	sync.Mutex
}

func (sc *safeCount) inc() {
	sc.Lock()
	defer sc.Unlock()
	sc.count++
}

type action struct {
	name  string
	count safeCount
	sync.Mutex
}

func newAction(name string) *action {
	return &action{
		name: name,
	}
}

func (a *action) doWithLock(worker int) {
	a.Lock()
	defer a.Unlock()
	a.do(worker)
}

func (a *action) do(worker int) {
	time.Sleep(time.Duration(actionNanoseconds+rand.Intn(actionNanoseconds)) * time.Nanosecond)
	a.count.inc()
}

type beverage struct {
	name    string
	actions []*action
	count   safeCount
	sync.Mutex
}

func newBeverage() *beverage {
	return &beverage{
		name: "latte",
		actions: []*action{
			newAction("grindBeans"),
			newAction("makeEspresso"),
			newAction("steamMilk"),
			newAction("makeLatte"),
		},
	}
}

func (b *beverage) doBeverageLock(worker int) {
	b.Lock()
	defer b.Unlock()
	for _, a := range b.actions {
		a.do(worker)
	}
	b.count.inc()
}

func (b *beverage) doActionLock(worker int) {
	for _, a := range b.actions {
		a.doWithLock(worker)
	}
	b.count.inc()
}

func useMutexes(beverageLock bool, workers int) {
	b := newBeverage()
	var wg sync.WaitGroup
	wg.Add(workers)
	done := make(chan int)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			for {
				select {
				case <-done:
					wg.Done()
					return
				default:
					if beverageLock {
						b.doBeverageLock(worker)
					} else {
						b.doActionLock(worker)
					}
				}
			}
		}(worker)
	}
	time.Sleep(time.Duration(testSeconds) * time.Second)
	result := b.count.count
	close(done)
	wg.Wait()
	fmt.Println("useMutexes:", "beverageLock:", beverageLock, "workers:", workers, "RESULT:", result)
}

func useChannels(cap, workers int) {
	b := newBeverage()
	var wg sync.WaitGroup
	wg.Add(len(b.actions))
	done := make(chan int)
	in := make(chan int, cap)
	enter := in
	for _, a := range b.actions {
		out := make(chan int, cap)
		go func(in, out chan int, a *action) {
			for {
				select {
				case <-done:
					wg.Done()
					return
				case worker := <-in:
					a.do(worker)
					select {
					case out <- worker:
					case <-done:
					}
				}
			}
		}(in, out, a)
		in = out
	}
	latteReady := in

	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			for {
				select {
				case <-done:
					wg.Done()
					return
				case enter <- worker:
				}
			}
		}(worker)
	}

	wg.Add(1)
	go func() {
		for {
			select {
			case <-done:
				wg.Done()
				return
			case <-latteReady:
				b.count.inc()
			}
		}
	}()

	time.Sleep(time.Duration(testSeconds) * time.Second)
	result := b.count.count
	close(done)
	wg.Wait()
	fmt.Println("useChannels:", "cap:", cap, "workers:", workers, "RESULT:", result)

}

func main() {
	fmt.Println("LOCK BEVERAGE")
	for workers := 1; workers < 10; workers++ {
		useMutexes(true, workers)
	}

	fmt.Println("LOCK ACTION")
	for workers := 1; workers < 10; workers++ {
		useMutexes(false, workers)
	}

	fmt.Println("CHANNELS")
	for capacity := 0; capacity < 10; capacity++ {
		useChannels(capacity, 1)
	}
}
