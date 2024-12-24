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
	timer := time.After(time.Duration(n) * time.Second)

	go func(ch chan int) {
		var num int
		for {
			select {
			case <-timer:
				fmt.Println("timer triggered")
				close(ch)
				fmt.Println("close main channel")
				return
			default:
				ch <- num
				num++
			}
		}
	}(ch)

	for num := range ch {
		fmt.Println(num)
	}
	fmt.Println("program exit")
}
