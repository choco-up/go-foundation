package database

import (
	"fmt"
	"strings"
)

type OrderBy struct {
	Key string
	Order string
}

func CreateOrderByStr(orders []OrderBy) string {
	var filters []string
	for _, item := range orders {
		filters = append(filters, fmt.Sprintf("%s %s", item.Key, item.Order))
	}
	var qStr = strings.Join(filters[:], fmt.Sprintf(" %s ", ", "))
	return qStr
}
