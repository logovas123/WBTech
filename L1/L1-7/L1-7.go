package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

// в мапу нельзя конкурентно записывать так как под капотом все время идёт перераспределение памяти

const amount = 5

func main() {
	m := make(map[int]int, amount)
	mx := &sync.Mutex{}     // мьютекс для конкурентоной записи
	wg := &sync.WaitGroup{} // для запуска пачки горутин
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		wg.Wait()
		signal.Stop(signalChan)
		fmt.Println("signal channel close")
		fmt.Println("Program complete!")
	}()

	// graceful shutdown
	go func() {
		<-signalChan
		cancel()
	}()

	// запускаем пачку горуттин
	for i := 1; i <= amount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var index int64

			// в бесконечном цикле либо ждём заершения контекста,
			// либо конкурентно записываем в мапу
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("goroutine %v complete\n", i)
					return
				default:
					index = rand.Int63n(int64(amount)) // запись будет в случайный индекс

					// чтобы конкурентно записываить в мапу нужно использовать мьютексы
					mx.Lock()
					m[int(index+1)] = i
					fmt.Printf("goroutine %v change value of key %v\n", i, index+1)
					fmt.Println(m)
					mx.Unlock()

					time.Sleep(time.Second)
				}
			}
		}(i)
	}
}
