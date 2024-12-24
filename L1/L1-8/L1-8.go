package main

import (
	"fmt"
)

func main() {
	var num, i, value int64

	fmt.Print("Введите число: ")
	fmt.Scan(&num)
	fmt.Printf("Код числа: %b\n", num)

	fmt.Print("Введите i-й бит: ")
	fmt.Scan(&i)

	fmt.Print("Введите на какое значение хотите изменить i-й бит (1 или 0): ")
	fmt.Scan(&value)

	if value == 1 {
		num = num | (1 << i)
	} else {
		num = num &^ (1 << i)
	}

	fmt.Printf("Ответ: %v, код: %b\n", num, num)
}
