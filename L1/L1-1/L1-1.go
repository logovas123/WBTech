package main

import "fmt"

type Human struct {
	FirstName string
	LastName  string
	Age       int
}

func (h Human) Go() string {
	return "human go"
}

func (h Human) SayHello() string {
	return "human say \"hello\""
}

type Action struct {
	Human
}

/* func (h Action) Go() string {
	return "Action: go"
} */

func main() {
	action := Action{
		Human: Human{Age: 39},
	}

	fmt.Println("Age:", action.Age)
	fmt.Println(action.SayHello())
	fmt.Println(action.Go())
}
