package storage

import "fmt"

var (
	ErrorOrderExist    = fmt.Errorf("order with this UID exist")
	ErrorOrderNotExist = fmt.Errorf("order not exist")
)
