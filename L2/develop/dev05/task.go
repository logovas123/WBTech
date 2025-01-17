package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// конфиг
type config struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	pattern    string
	file       string
}

// функция для парсинга командной строки
func parseArgs() config {
	var cfg config

	flag.IntVar(&cfg.after, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&cfg.before, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&cfg.context, "C", 0, "(A+B) печатать ±N строк вокруг совпадения")
	flag.BoolVar(&cfg.count, "c", false, "количество строк")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&cfg.invert, "v", false, "вместо совпадения, исключать")
	flag.BoolVar(&cfg.fixed, "F", false, "точное совпадение со строкой, не паттерн")
	flag.BoolVar(&cfg.lineNum, "n", false, "напечатать номер строки")
	flag.Parse()

	// обязательно наличие минимум двух аргументов: паттерна и названия файла
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Error: args < 2")
		return config{}
	}

	cfg.pattern = args[0]
	cfg.file = args[1]

	return cfg
}

// функция для чтения файла
func readLines(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// функция для сравнения двух строк
func matches(line string, cfg config) bool {
	var pattern string
	if cfg.ignoreCase {
		line = strings.ToLower(line)
		pattern = strings.ToLower(cfg.pattern)
	} else {
		pattern = cfg.pattern
	}

	// если нужно полное совпадение строки, то возвращаем полное сравнение, иначе ищем наличие подстроки
	if cfg.fixed {
		return line == pattern
	}

	return strings.Contains(line, pattern)
}

// функция фильтрации
func filterLines(lines []string, cfg config) {
	var matchedLines []int
	for i, line := range lines {
		match := matches(line, cfg) // функция для сравнения двух строк (в зависимости от параметров)
		if cfg.invert {
			match = !match
		}
		if match {
			matchedLines = append(matchedLines, i)
		}
	}

	// выводим только количество строк, если параметр активен
	if cfg.count {
		fmt.Println(len(matchedLines))
		return
	}

	// определяем границы, которые надо вывести
	// заметим
	var start, end int

	// добавляем в карту индексы для вывода, исключая дубликаты
	withdraw := make(map[int]struct{})
	for _, l := range matchedLines {
		switch {
		case cfg.context != 0:
			start = l - cfg.context
			if start < 0 {
				start = 0
			}
			end = l + cfg.context
			if end >= len(lines) {
				end = len(lines) - 1
			}
		default:
			start = l - cfg.before
			if start < 0 {
				start = 0
			}
			end = l + cfg.after
			if end >= len(lines) {
				end = len(lines) - 1
			}
		}

		for j := start; j <= end; j++ {
			withdraw[j] = struct{}{}
		}
	}

	// упорядочили индексы и вывели строки, по их индексам
	// при наличии нужного параметра печатаем номера строк
	indexes := make([]int, 0, len(withdraw))
	for i := range withdraw {
		indexes = append(indexes, i)
	}

	sort.Ints(indexes)

	for _, i := range indexes {
		line := lines[i]
		if cfg.lineNum {
			fmt.Printf("%v:", i+1)
		}
		fmt.Println(line)
	}
}

func main() {
	cfg := parseArgs() // парсим аргументы командной строки в конфиг

	lines, err := readLines(cfg.file) // функция для чтения файла
	if err != nil {
		fmt.Println("Error read file: ", err)
		return
	}

	filterLines(lines, cfg) // функция фильтрации строк
}
