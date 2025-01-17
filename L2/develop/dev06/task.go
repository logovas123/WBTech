package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// конфиг
type config struct {
	fields    string
	delimiter string
	separated bool
}

func parseArgs() config {
	var cfg config

	flag.StringVar(&cfg.fields, "f", "", "Выбрать поля (колонки)")
	flag.StringVar(&cfg.delimiter, "d", "\t", "Использовать другой разделитель (по умолчанию TAB)")
	flag.BoolVar(&cfg.separated, "s", false, "Только строки с разделителем")
	flag.Parse()

	return cfg
}

func main() {
	cfg := parseArgs()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		// если флаг -s установлен, проверяем наличие разделителя, если нет разделителя то пропускаем строку
		if cfg.separated && !strings.Contains(line, cfg.delimiter) {
			continue
		}

		// разбиваем строку по разделителю
		fields := strings.Split(line, cfg.delimiter)

		// парсим и выбираем поля
		if cfg.fields != "" {
			fieldIndexes := parseFields(cfg.fields)

			// выводим только нужные поля
			for _, index := range fieldIndexes {
				if index >= 0 && index < len(fields) {
					fmt.Print(fields[index])
					if index != fieldIndexes[len(fieldIndexes)-1] {
						fmt.Print(cfg.delimiter)
					}
				}
			}
		} else {
			// если не указан флаг -f, выводим всю строку
			fmt.Print(line)
		}

		fmt.Println()
	}

	// обработка ошибок
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scan:", err)
	}
}

// функция для парсинга полей, которые нужно выбрать
func parseFields(NumsFields string) []int {
	fieldIndexes := make([]int, 0)
	for _, field := range strings.Split(NumsFields, ",") {
		var index int
		index, err := strconv.Atoi(field)
		if err == nil && index > 0 {
			fieldIndexes = append(fieldIndexes, index-1) // индексация с 0
		}
	}
	return fieldIndexes
}
