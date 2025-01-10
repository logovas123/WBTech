package storage

import "fmt"

// кастомные ошибки
var (
	ErrorOrderExist    = fmt.Errorf("order with this UID exist")
	ErrorOrderNotExist = fmt.Errorf("order not exist")
)
