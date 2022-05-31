package database

import (
	"fmt"
	"strconv"
	"strings"
)

type Filter interface {
	ConstructFilter() string
}

func ConstructFilterStr(subFilters []Filter) string {
	var filters []string
	for _, item := range subFilters {
		filters = append(filters, item.ConstructFilter())
	}
	var qStr = strings.Join(filters[:], fmt.Sprintf(" %s ", "AND"))
	return qStr
}

type StringFilter struct {
	Key string
	Inclusive bool
	Values []string
}

func (sq StringFilter) ConstructFilter() string {
	var op string
	if sq.Inclusive {
		op = "OR"
	} else {
		op = "AND"
	}

	var queries []string
	for _, item := range sq.Values {
		queries = append(queries, fmt.Sprintf("%s = '%s'", sq.Key, item))
	}
	var qStr = strings.Join(queries[:], fmt.Sprintf(" %s ", op))
	if len(qStr) == 0 {
		return qStr
	}
	return "( " + qStr + " )"
}

type FloatFilter struct {
	Key string
	Inclusive bool
	Values []float64
}

func (sq FloatFilter) ConstructFilter() string {
	var op string
	if sq.Inclusive {
		op = "OR"
	} else {
		op = "AND"
	}

	var queries []string
	for _, item := range sq.Values {
		queries = append(queries, fmt.Sprintf("%s = %f", sq.Key, item))
	}
	var qStr = strings.Join(queries[:], fmt.Sprintf(" %s ", op))
	if len(qStr) == 0 {
		return qStr
	}
	return "( " + qStr + " )"
}

type IntFilter struct {
	Key string
	Inclusive bool
	Values []int64
}

func (sq IntFilter) ConstructFilter() string {
	var op string
	if sq.Inclusive {
		op = "OR"
	} else {
		op = "AND"
	}

	var queries []string
	for _, item := range sq.Values {
		queries = append(queries, fmt.Sprintf("%s = %d", sq.Key, item))
	}
	var qStr = strings.Join(queries[:], fmt.Sprintf(" %s ", op))
	if len(qStr) == 0 {
		return qStr
	}
	return "( " + qStr + " )"
}

func CreateStringFilter(key string, queries []string) *StringFilter {
	fmt.Println("Key: " + key)
	fmt.Printf("%v\n", queries)
	var filter *StringFilter
	if len(queries) > 0 {
		filter = &StringFilter{
			Key:       key,
			Inclusive: false,
			Values:    nil,
		}
		if len(queries[0]) > 3 && queries[0][0:3] == "in:" {
			queries[0] = queries[0][3:]
			filter.Inclusive = true
		}
		filter.Values = queries
	}
	fmt.Printf("%v\n", filter)
	return filter
}

func CreateFloatFilter(key string, queries []string) (*FloatFilter, error) {
	fmt.Println("Key: " + key)
	fmt.Printf("%v\n", queries)
	var filter *FloatFilter
	if len(queries) > 0 {
		filter = &FloatFilter{
			Key:       key,
			Inclusive: false,
			Values:    nil,
		}
		if len(queries[0]) > 3 && queries[0][0:3] == "in:" {
			queries[0] = queries[0][3:]
			filter.Inclusive = true
		}

		var values []float64
		for _, item := range queries {
			val, err := strconv.ParseFloat(item, 10)
			if err != nil {
				return nil, fmt.Errorf("value has to be a floating number")
			}
			values = append(values, val)
		}
		filter.Values = values
	}
	fmt.Printf("%v\n", filter)
	return filter, nil
}