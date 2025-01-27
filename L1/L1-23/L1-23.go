package main

import "fmt"

func main() {
	sl := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sl = deleteElemFromSlice(sl, 6)
	fmt.Println(sl)
}

// удалить элемент из слайса можно только следующим способом через append
func deleteElemFromSlice(sl []int, ind int) []int {
	return append(sl[:ind], sl[ind+1:]...)
}
