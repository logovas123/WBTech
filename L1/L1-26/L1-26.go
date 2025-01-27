package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

/*
Приводим символы(слайс рун) к нижнему регистру. Кладём символы в мапу, в итоге остаются только уникальные значения.
И сравниваем количество символов мапы и исходной строки. Если равны, значит все символы в строке уникальные
*/

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		checkSet := map[rune]struct{}{}
		chars := []rune(strings.TrimSpace(scanner.Text()))

		for _, v := range chars {
			checkSet[unicode.ToLower(v)] = struct{}{}
		}

		if len(checkSet) == len(chars) {
			fmt.Println(true)
		} else {
			fmt.Println(false)
		}
	}
}
