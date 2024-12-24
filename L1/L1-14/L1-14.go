package main

import "fmt"

func main() {
	data := []interface{}{1, "str", false, make(chan interface{}), struct{}{}}
	for i, v := range data {
		fmt.Printf("value %v: ", i+1)
		do(v)
	}
}

func do(i interface{}) {
	switch i.(type) {
	case int:
		fmt.Println("This is int!")
	case string:
		fmt.Println("This is string!")
	case bool:
		fmt.Println("This is bool!")
	case chan interface{}:
		fmt.Println("This is channel!")
	default:
		fmt.Println("I don't know this type :(")
	}
}
