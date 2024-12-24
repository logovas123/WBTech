package main

import (
	"fmt"

	"L1-24/point"
)

func main() {
	a, b := point.NewPoint(3, 4), point.NewPoint(11, 20)

	dist := a.GetDistance(b)
	fmt.Println(dist)
}
