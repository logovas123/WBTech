package main

import (
	"fmt"
	"sort"
)

func main() {
	temperature := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}

	m := make(map[int][]float64) // создаём мапу для гибкого обьединения в группам

	// итерируемся по значениям температуры
	for _, t := range temperature {
		key := int(t) / 10 * 10 // убираем все остатки в числе, оставляя только число кратное 10 (ключ)

		// если ключа не существует, то для ключа создаём массив
		// если существует то аппендим
		if _, ok := m[key]; !ok {
			m[key] = []float64{t}
			continue
		}
		m[key] = append(m[key], t)
	}

	// операции для последовательного вывода ключей
	keys := make([]int, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf("%v:%v ", key, m[key])
	}

	fmt.Println()
}
