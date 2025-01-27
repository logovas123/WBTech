package main

import (
	"fmt"
	"time"
)

// используется паттерн pipe
// суть патерна в том, что программы общаются между собой каналами последовательно
// выходной канал одной программы, явялется входным для другой

func main() {
	c := StartConveyor()
	out := MultiplyBy2(c)

	// конец конвейра. полученные данные выводим на экран
	// так как данные больше никуда не передаются, поэтому это конец конвейра
	for x := range out {
		fmt.Println(x)
	}
	fmt.Println("all numbers print")
}

/*
Начало конвейра. Функция создаёт канал и параллельно в горутине начинает в него писать.
Сам канал возвращается как резулльтат функции, и передаётся в следующую программу.
*/
func StartConveyor() chan int {
	sliceX := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 64, 376, 653, 2, 6, 665465, 34, 43, 4545, 43, 6, 4535, 425, 435}

	out := make(chan int)

	go func() {
		for _, v := range sliceX {
			out <- v
		}
		close(out)
		fmt.Println("all numbers send in channel")
	}()

	return out
}

// программа принимает выходной канал прошлой программы и начинает читать из него
// также создаёт свой выходной канал и начинает писать в него полученные данные и возвращает его, как результат функции
// для слудюущей программы
// это промежуточные программы - их может много
func MultiplyBy2(in chan int) chan int {
	out := make(chan int)
	go func() {
		for x := range in {
			out <- x * 2
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("all numbers multiply by 2")
		close(out)
	}()

	return out
}
