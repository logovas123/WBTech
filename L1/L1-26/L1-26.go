package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

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
