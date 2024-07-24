package mysql

import (
	"fmt"

	"gorm.io/gorm"
)

type OrderType string

const (
	OrderType_ASC  OrderType = "ASC"
	OrderType_DESC OrderType = "DESC"
)

type OrderField struct {
	Field     string
	OrderType OrderType
}

func AddDbOrders(db *gorm.DB, orderFields []OrderField) *gorm.DB {
	for _, sortField := range orderFields {
		db = db.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.OrderType))
	}
	return db
}
