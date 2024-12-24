package main

import "fmt"

func main() {
	sl := []string{"cat", "cat", "dog", "cat", "tree"}

	fmt.Println(createSet(sl))
}

func createSet(sl []string) map[string]struct{} {
	resultSet := make(map[string]struct{})

	for _, v := range sl {
		resultSet[v] = struct{}{}
	}

	return resultSet
}
