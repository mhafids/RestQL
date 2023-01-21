package parser

import (
	"errors"
	"reflect"
	"restql/constants"
	"restql/repository"
	"strings"
)

type paramMongoTranslate struct {
	data interface{}
	key  string
}

type MongoModel struct {
	repo repository.Repository
}

func NewMongoModel(repo repository.Repository) Parser {
	return &MongoModel{
		repo: repo,
	}
}

func (mts *MongoModel) Command(data ModelColumn, model map[string]interface{}) (repo map[string]repository.Repository, err error) {

	return
}

func (mts *MongoModel) Query(data ModelColumn, model map[string]interface{}) (repo map[string]repository.Repository, err error) {
	repo = make(map[string]repository.Repository, 0)
	for key, value := range data {
		if value.Find != nil || value.Filter != nil || value.Where != nil {
			err = mts.filter(value, model[key])
			if err != nil {
				return
			}
		}

		if value.Orderby != "" || value.Sort != "" || value.Sortby != "" {
			err = mts.sortby(value, model[key])
			if err != nil {
				return
			}
		}

		if value.Limit > 0 {
			err = mts.limit(value)
			if err != nil {
				return
			}
		}

		if value.Offset > 0 || value.Skip > 0 {
			err = mts.offset(value)
			if err != nil {
				return
			}
		}

		if len(value.Select) > 0 {
			err = mts.selects(value, model[key])
			if err != nil {
				return
			}
		}

		repo[key] = mts.repo
	}
	return
}

func (mts *MongoModel) QueryOne(data ModelActions, model interface{}) (repo repository.Repository, err error) {

	if data.Find != nil || data.Filter != nil || data.Where != nil {
		err = mts.filter(data, model)
		if err != nil {
			return
		}
	}

	if data.Orderby != "" || data.Sort != "" || data.Sortby != "" {
		err = mts.sortby(data, model)
		if err != nil {
			return
		}
	}

	if data.Limit > 0 {
		err = mts.limit(data)
		if err != nil {
			return
		}
	}

	if data.Offset > 0 || data.Skip > 0 {
		err = mts.offset(data)
		if err != nil {
			return
		}
	}

	if len(data.Select) > 0 {
		err = mts.selects(data, model)
		if err != nil {
			return
		}
	}

	repo = mts.repo

	return
}

func (mts *MongoModel) filter(data ModelActions, model interface{}) (err error) {
	Filter := data.Filter
	where := data.Where
	find := data.Find

	if where != nil {
		Filter = where
	} else if find != nil {
		Filter = find
	}

	filter, err := mts.parser(paramMongoTranslate{data: Filter})
	if err != nil {
		return
	}

	err = mts.repo.Filter(filter, model)
	if err != nil {
		return
	}
	return
}

func (mts *MongoModel) sortby(data ModelActions, model interface{}) (err error) {

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

	var userFields = mts.getFields(model)
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

		if !mts.stringInSlice(userFields, field) && field != "id" {
			err = errors.New("unknown field in sortBy query parameter")
			return
		}

		orderby = append(orderby, repository.ISortBy{
			Field: field,
			Sort:  strings.ToUpper(order),
		})
	}

	err = mts.repo.SortBy(orderby, model)
	if err != nil {
		return
	}
	return
}

func (mts *MongoModel) limit(data ModelActions) (err error) {
	// Limit
	var limit int = 10
	if data.Limit > 0 {
		limit = data.Limit
	}

	err = mts.repo.Limit(limit)
	if err != nil {
		return
	}
	return
}

func (mts *MongoModel) offset(data ModelActions) (err error) {
	strOffset := data.Offset
	strSkip := data.Skip

	offset := 0
	if strOffset > 0 {
		offset = strOffset
	} else if strSkip > 0 {
		offset = strSkip
	}

	err = mts.repo.Offset(offset)
	if err != nil {
		return
	}
	return
}

func (mts *MongoModel) selects(data ModelActions, model interface{}) (err error) {
	selectcheck := func(selects []string, model interface{}) error {
		var userFields = mts.getFields(model)
		for _, Select := range selects {
			if !mts.stringInSlice(userFields, Select) {
				return errors.New(Select + " field not found")
			}
		}

		return nil
	}

	err = selectcheck(data.Select, mts)
	if err != nil {
		return
	}

	err = mts.repo.Select(data.Select, model)
	if err != nil {
		return
	}
	return
}

func (mts *MongoModel) parser(param paramMongoTranslate) (filtering repository.IFilter, err error) {
	data := param.data
	rt := reflect.TypeOf(data)
	switch rt.Kind() {
	case reflect.Slice:
		items, errl := mts.loopmap(param)
		if errl != nil {
			err = errl
			return
		}
		filtering = repository.IFilter{
			Operator: constants.AND,
			Items:    items,
		}
		return

	case reflect.Array:
		items, errl := mts.loopmap(param)
		if errl != nil {
			err = errl
			return
		}
		filtering = repository.IFilter{
			Operator: constants.AND,
			Items:    items,
		}
		return

	case reflect.Map:
		for key, value := range data.(map[string]interface{}) {
			switch key {
			case constants.AND:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = repository.IFilter{
					Operator: constants.AND,
					Items:    items,
				}
				return
			case constants.OR:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = repository.IFilter{
					Operator: constants.OR,
					Items:    items,
				}
				return
			case constants.NOT:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = repository.IFilter{
					Operator: constants.NOT,
					Items:    items,
				}
				return
			case constants.NOR:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = repository.IFilter{
					Operator: constants.NOR,
					Items:    items,
				}
				return

			default:
				typeofvalue := reflect.TypeOf(value)
				if typeofvalue.Kind() == reflect.Map {
					items, errl := mts.loopmap(paramMongoTranslate{data: value, key: key})
					if errl != nil {
						err = errl
						return
					}
					filtering = repository.IFilter{
						Operator: constants.AND,
						Items:    items,
					}
					return
				} else {
					if param.key != "" {
						filtering, err = mts.operatorcase(key, param.key, value)
						if err != nil {
							return
						}
						return
					} else {
						filtering = repository.IFilter{
							Operator: constants.EQ,
							Field:    key,
							Value:    value,
						}
						return
					}
				}
			}
		}
		return
	}

	return
}

func (mts *MongoModel) loopmap(param paramMongoTranslate) (filtering []repository.IFilter, err error) {
	data := param.data
	rt := reflect.TypeOf(data)
	switch rt.Kind() {
	case reflect.Slice:
		for _, dataarr := range data.([]interface{}) {
			filter, errl := mts.parser(paramMongoTranslate{data: dataarr, key: param.key})
			if errl != nil {
				err = errl
				return
			}

			filtering = append(filtering, filter)
		}
		return
	case reflect.Array:
		for _, dataarr := range data.([]interface{}) {
			filter, errl := mts.parser(paramMongoTranslate{data: dataarr, key: param.key})
			if errl != nil {
				err = errl
				return
			}
			filtering = append(filtering, filter)
		}
		return
	case reflect.Map:
		for key, value := range data.(map[string]interface{}) {
			if param.key != "" {
				filter, errl := mts.operatorcase(key, param.key, value)
				if errl != nil {
					err = errl
					return
				}
				filtering = append(filtering, filter)
			}
		}
		return
	}
	return
}

func (mts *MongoModel) operatorcase(key string, upperkey string, value interface{}) (filtering repository.IFilter, err error) {
	switch strings.ToLower(key) {
	case constants.GT:
		filtering = repository.IFilter{
			Operator: constants.GT,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.GTE:
		filtering = repository.IFilter{
			Operator: constants.GTE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.ILIKE:
		filtering = repository.IFilter{
			Operator: constants.ILIKE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.LIKE:
		filtering = repository.IFilter{
			Operator: constants.LIKE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.IN:
		filtering = repository.IFilter{
			Operator: constants.IN,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.LT:
		filtering = repository.IFilter{
			Operator: constants.LT,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.LTE:
		filtering = repository.IFilter{
			Operator: constants.LTE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.NE:
		filtering = repository.IFilter{
			Operator: constants.NE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.NIN:
		filtering = repository.IFilter{
			Operator: constants.NIN,
			Field:    upperkey,
			Value:    value,
		}
		return

	case constants.AND:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = repository.IFilter{
			Operator: constants.AND,
			Items:    items,
		}
		return
	case constants.OR:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = repository.IFilter{
			Operator: constants.OR,
			Items:    items,
		}
		return
	case constants.NOT:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = repository.IFilter{
			Operator: constants.NOT,
			Items:    items,
		}
		return
	case constants.NOR:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = repository.IFilter{
			Operator: constants.NOR,
			Items:    items,
		}
		return

	default:
		err = errors.New("\"" + key + "\" operator not found")
		return
	}

	return
}

func (mts *MongoModel) getFields(Interfacefield interface{}) []string {
	var field []string
	v := reflect.ValueOf(Interfacefield)
	for i := 0; i < v.Type().NumField(); i++ {
		field = append(field, v.Type().Field(i).Tag.Get("json"))
	}
	return field
}

func (mts *MongoModel) stringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}
