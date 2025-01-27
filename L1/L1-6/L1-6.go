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
	var flag int64          // создаём флаг для завершения горутины по флагу
	wg := &sync.WaitGroup{} // нужно для завершения всех горутин (которые будут остановлены разными способами)

	channelForComplete := make(chan struct{}) // создаём канал, который будет использоваться для закрытия горутины через канал

	timer := time.After(10 * time.Second) // канал-таймер для завершения горутины по таймеру

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx) // контекст с закрытием для завершения горутины через контекст

	signalChan := make(chan os.Signal, 1) // создаём канал который будет ожидать сигнал прерывания (сигнал для завершения всех горутин)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		wg.Wait() // ждём завершения всех горутин
		fmt.Println("all goroutines complete success")
		signal.Stop(signalChan)             // канал больше не ждёт сигналы
		fmt.Println("signal channel close") // конец программы
	}()

	// в отдельной горутине ждём сигнал прерывания
	go func() {
		select {
		case <-signalChan:
			cancel()                         // закрываем контекст
			channelForComplete <- struct{}{} // отправляем значение в канал
			close(channelForComplete)        // закрываем канал
			atomic.StoreInt64(&flag, 1)      // атомарно изменяем значение flag по указателю
		}
	}()

	// запускаем горутины
	wg.Add(4)
	go goroutineCompleteByChannel(channelForComplete, wg)
	go goroutineCompleteByContext(ctx, wg)
	go goroutineCompleteByTimer(timer, wg)
	go goroutineCompleteByFlag(&flag, wg)
}

// когда ch получит значение, горутина завершится
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

// когда придёт сигнал о закрытии контекста горутина завершится
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

// эта горутина завершается сама без сигнала прерывания, когда придёт сигнал с таймера
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

// данная горутина в бесконечном цикле атомарно проверяет значение флага, и завершается, когда оно равно единице
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
