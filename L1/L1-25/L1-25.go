package main

import (
	"fmt"
	"time"
)

func main() {
	duration := 5
	fmt.Printf("program sleep %vs...\n", duration)
	sleep(time.Duration(duration) * time.Second)
	fmt.Printf("program complete after %vs\n", duration)
}

/*
Функцию sleep можно реализовать с помощью функции time.After(), которая возвращает канал, который вернёт значение
через определённый промежуток времени.
*/
func sleep(t time.Duration) {
	<-time.After(t) // блокиремся(спим), пока не придёт значение в канал
}
