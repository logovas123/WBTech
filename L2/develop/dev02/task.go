package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		extractString := scanner.Text()
		unpackString := getUnpackString(extractString)
		fmt.Printf("\"%s\" => \"%s\"\n", extractString, unpackString)
	}
}

/*
Решал задачу по следующему принципу:
проходился по каждому символу строки и сравнивал текущий и предыдущий символ. И в зависимости от того какой порядок
числа и буквы сейчас идёт производил нужные действия.
Переменная num хранит число строкой. К этой строке конкатенируются цифры, которые идут друг за другом в строке, так как
число может быть двузначными и более. В переменной letter храним букву, которая потом num раз будет добалвена к итоговой строке.
*/
func getUnpackString(s string) string {
	var (
		result string
		prev   rune
		letter rune
		num    string
		n      int
	)

	if s == "" || unicode.IsDigit([]rune(s)[0]) {
		return ""
	}

	if len(s) == 1 {
		return s
	}

	for i, r := range s {
		switch {
		case i == 0 && unicode.IsLetter(r):
			prev = r
			letter = r
		case unicode.IsLetter(r) && unicode.IsLetter(prev):
			result += string(prev)
			prev = r
			letter = r
		case unicode.IsDigit(r) && unicode.IsLetter(prev):
			num += string(r)
			prev = r
		case unicode.IsLetter(r) && unicode.IsDigit(prev):
			n, _ = strconv.Atoi(num)
			num = ""
			result += strings.Repeat(string(letter), n)
			prev = r
			letter = r
		case unicode.IsDigit(r) && unicode.IsDigit(prev):
			num += string(r)
			prev = r
		}
	}

	// этот switch нужен чтобы обработать последний символ
	switch {
	case unicode.IsLetter(prev):
		result += string(prev)
	case unicode.IsDigit(prev):
		n, _ = strconv.Atoi(num)
		num = ""
		result += strings.Repeat(string(letter), n)
	}

	return result
}
