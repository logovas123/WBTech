package main

/*
Паттерн строитель используется когда сложный продукт состоит из нескольких этапов строительства.
И нам нужно скрыть все этапы строительства от пользователя и выдать только готовый продукт.

В примере есть Director который упрвляет процессом строительства. Пользователь взаимодействет только с ним.
В зависимости от того какой дом просит пользователь, тот и строитель будет назначен на строительство дома.
Каждый из строителей имеет одинаковые этапы строительства, но их реализация разная.
Директор лишь говорит какой этап строительства выполнить.
В результате каждый строитель выполняя одинаковые этапы (с разной реализацией), возвращает готовый дом того типа который
попросил пользователь.
*/

// интерфейс определяет этапы строительства дома
type IBuilder interface {
	setWindowType()
	setDoorType()
	setNumFloor()
	getHouse() House
}

// возрващаем строителя, который построит определённый тип дома
func getBuilder(builderType string) IBuilder {
	switch builderType {
	case "normal":
		return newNormalBuilder()
	case "igloo":
		return newIglooBuilder()
	}

	return nil
}

// -----
// описание строителя который строит обычный дом
type NormalBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func newNormalBuilder() *NormalBuilder {
	return &NormalBuilder{}
}

func (b *NormalBuilder) setWindowType() {
	b.windowType = "Wooden Window"
}

func (b *NormalBuilder) setDoorType() {
	b.doorType = "Wooden Door"
}

func (b *NormalBuilder) setNumFloor() {
	b.floor = 2
}

// возвращаем готовый дом
func (b *NormalBuilder) getHouse() House {
	return House{
		doorType:   b.doorType,
		windowType: b.windowType,
		floor:      b.floor,
	}
}

// -----
// описание строителя который строит иглу
type IglooBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func newIglooBuilder() *IglooBuilder {
	return &IglooBuilder{}
}

func (b *IglooBuilder) setWindowType() {
	b.windowType = "Snow Window"
}

func (b *IglooBuilder) setDoorType() {
	b.doorType = "Snow Door"
}

func (b *IglooBuilder) setNumFloor() {
	b.floor = 1
}

// возвращаем готовый дом
func (b *IglooBuilder) getHouse() House {
	return House{
		doorType:   b.doorType,
		windowType: b.windowType,
		floor:      b.floor,
	}
}

// структура которая представляет собой готовый дом (результат работы)
type House struct {
	windowType string
	doorType   string
	floor      int
}

// -----
// директор управляет процессом
type Director struct {
	builder IBuilder
}

func newDirector() *Director {
	return &Director{}
}

// устанавливаем нового строителя, которым будет управлять директор
func (d *Director) setBuilder(b IBuilder) {
	d.builder = b
}

// строим и возвращаем готовый дом
func (d *Director) buildHouse() House {
	d.builder.setDoorType()
	d.builder.setWindowType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}

/*
func main() {
	// создаём новых строителей
	normalBuilder := getBuilder("normal")
	iglooBuilder := getBuilder("igloo")

	director := newDirector() // создали директора

	director.setBuilder(normalBuilder)   // устанавливаем строителя
	normalHouse := director.buildHouse() // строим и возвращаем готовый дом, пользователь получает готовый дом, не зная о внутренних этапах строительства

	fmt.Printf("Normal House Door Type: %s\n", normalHouse.doorType)
	fmt.Printf("Normal House Window Type: %s\n", normalHouse.windowType)
	fmt.Printf("Normal House Num Floor: %d\n", normalHouse.floor)

	// построили также другой дом, пользователь только сказал название дома и получил его
	director.setBuilder(iglooBuilder)
	iglooHouse := director.buildHouse()

	fmt.Printf("\nIgloo House Door Type: %s\n", iglooHouse.doorType)
	fmt.Printf("Igloo House Window Type: %s\n", iglooHouse.windowType)
	fmt.Printf("Igloo House Num Floor: %d\n", iglooHouse.floor)
}
*/
