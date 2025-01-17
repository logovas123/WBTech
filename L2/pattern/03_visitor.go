package main

import (
	"fmt"
	"math"
)

/*
Иногда возникает необходимость добалвения нового функционала в уже существующую структуру. Однако когда мы добавляем новое поведения,
то рискуем "сломать уже существующий код". С такой задачей справляется паттерн Посетитель. Паттерн позволяет добавлять новое поведение в
в структру без изменения самой струткуры.

В примере у нас есть фигуры, к которым необходимо добавить новый метод. Для этого опрделяется новый интерфейс Посетитель. К каждой
фигуре мы добавляем метод accept(), который будет принимать посетителя и выполнять метод, реализация которого определена в структуре,
которая реализует интерфейс посетителя. Остаётся только реализовать струткуру, и методы струткуры, которые добавляют новый функционал.
Заметим, что в уже существующий код мы добавили только метод accept(), поэтому поведение "старого" кода никак не изменилось.
*/

// интерфейс, который опредеялет поведение фигуры
type Shape interface {
	getType() string
	accept(Visitor) // добавили метод accept, через который обращаемся к методам с новым функционалом
}

// существующие фигуры

// -----
type Square struct {
	side int
}

func (s *Square) accept(v Visitor) {
	v.visitForSquare(s)
}

func (s *Square) getType() string {
	return "Square"
}

// -----
type Circle struct {
	radius int
}

func (c *Circle) accept(v Visitor) {
	v.visitForCircle(c)
}

func (c *Circle) getType() string {
	return "Circle"
}

// -----
type Rectangle struct {
	l int
	b int
}

func (t *Rectangle) accept(v Visitor) {
	v.visitForrectangle(t)
}

func (t *Rectangle) getType() string {
	return "Rectangle"
}

// -----
// определяем интерфейс посетителя
type Visitor interface {
	visitForSquare(*Square)
	visitForCircle(*Circle)
	visitForrectangle(*Rectangle)
}

// добавляем структуры которые реализует интерфейс посетителя и определеляют новое поведение

// -----
// поведение для расчёта площади
type AreaCalculator struct {
	area float64
}

func (a *AreaCalculator) visitForSquare(s *Square) {
	a.area = math.Pow(float64(s.side), 2)
	fmt.Println("Calculating area for square. Value:", a.area)
}

func (a *AreaCalculator) visitForCircle(s *Circle) {
	a.area = 2 * math.Pi * math.Pow(float64(s.radius), 2)
	fmt.Println("Calculating area for circle. Value:", a.area)
}

func (a *AreaCalculator) visitForrectangle(s *Rectangle) {
	a.area = float64(0.5 * float64(s.l*s.b))
	fmt.Println("Calculating area for rectangle. Value:", a.area)
}

// -----
// поведение для расчёта периметра
type PerimeterCalculator struct {
	perimetr float64
}

func (p *PerimeterCalculator) visitForSquare(s *Square) {
	p.perimetr = float64(s.side) * 4
	fmt.Println("Calculating perimter for square. Value:", p.perimetr)
}

func (p *PerimeterCalculator) visitForCircle(c *Circle) {
	p.perimetr = 2 * math.Pi * float64(c.radius)
	fmt.Println("Calculating perimter for circle. Value:", p.perimetr)
}

func (p *PerimeterCalculator) visitForrectangle(t *Rectangle) {
	p.perimetr = 10
	fmt.Println("Calculating perimter for rectangle. Value:", p.perimetr)
}

/*
func main() {
	// создаём фигуры
	square := &Square{side: 2}
	circle := &Circle{radius: 3}
	rectangle := &Rectangle{l: 2, b: 3}

	// поведение для расчёта площади
	areaCalculator := &AreaCalculator{}
	square.accept(areaCalculator)
	circle.accept(areaCalculator)
	rectangle.accept(areaCalculator)

	fmt.Println()

	// новое поведение для расчёта периметра, при этом старый код ни как не поменялся
	middleCoordinates := &PerimeterCalculator{}
	square.accept(middleCoordinates)
	circle.accept(middleCoordinates)
	rectangle.accept(middleCoordinates)
}
*/
