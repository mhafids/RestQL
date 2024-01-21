package rawparser

import (
	"sync"

	"github.com/mhafids/RestQL/repository"
)

// NewRawModelParser initial new parser Rawmodel
func NewRawModelParser(repo repository.Repository) *RawModelParser {
	return &RawModelParser{
		repo:   repo,
		txrepo: repo,
		mtx:    sync.Mutex{},
	}
}

// RawModelParser is model for raw RestQuery
type RawModelParser struct {
	repo   repository.Repository
	txrepo repository.Repository
	mtx    sync.Mutex
}

// ToORM raw parser to Repository ORM
func (mdl *RawModelParser) ToORM() (orm repository.IORM, err error) {
	orm, err = mdl.repo.ToORM()
	if err != nil {
		return
	}
	return
}

type filtermodel struct {
	Find   repository.IFilter `json:"find"`
	Filter repository.IFilter `json:"filter"`
	Where  repository.IFilter `json:"where"`
}

type sortbymodel struct {
	Sortby  []string `json:"sortby"`
	Sort    []string `json:"sort"`
	Orderby []string `json:"orderby"`
}

type limitmodel struct {
	Limit int `json:"limit"`
}

type offsetmodel struct {
	Offset int `json:"offset"`
	Skip   int `json:"skip"`
}

type selectmodel struct {
	Select []string `json:"select"`
}

type groupbymodel struct {
	GroupBy []string `json:"group_by"`
}

type rawmodelactions struct {
	filtermodel
	sortbymodel
	limitmodel
	offsetmodel
	selectmodel
	groupbymodel
}

type updatemodel struct {
	Datas []repository.ICommand `json:"datas"`

	filtermodel
	sortbymodel
	limitmodel
}

type insertmodel struct {
	Datas []repository.ICommand `json:"datas"`
}

type deletemodel struct {
	Datas []repository.ICommand `json:"datas"`

	filtermodel
	sortbymodel
	limitmodel
	offsetmodel
}
