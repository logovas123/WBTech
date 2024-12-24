package main

import (
	"fmt"
	"time"
)

func main() {
	c := StartConveyor()
	out := MultiplyBy2(c)

	for x := range out {
		fmt.Println(x)
	}
	fmt.Println("all numbers print")
}

func StartConveyor() chan int {
	sliceX := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 64, 376, 653, 2, 6, 665465, 34, 43, 4545, 43, 6, 4535, 425, 435}

	out := make(chan int)

	go func() {
		for _, v := range sliceX {
			out <- v
		}
		close(out)
	}()

	fmt.Println("all numbers send in channel")

	return out
}

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
