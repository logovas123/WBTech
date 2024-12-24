package main

import (
	"fmt"
)

func main() {
	temperature := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}

	m := make(map[int][]float64)

	for _, t := range temperature {
		key := int(float64(int(t)/10) * 10)

		if _, ok := m[key]; !ok {
			m[key] = []float64{t}
			continue
		}
		m[key] = append(m[key], t)
	}

	fmt.Println(m)
}
