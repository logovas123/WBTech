package main

import (
	"fmt"
	"time"
)

func main() {
	var n int64
	fmt.Print("Введите время в секундах после которого программа завершится: ")
	fmt.Scan(&n)

	ch := make(chan int)
	timer := time.After(time.Duration(n) * time.Second) // создаём таймер (канал в который придёт значение после промежутка времени введённого пользователем)

	// создаём горутину, которая завершится после того, как среагирует таймер
	go func(ch chan int) {
		var num int
		// в бесконечном цикле отправляем значения в канал
		// в select ожидаем когда "сработает" таймер(по умолчанию отправляем значение в канал), закрываем канал, завершаем функцию
		for {
			select {
			case <-timer:
				fmt.Println("timer triggered")
				close(ch)
				fmt.Println("close main channel")
				return
			default:
				ch <- num // кладём в канал значение
				num++
				time.Sleep(time.Second / 2)
			}
		}
	}(ch)

	// читаем значения из канала пока он не закроется
	for num := range ch {
		fmt.Println(num)
	}
	fmt.Println("program exit") // после закрытия канала программа завершается
}
