package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/mhafids/RestQL/repository"
)

type RawModel struct {
	repo    repository.Repository
	bufpool sync.Pool
}

type rawmodelactions struct {
	// Command
	Insert interface{} `json:"insert"`
	Update interface{} `json:"update"`
	Delete interface{} `json:"delete"`

	// Query
	Find   repository.IFilter `json:"find"`
	Filter repository.IFilter `json:"filter"`
	Where  repository.IFilter `json:"where"`

	Sortby  string `json:"sortby"`
	Sort    string `json:"sort"`
	Orderby string `json:"orderby"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`
	Skip   int `json:"skip"`

	Select []string `json:"select"`
}

// NewRawModel initial new parser Rawmodel
func NewRawModel(repo repository.Repository) Parser {
	return &RawModel{
		repo: repo,
		bufpool: sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}
}

// Command for raw command to parser
func (mdl *RawModel) Command(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error) {
	var datamodel map[string]rawmodelactions
	json.Unmarshal([]byte(data), &datamodel)
	data = ""


	for _, value := range datamodel {
		if value.Insert != nil {

		}
		if value.Delete != nil {

		}
		if value.Update != nil {

		}
	}
	return
}

// Query for multiple raw query to parser
func (mdl *RawModel) Query(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error) {

	var datamodel map[string]rawmodelactions
	json.Unmarshal([]byte(data), &datamodel)
	data = ""

	repo = make(map[string]repository.Repository, 0)
	for key, value := range datamodel {
		if value.Find.Operator != "" || value.Filter.Operator != "" || value.Where.Operator != "" {
			err = mdl.filter(value, model[key])
			if err != nil {
				return
			}
		}

		if value.Orderby != "" || value.Sort != "" || value.Sortby != "" {
			err = mdl.sortby(value, model[key])
			if err != nil {
				return
			}
		}

		if value.Limit > 0 {
			err = mdl.limit(value)
			if err != nil {
				return
			}
		}

		if value.Offset > 0 || value.Skip > 0 {
			err = mdl.offset(value)
			if err != nil {
				return
			}
		}

		if len(value.Select) > 0 {
			err = mdl.selects(value, model[key])
			if err != nil {
				return
			}
		}

		repo[key] = mdl.repo
	}

	return
}

// QueryOne for single raw parser to Repository ORM
func (mdl *RawModel) QueryOne(data string, model interface{}) (repo repository.Repository, err error) {
	var datamodel rawmodelactions
	json.Unmarshal([]byte(data), &datamodel)
	data = ""

	if datamodel.Find.Operator != "" || datamodel.Filter.Operator != "" || datamodel.Where.Operator != "" {
		err = mdl.filter(datamodel, model)
		if err != nil {
			return
		}
	}

	if datamodel.Orderby != "" || datamodel.Sort != "" || datamodel.Sortby != "" {
		err = mdl.sortby(datamodel, model)
		if err != nil {
			return
		}
	}

	if datamodel.Limit > 0 {
		err = mdl.limit(datamodel)
		if err != nil {
			return
		}
	}

	if datamodel.Offset > 0 || datamodel.Skip > 0 {
		err = mdl.offset(datamodel)
		if err != nil {
			return
		}
	}

	if len(datamodel.Select) > 0 {
		err = mdl.selects(datamodel, model)
		if err != nil {
			return
		}
	}

	repo = mdl.repo

	return
}

// ToORM raw parser to Repository ORM
func (mdl *RawModel) ToORM() (orm repository.IORM, err error) {
	orm, err = mdl.repo.ToORM()
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) filter(data rawmodelactions, model interface{}) (err error) {
	Filter := data.Filter

	if data.Where.Operator != "" {
		Filter = data.Where
	} else if data.Find.Operator != "" {
		Filter = data.Find
	}

	err = mdl.repo.Filter(Filter, model)
	Filter = repository.IFilter{}

	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) sortby(data rawmodelactions, model interface{}) (err error) {

	// Sort
	sortBy := ""

	if data.Sort != "" {
		sortBy = data.Sort
	}

	if data.Orderby != "" {
		sortBy = data.Orderby
	}

	if sortBy == "" {
		// id.asc is the default sort query
		sortBy = "id asc"
	}

	var userFields = mdl.getFields(model)
	sortBy = strings.ReplaceAll(sortBy, ", ", ",")
	commasplit := strings.Split(sortBy, ",")
	var orderby []repository.ISortBy

	for _, cs := range commasplit {
		splits := strings.Split(cs, " ")
		if len(splits) != 2 {
			splits = append(splits, "asc")
		}
		field, order := splits[0], splits[1]
		order = strings.ToLower(order)

		if order != "desc" && order != "asc" {
			err = errors.New("malformed order direction in sortBy query parameter, should be asc or desc")
			return
		}

		if !mdl.stringInSlice(userFields, field) && field != "id" {
			err = errors.New("unknown field in sortBy query parameter")
			return
		}

		orderby = append(orderby, repository.ISortBy{
			Field: field,
			Sort:  strings.ToUpper(order),
		})
	}

	err = mdl.repo.SortBy(orderby, model)
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) limit(data rawmodelactions) (err error) {
	// Limit
	var limit int = 10
	if data.Limit > 0 {
		limit = data.Limit
	}

	err = mdl.repo.Limit(limit)
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) offset(data rawmodelactions) (err error) {
	strOffset := data.Offset
	strSkip := data.Skip

	offset := 0
	if strOffset > 0 {
		offset = strOffset
	} else if strSkip > 0 {
		offset = strSkip
	}

	err = mdl.repo.Offset(offset)
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) selects(data rawmodelactions, model interface{}) (err error) {
	selectcheck := func(selects []string, model interface{}) error {
		var userFields = mdl.getFields(model)
		for _, Select := range selects {
			if !mdl.stringInSlice(userFields, Select) {
				return errors.New(Select + " field not found")
			}
		}

		return nil
	}

	err = selectcheck(data.Select, mdl)
	if err != nil {
		return
	}

	err = mdl.repo.Select(data.Select, model)
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) getFields(Interfacefield interface{}) []string {
	var field []string
	v := reflect.ValueOf(Interfacefield)
	for i := 0; i < v.Type().NumField(); i++ {
		field = append(field, v.Type().Field(i).Tag.Get("json"))
	}
	return field
}

func (mdl *RawModel) stringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}

// func (mdl *RawModel) transcode(in, out interface{}) {

// 	buf := mdl.bufpool.Get().(*bytes.Buffer)
// 	json.NewEncoder(buf).Encode(in)
// 	json.NewDecoder(buf).Decode(out)
// 	buf.Reset()
// 	mdl.bufpool.Put(buf)
// 	in = nil
// }

func (mdl *RawModel) transcode(in map[string]interface{}, out *repository.IFilter) (err error) {
	data, err := json.Marshal(in)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return
	}
	return
}
