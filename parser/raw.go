package parser

import (
	"errors"
	"reflect"
	"restql/repository"
	"strings"
)

type RawModel struct {
	repo  repository.Repository
}

func NewRawModel(repo repository.Repository) Parser {
	return &RawModel{
		repo:  repo,
	}
}

func (mdl *RawModel) Command(data ModelColumn, model map[string]interface{}) (repo map[string]repository.Repository, err error) {
	for _, value := range data {
		if value.Insert != nil {

		}
		if value.Delete != nil {

		}
		if value.Update != nil {

		}
	}
	return
}

func (mdl *RawModel) Query(data ModelColumn, model map[string]interface{}) (repo map[string]repository.Repository, err error) {
	for key, value := range data {
		if value.Find != nil || value.Filter != nil || value.Where != nil {
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

func (mdl *RawModel) ToORM() (orm repository.IORM, err error) {
	orm, err = mdl.repo.ToORM()
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) filter(data ModelActions, model interface{}) (err error) {
	Filter := data.Filter
	where := data.Where
	find := data.Find

	if where != nil {
		Filter = where
	} else if find != nil {
		Filter = find
	}

	err = mdl.repo.Filter(Filter.(repository.IFilter), model)

	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) sortby(data ModelActions, model interface{}) (err error) {

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

func (mdl *RawModel) limit(data ModelActions) (err error) {
	// Limit
	var limit int = 10
	if data.Limit > 0 {
		limit = data.Limit
	}

	err = mdl.repo.Limit(int64(limit))
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) offset(data ModelActions) (err error) {
	strOffset := data.Offset
	strSkip := data.Skip

	offset := 0
	if strOffset > 0 {
		offset = strOffset
	} else if strSkip > 0 {
		offset = strSkip
	}

	err = mdl.repo.Offset(int64(offset))
	if err != nil {
		return
	}
	return
}

func (mdl *RawModel) selects(data ModelActions, model interface{}) (err error) {
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
