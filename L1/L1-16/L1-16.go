package main

import "fmt"

func main() {
	array := []int{13, 45, 5, 43, 5, 4, 5, 6456, 4, 4765765, 8, 757, 7657, 67, 545, 6465, 8, 895, 3, 5, 76584, 43, 2, 342, 565, 598, 63, 4321}
	fmt.Println(quicksort(array))
}

func quicksort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	pivot := arr[0]

	left := make([]int, 0)
	right := make([]int, 0)

	for i := 1; i < len(arr); i++ {
		if pivot > arr[i] {
			left = append(left, arr[i])
		} else {
			right = append(right, arr[i])
		}
	}

	result := append(quicksort(left), pivot)
	result = append(result, quicksort(right)...)

	return result
}
