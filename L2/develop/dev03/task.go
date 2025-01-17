package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// парсинг флагов командной строки
	column := flag.Int("k", 0, "Номер колонки для сортировки (начиная с 1, 0 - вся строка)")
	numeric := flag.Bool("n", false, "Сортировать по числовому значению")
	reverse := flag.Bool("r", false, "Сортировать в обратном порядке")
	unique := flag.Bool("u", false, "Не выводить повторяющиеся строки")

	flag.Parse()

	// проверка случая если файл не указан
	if flag.NArg() < 1 {
		fmt.Println("Error: file name not specified")
		return
	}

	// открываем файл и копируем из него строки в слайс
	fileName := flag.Arg(0)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error open file: ", err)
		return
	}
	defer file.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error read file: ", err)
		return
	}

	// обрабатываем флаг, вызываем функцию, которая удалит повторяющиеся строки
	if *unique {
		lines = DeleteDublicateLines(lines)
	}

	/*
		Сортировка слайса. Сначала получаем колонку по которой сортируем. Проверяем флаг на сортировку по числам или
		лексиграфически. Проверяем прямой или обратный порядок нужно вернуть.
	*/
	sort.SliceStable(lines, func(i, j int) bool {
		col1 := GetColumn(lines[i], *column)
		col2 := GetColumn(lines[j], *column)

		if *numeric {
			num1, err1 := strconv.Atoi(col1)
			num2, err2 := strconv.Atoi(col2)
			if err1 == nil && err2 == nil {
				if *reverse {
					return num1 > num2
				}
				return num1 < num2
			}
		}

		if *reverse {
			return col1 > col2
		}
		return col1 < col2
	})

	fileResult, err := os.Create("result.txt")
	if err != nil {
		fmt.Println("Error create result file: ", err)
		return
	}

	defer fileResult.Close()

	for _, l := range lines {
		fileResult.WriteString(l + "\n")
	}
}

// добавляем строки в мапу как ключи и получаем только уникальные строки
func DeleteDublicateLines(lines []string) []string {
	linesMap := make(map[string]struct{})
	result := make([]string, 0)

	for _, l := range lines {
		if _, ok := linesMap[l]; !ok {
			result = append(result, l)
		}
	}

	return result
}

// Если колонка это 0, то вернём саму строку, иначе колонку
func GetColumn(line string, col int) string {
	if col == 0 {
		return line
	}
	columns := strings.Split(line, " ")
	if col <= len(columns) {
		return columns[col-1]
	}

	return ""
}
