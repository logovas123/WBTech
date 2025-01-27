package main

import "fmt"

/*
Суть бинарного поиска.
Отсортированный массив делим пополам и находим середину.
Сравниваем элемент с искомым элементом.
Если искомый элемент больше среднего, продолжаем поиск в правой части. Иначе в левой.
Повторяем тот же алгоритм для выбранной части.
*/

func main() {
	arr := []int{1, 2, 4, 5, 6, 13, 435, 786, 989}
	index := binarySearch(arr, 435)
	if index == -1 {
		fmt.Println("element not found")
		return
	}
	fmt.Println("element found, index:", index)
}

func binarySearch(arr []int, elem int) int {
	minIndex := 0
	maxIndex := len(arr) - 1

	for maxIndex >= minIndex {
		midIndex := minIndex + (maxIndex-minIndex)/2
		if arr[midIndex] == elem {
			return midIndex
		} else if arr[midIndex] < elem {
			minIndex = midIndex + 1
		} else {
			maxIndex = midIndex - 1
		}
	}

	// если не нашли элемент то возвращаем -1
	return -1
}
