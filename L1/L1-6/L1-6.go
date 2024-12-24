package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var flag int64
	wg := &sync.WaitGroup{}

	channelForComplete := make(chan struct{})

	timer := time.After(10 * time.Second)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		wg.Wait()
		fmt.Println("all goroutines complete success")
		signal.Stop(signalChan)
		fmt.Println("signal channel close")
	}()

	go func() {
		select {
		case <-signalChan:
			cancel()
			channelForComplete <- struct{}{}
			close(channelForComplete)
			atomic.StoreInt64(&flag, 1)
		}
	}()

	wg.Add(4)
	go goroutineCompleteByChannel(channelForComplete, wg)
	go goroutineCompleteByContext(ctx, wg)
	go goroutineCompleteByTimer(timer, wg)
	go goroutineCompleteByFlag(&flag, wg)
}

func goroutineCompleteByChannel(ch chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ch:
			fmt.Println("ByChannel complete")
			return
		default:
			fmt.Println("ByChannel work")
			time.Sleep(500 * time.Millisecond) // some work
		}
	}
}

func goroutineCompleteByContext(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ByContext complete")
			return
		default:
			fmt.Println("ByContext work")
			time.Sleep(500 * time.Millisecond) // some work
		}
	}
}

func goroutineCompleteByTimer(timer <-chan time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-timer:
			fmt.Println("ByTimer was completed 10 seconds after start program")
			return
		default:
			fmt.Println("ByTimer work")
			time.Sleep(500 * time.Millisecond) // some work
		}
	}
}

func goroutineCompleteByFlag(flag *int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if atomic.LoadInt64(flag) == 1 {
			fmt.Println("ByFlag complete")
			return
		}
		fmt.Println("ByFlag work")
		time.Sleep(500 * time.Millisecond) // some work
	}
}
