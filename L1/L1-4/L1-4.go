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

var flagForComplete bool

func main() {
	wg := &sync.WaitGroup{} // для ожидания завершения всех воркеров

	// создаём контекст с закрытием (нужен для синхронизации между горутиной ждущей сигнал interrupt и горутиной main)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)   // создаём канал, который будет ждать сигнал
	signal.Notify(signalChan, os.Interrupt) // заставляем канал ожидать только сигнала Interrupt

	ch := make(chan interface{}, 50) // по каналу передаётся тип interface, тем самым мы можем передавать произвольные данные

	defer func() {
		wg.Wait()                           // ждём заверешения всех горутин
		signal.Stop(signalChan)             // прекращаем передачу входящих сигналов
		fmt.Println("signal channel close") // после этого вывода программа полностью и корректно завершается
	}()

	// в отдельной горутине ждём сигнал(по нажатию ctrl+c), затем закрываем контекст
	go func() {
		select {
		case <-signalChan:
			cancel()
		}
	}()

	var n int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&n)

	// запускаем воркеры
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go worker(ch, i, wg)
	}

	// массив с произвольными данными
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

	// в бесконечном цикле через select ждём сигнала о закрытии контекста
	// по умолчанию будем класть в канал случайное значение
	for {
		select {
		case <-ctx.Done(): // при получении сигнала о закрытии контекста мы закрываем канал, тем самым завершая работу всех воркеров и через return завершаем работу функции функции main (переходим к блоку defer)
			close(ch)
			fmt.Println("main channel close")
			return
		default:
			randomIndex := rand.Int63n(int64(len(slAny))) // генерируем случайный индекс
			ch <- slAny[randomIndex]                      // кладём в канал значение
			time.Sleep(1 * time.Second)                   // спим
		}
	}
}

// каждый воркер будет читать из канала, пока канала не закроется, при закрытии канала воркер завершает свою работу
func worker(ch chan interface{}, i int, wg *sync.WaitGroup) {
	defer wg.Done()

	for anything := range ch {
		fmt.Printf("worker %v, value %v\n", i, anything)
	}
	fmt.Printf("worker %v complete\n", i)
}
