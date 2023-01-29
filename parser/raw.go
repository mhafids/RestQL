package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/mhafids/RestQL/repository"
	"github.com/mhafids/RestQL/utils"
)

type RawModel struct {
	repo repository.Repository
}

type rawmodelactions struct {
	// Command
	Insert repository.Insert `json:"insert"`
	Create repository.Insert `json:"create"`
	Update repository.Update `json:"update"`
	Delete repository.Delete `json:"delete"`

	// QueryBatch
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
	}
}

// Command for raw command to parser
func (mdl *RawModel) CommandBatch(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error) {
	var datamodel map[string]rawmodelactions
	json.Unmarshal([]byte(data), &datamodel)
	data = ""

	for key, value := range datamodel {
		repo[key], err = mdl.command(value, model[key])
		if err != nil {
			return
		}
	}
	return
}

func (mdl *RawModel) Command(data string, model map[string]interface{}) (repo repository.Repository, err error) {

	return
}

func (mdl *RawModel) command(datamodel rawmodelactions, model interface{}) (repo repository.Repository, err error) {
	if len(datamodel.Insert.Datas) > 0 {

	}
	if datamodel.Delete.Find.Operator != "" || datamodel.Delete.Filter.Operator != "" || datamodel.Delete.Where.Operator != "" {

	}
	if datamodel.Update.Find.Operator != "" && len(datamodel.Update.Datas) > 0 || datamodel.Update.Filter.Operator != "" && len(datamodel.Update.Datas) > 0 || datamodel.Update.Where.Operator != "" && len(datamodel.Update.Datas) > 0 {

	}
	return
}

// QueryBatch for multiple raw query to parser
func (mdl *RawModel) QueryBatch(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error) {
	var datamodel map[string]rawmodelactions
	json.Unmarshal([]byte(data), &datamodel)
	data = ""

	repo = make(map[string]repository.Repository, 0)
	for key, value := range datamodel {
		repo[key], err = mdl.query(value, model[key])
		if err != nil {
			return
		}
	}

	return
}

// Query for single raw parser to Repository ORM
func (mdl *RawModel) Query(data string, model interface{}) (repo repository.Repository, err error) {
	datamodel := rawmodelactions{}
	json.Unmarshal([]byte(data), &datamodel)
	data = ""

	repo, err = mdl.query(datamodel, model)

	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) query(datamodel rawmodelactions, model interface{}) (repo repository.Repository, err error) {
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

/* Query */

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

	userFields := utils.Bufpool.Get().(*bytes.Buffer)
	userFields.Reset()
	defer utils.Bufpool.Put(userFields)

	mdl.getFields(userFields, model)
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
		userFields := utils.Bufpool.Get().(*bytes.Buffer)
		userFields.Reset()
		defer utils.Bufpool.Put(userFields)
		mdl.getFields(userFields, model)
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

func (mdl *RawModel) getFields(buffer *bytes.Buffer, Interfacefield interface{}) {
	v := reflect.ValueOf(Interfacefield)
	for i := 0; i < v.Type().NumField(); i++ {
		buffer.WriteString(v.Type().Field(i).Tag.Get("json"))
		if i+1 < v.Type().NumField() {
			buffer.WriteByte('.')
		}
	}
}

func (mdl *RawModel) stringInSlice(bufferfield *bytes.Buffer, s string) bool {
	return bytes.Contains(bufferfield.Bytes(), []byte(s))
}
