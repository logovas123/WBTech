package main

import "fmt"

// в go множества реализуют мапы
// как значения для ключей можно использовать пустую структуру, так как оно весит 0 байт

func main() {
	a := map[int]struct{}{1: {}, 453: {}, 65: {}, 4: {}, 6: {}, 24: {}, 2: {}, 36: {}, 7: {}, 5: {}, 7643: {}}
	b := map[int]struct{}{345: {}, 98: {}, 89: {}, 1: {}, 453: {}, 47: {}, 765: {}, 4: {}, 6: {}, 2: {}, 435: {}, 36: {}, 7: {}, 49: {}, 37: {}, 7643: {}}
	fmt.Println(IntersectionOfSets(a, b))
}

// обьединение множеств: если ключ одной мапы есть в другой мапе, то это их пересечение - добавляем в результирующую мапу
func IntersectionOfSets(a, b map[int]struct{}) map[int]struct{} {
	resultSet := make(map[int]struct{})

	for ka := range a {
		if _, ok := b[ka]; ok {
			resultSet[ka] = struct{}{}
		}
	}

	return resultSet
}
