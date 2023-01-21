package repo

import (
	"RestQL/constants"
	"RestQL/repository"
	"errors"
	"reflect"
	"strings"
)

func (query *Repo) filterDB(filter repository.IFilter, model interface{}) (filterProcessed repository.IFilterProcessed, err error) {
	var userFields = getFields(model)
	fields, values, err := operatorComparison(filter, userFields)
	if err != nil {
		return
	}
	filterProcessed = repository.IFilterProcessed{
		Field:  fields,
		Values: values,
	}
	return
}

func operatorComparison(filter repository.IFilter, model []string) (fields string, values []interface{}, err error) {

	switch strings.ToLower(filter.Operator) {
	case constants.EQ:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " = ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NE:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " <> ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LIKE:
		if stringInSlice(model, filter.Field) {
			values = append(values, "%"+filter.Value.(string)+"%")
			fields += filter.Field + " LIKE ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.ILIKE:
		if stringInSlice(model, filter.Field) {
			values = append(values, "%"+filter.Value.(string)+"%")
			fields += "LOWER(" + filter.Field + ") LIKE LOWER(?)"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.GT:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " > ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.GTE:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " >= ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LT:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " < ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LTE:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " <= ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NIN:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " NOT IN ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.IN:
		if stringInSlice(model, filter.Field) {
			values = append(values, filter.Value)
			fields += filter.Field + " IN ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NOT:
		if len(filter.Items) > 0 {
			var fieldsNAND []string
			for _, item := range filter.Items {
				field, value, errl := operatorComparison(item, model)

				if errl != nil {
					err = errl
					break
				}
				if field == "" {
					continue
				}
				fieldsNAND = append(fieldsNAND, field)
				values = append(values, value...)
			}

			if err != nil {
				return
			}

			if len(fieldsNAND) > 0 {
				fields = "( NOT " + strings.Join(fieldsNAND, " AND NOT ") + ")"
			}
		}

	case constants.NOR:
		if len(filter.Items) > 0 {
			var fieldsNOR []string
			for _, item := range filter.Items {
				field, value, errl := operatorComparison(item, model)

				if errl != nil {
					err = errl
					break
				}
				if field == "" {
					continue
				}
				fieldsNOR = append(fieldsNOR, field)
				values = append(values, value...)
			}

			if err != nil {
				return
			}

			if len(fieldsNOR) > 0 {
				fields = "( NOT " + strings.Join(fieldsNOR, " OR NOT ") + ")"
			}
		}

	case constants.AND:
		if len(filter.Items) > 0 {
			var fieldsAND []string
			for _, item := range filter.Items {
				field, value, errl := operatorComparison(item, model)

				if errl != nil {
					err = errl
					break
				}
				if field == "" {
					continue
				}
				fieldsAND = append(fieldsAND, field)
				values = append(values, value...)
			}

			if err != nil {
				return
			}

			if len(fieldsAND) > 0 {
				fields = "(" + strings.Join(fieldsAND, " AND ") + ")"
			}
		}

	case constants.OR:
		if len(filter.Items) > 0 {
			var fieldsOR []string
			for _, item := range filter.Items {
				field, value, errl := operatorComparison(item, model)

				if errl != nil {
					err = errl
					break
				}
				if field == "" {
					continue
				}
				fieldsOR = append(fieldsOR, field)
				values = append(values, value...)
			}

			if err != nil {
				return
			}

			if len(fieldsOR) > 0 {
				fields = "(" + strings.Join(fieldsOR, " OR ") + ")"
			}
		}

	default:
		err = errors.New("Operator \"" + filter.Operator + "\" not available")
		if err != nil {
			return
		}
	}

	// fmt.Println(strings.ToLower(filter.Operator), fields)
	return
}

func getFields(Interfacefield interface{}) []string {
	var field []string
	v := reflect.ValueOf(Interfacefield)
	for i := 0; i < v.Type().NumField(); i++ {
		field = append(field, v.Type().Field(i).Tag.Get("json"))
	}
	return field
}
