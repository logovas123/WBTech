package main

import "fmt"

// родительская структура
type Human struct {
	FirstName string
	LastName  string
	Age       int
}

// методы закреплены за структурой Human
func (h Human) Go() string {
	return "human go"
}

func (h Human) SayHello() string {
	return "human say \"hello\""
}

// дочерняя структура, которая наследуется от структуры Human
// структура Human встроена в структуру Action, это позволяет структуре Action вызывать методы, которые принадлежат структуре Human
type Action struct {
	Human
}

// если у родительской и дочерней структур одининаковые имена, то будет вызван метод дочерней структуры
/* func (h Action) Go() string {
	return "Action: go"
} */

func main() {
	action := Action{
		Human: Human{Age: 39},
	}

	// теперь можно обратиться к полю Age напрямую
	fmt.Println("Age:", action.Age)
	// можно вызывать методы, которые принадлежат структуре Human
	fmt.Println(action.SayHello())
	fmt.Println(action.Go())
}
