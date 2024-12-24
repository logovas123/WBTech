package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
)

var flagForComplete bool

func main() {
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	ch := make(chan interface{}, 50)

	defer func() {
		wg.Wait()
		signal.Stop(signalChan)
		fmt.Println("signal channel close")
	}()

	go func() {
		select {
		case <-signalChan:
			cancel()
		}
	}()

	var n int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&n)

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go worker(ch, i, wg)
	}

	slAny := []interface{}{
		1,
		"string",
		true,
		struct {
			num int
		}{
			num: 2,
		},
	}

	for {
		select {
		case <-ctx.Done():
			close(ch)
			fmt.Println("main channel close")
			return
		default:
			randomIndex := rand.Int63n(int64(len(slAny)))
			ch <- slAny[randomIndex]
		}
	}
}

func worker(ch chan interface{}, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	for anything := range ch {
		fmt.Printf("worker %v, value %v\n", i, anything)
	}
	fmt.Printf("worker %v complete\n", i)
}
