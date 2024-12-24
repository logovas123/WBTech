package main

import (
	"fmt"
	"math/big"
)

func main() {
	a := big.NewFloat(1 << 24)
	b := big.NewFloat(1 << 21)

	sum := new(big.Float)
	subtr := new(big.Float)
	mult := new(big.Float)
	div := new(big.Float)

	sum.Add(a, b)
	subtr.Sub(a, b)
	mult.Mul(a, b)
	div.Quo(a, b)

	fmt.Printf("a = %g, b = %g\n", a, b)
	fmt.Printf("addition: %g\n", sum)
	fmt.Printf("subtraction: %g\n", subtr)
	fmt.Printf("multiplication: %g\n", mult)
	fmt.Printf("division: %g\n", div)
}
