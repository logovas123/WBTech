package main

import (
	"fmt"

	"L1-24/point"
)

/*
Инкапсуляция в го реализуется на уровне пакета
Если название переменной(поля структуры) начинается с большой буквы, то она экспортируемая, если с маленькой то неэкспортируемая
Также работает и с методами структуры.
Доступ к полям структуры осуществляется через методы.
*/
func main() {
	a, b := point.NewPoint(3, 4), point.NewPoint(11, 20)

	dist := a.GetDistance(b)
	fmt.Println(dist)
}
