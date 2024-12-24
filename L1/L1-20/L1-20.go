package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "snow dog sun"

	str = strings.TrimSpace(str)
	newStr := GetNewStr(str)

	fmt.Printf("%s - %s\n", str, newStr)
}

func GetNewStr(s string) string {
	sliceStr := strings.Split(s, " ")

	result := make([]string, len(sliceStr))
	for i := len(sliceStr) - 1; i >= 0; i-- {
		result[len(sliceStr)-1-i] = sliceStr[i]
	}

	return strings.Join(result, " ")
}
