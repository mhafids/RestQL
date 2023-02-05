package repository

import "github.com/mhafids/RestQL/constants"

type IFilter struct {
	Operator constants.Operator `json:"op"`
	Field    string             `json:"field"`
	Items    []IFilter          `json:"items"`
	Value    interface{}        `json:"value"`
}

type IFilterProcessed struct {
	Field  string
	Values []any
}

type ISortBy struct {
	Field string
	Sort  string
}
