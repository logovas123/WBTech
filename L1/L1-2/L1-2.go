package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	array := [5]int{2, 4, 6, 8, 10}

	for _, v := range array {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			fmt.Println(v * v)
		}(v)
	}
	wg.Wait()
}
