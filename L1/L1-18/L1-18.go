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

type counterStruct struct {
	counter int64
}

func newCounterStruct() *counterStruct {
	return &counterStruct{}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	defer signal.Stop(sigChannel)

	// graceful shutdown
	go func() {
		<-sigChannel
		fmt.Println("get signal of complete")
		cancel()
	}()

	wg := &sync.WaitGroup{}

	structCount := newCounterStruct() // создали структуру счётчик

	// запускаем пачку горутин
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go incrementCounter(ctx, wg, structCount)
	}

	wg.Wait()
	fmt.Println("all goroutines complete")
	fmt.Println("counter =", structCount.counter)
}

// в бесконечном цикле инкрементируем счётик, лучше всего это делать с помощью пакета атомик
// слишком жирно выделять на инкремент целый мьютекс
func incrementCounter(ctx context.Context, wg *sync.WaitGroup, c *counterStruct) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			atomic.AddInt64(&c.counter, 1)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
