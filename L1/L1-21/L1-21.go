package main

import "fmt"

// паттерн адаптер позволяет связать несовместимые объекты

// интерфейс, который ожидает клиент, Adapter реализует этот интерфейс
type Target interface {
	Request()
}

// тип который нужно адаптировать
type Adaptee struct{}

// создаём новый адаптер
func NewAdapter(adaptee *Adaptee) Target {
	return &Adapter{adaptee}
}

// метод, к которому у нас нет доступа
func (a *Adaptee) SpecificRequest() {
	fmt.Println("request done")
}

// адаптер для Adaptee
type Adapter struct {
	*Adaptee
}

// метод, который позволяет вызвать недоступный метод
func (a *Adapter) Request() {
	a.SpecificRequest()
}

func main() {
	adaptee := &Adaptee{}

	target := NewAdapter(adaptee)

	target.Request() // должен вывестись "request done"
	fmt.Println("success adapt")
}
