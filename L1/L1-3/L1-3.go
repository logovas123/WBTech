package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	TotalMutex()
	// TotalAtomic()
}

func TotalMutex() {
	wg := sync.WaitGroup{}
	mx := &sync.Mutex{}
	var total int

	array := [5]int{2, 4, 6, 8, 10}

	for _, v := range array {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			mx.Lock()
			total += v * v
			mx.Unlock()
		}(v)
	}
	wg.Wait()
	fmt.Println(total)
}

func TotalAtomic() {
	wg := sync.WaitGroup{}
	var total int64

	array := [5]int64{2, 4, 6, 8, 10}

	for _, v := range array {
		wg.Add(1)
		go func(v int64) {
			defer wg.Done()
			atomic.AddInt64(&total, int64(v*v))
		}(v)
	}
	wg.Wait()
	fmt.Println(total)
}
