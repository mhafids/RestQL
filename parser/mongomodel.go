package parser

import (
	"errors"
	"reflect"
	"strings"
)

type paramMongoTranslate struct {
	data interface{}
	key  string
}

type MongoModel struct {
}

func NewMongoModel() *MongoModel {
	return &MongoModel{}
}

func (mts *MongoModel) Query(data interface{}) (documentquery map[string]IFilter, err error) {
	documentquery = make(map[string]IFilter, 0)
	if reflect.TypeOf(data).Kind() == reflect.Map {
		for key, value := range data.(map[string]interface{}) {
			doc, errl := mts.parser(paramMongoTranslate{data: value})
			if errl != nil {
				err = errl
				break
			}
			documentquery[key] = doc
		}

		if err != nil {
			return
		}
		return
	} else {
		err = errors.New("query must objects")
		return
	}
	return
}

func (mts *MongoModel) parser(param paramMongoTranslate) (filtering IFilter, err error) {
	data := param.data
	rt := reflect.TypeOf(data)
	switch rt.Kind() {
	case reflect.Slice:
		items, errl := mts.loopmap(param)
		if errl != nil {
			err = errl
			return
		}
		filtering = IFilter{
			Operator: AND,
			Items:    items,
		}
		return

	case reflect.Array:
		items, errl := mts.loopmap(param)
		if errl != nil {
			err = errl
			return
		}
		filtering = IFilter{
			Operator: AND,
			Items:    items,
		}
		return

	case reflect.Map:
		for key, value := range data.(map[string]interface{}) {
			switch key {
			case AND:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = IFilter{
					Operator: AND,
					Items:    items,
				}
				return
			case OR:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = IFilter{
					Operator: OR,
					Items:    items,
				}
				return
			case NOT:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = IFilter{
					Operator: NOT,
					Items:    items,
				}
				return
			case NOR:
				items, errl := mts.loopmap(paramMongoTranslate{data: value})
				if errl != nil {
					err = errl
					return
				}
				filtering = IFilter{
					Operator: NOR,
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
					filtering = IFilter{
						Operator: AND,
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
						filtering = IFilter{
							Operator: EQ,
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

func (mts *MongoModel) loopmap(param paramMongoTranslate) (filtering []IFilter, err error) {
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

func (mts *MongoModel) operatorcase(key string, upperkey string, value interface{}) (filtering IFilter, err error) {
	switch strings.ToLower(key) {
	case GT:
		filtering = IFilter{
			Operator: GT,
			Field:    upperkey,
			Value:    value,
		}
		return

	case GTE:
		filtering = IFilter{
			Operator: GTE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case ILIKE:
		filtering = IFilter{
			Operator: ILIKE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case LIKE:
		filtering = IFilter{
			Operator: LIKE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case IN:
		filtering = IFilter{
			Operator: IN,
			Field:    upperkey,
			Value:    value,
		}
		return

	case LT:
		filtering = IFilter{
			Operator: LT,
			Field:    upperkey,
			Value:    value,
		}
		return

	case LTE:
		filtering = IFilter{
			Operator: LTE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case NE:
		filtering = IFilter{
			Operator: NE,
			Field:    upperkey,
			Value:    value,
		}
		return

	case NIN:
		filtering = IFilter{
			Operator: NIN,
			Field:    upperkey,
			Value:    value,
		}
		return

	case AND:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = IFilter{
			Operator: AND,
			Items:    items,
		}
		return
	case OR:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = IFilter{
			Operator: OR,
			Items:    items,
		}
		return
	case NOT:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = IFilter{
			Operator: NOT,
			Items:    items,
		}
		return
	case NOR:
		items, errl := mts.loopmap(paramMongoTranslate{data: value, key: upperkey})
		if errl != nil {
			err = errl
			return
		}
		filtering = IFilter{
			Operator: NOR,
			Items:    items,
		}
		return

	default:
		err = errors.New("\"" + key + "\" operator not found")
		return
	}

	return
}
