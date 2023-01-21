package parser

import "RestQL/repository"

type ModelColumn map[string]ModelActions

type Parser interface {
	QueryOne(data ModelActions, model interface{}) (repo repository.Repository, err error)
	Query(data ModelColumn, model map[string]interface{}) (repo map[string]repository.Repository, err error)
	Command(data ModelColumn, model map[string]interface{}) (repo map[string]repository.Repository, err error)
}

type ModelActions struct {
	// Command
	Insert interface{} `json:"insert"`
	Update interface{} `json:"update"`
	Delete interface{} `json:"delete"`

	// Query
	Find   interface{} `json:"find"`
	Filter interface{} `json:"filter"`
	Where  interface{} `json:"where"`

	Sortby  string `json:"sortby"`
	Sort    string `json:"sort"`
	Orderby string `json:"orderby"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`
	Skip   int `json:"skip"`

	Select []string `json:"select"`
}
