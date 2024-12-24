package main

import "fmt"

func main() {
	str := "главрыба"

	fmt.Println(GetInvertedString(str))
}

func GetInvertedString(s string) string {
	r := []rune(s)
	result := make([]rune, len(r))

	for i := len(r) - 1; i >= 0; i-- {
		result[len(r)-1-i] = r[i]
	}

	return string(result)
}
