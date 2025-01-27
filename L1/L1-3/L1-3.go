package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// задачу можно решить через мьютексы или через атомики
	TotalMutex()
	TotalAtomic()
}

func TotalMutex() {
	wg := sync.WaitGroup{} // для ожидания завершения всех горутин
	mx := &sync.Mutex{}    // мьютекс нужен для избежания конкурентного доступа к переменнойц total
	var total int

	array := [5]int{2, 4, 6, 8, 10}

	for _, v := range array {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			mx.Lock() // лочим мьютекс избегая конкурентного доступа к total
			total += v * v
			mx.Unlock()
		}(v)
	}
	wg.Wait() // ждём все горутины
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
			atomic.AddInt64(&total, v*v) // каждая горутина атомарно прибавляет к total квардат числа
		}(v)
	}
	wg.Wait()
	fmt.Println(total)
}
