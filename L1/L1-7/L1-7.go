package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const amount = 5

func main() {
	m := make(map[int]int, amount)
	mx := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := 1; i <= amount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var index int64
			for {
				index = rand.Int63n(int64(amount))

				mx.Lock()
				m[int(index+1)] = i
				fmt.Printf("goroutine %v change value of key %v\n", i, index+1)
				fmt.Println(m)
				mx.Unlock()

				time.Sleep(time.Second)
			}
		}(i)
	}
	wg.Wait()
}
