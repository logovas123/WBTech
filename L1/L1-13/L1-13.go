package main

import "fmt"

// язык позволяет менять местами значения переменных, без создания временной переменной

func main() {
	a, b := 1, 2
	fmt.Printf("before: a=%v, b=%v\n", a, b)
	a, b = b, a
	fmt.Printf("after: a=%v, b=%v\n", a, b)
}
