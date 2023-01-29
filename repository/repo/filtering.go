package repo

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/mhafids/RestQL/constants"
	"github.com/mhafids/RestQL/repository"
	"github.com/mhafids/RestQL/utils"
)

func (query *Repo) filterDB(filter repository.IFilter, model interface{}) (filterProcessed repository.IFilterProcessed, err error) {
	userFields := utils.Bufpool.Get().(*bytes.Buffer)
	userFields.Reset()
	defer utils.Bufpool.Put(userFields)

	valuebuffer := utils.Bufpool.Get().(*bytes.Buffer)
	valuebuffer.Reset()
	defer utils.Bufpool.Put(valuebuffer)

	fields := utils.Bufpool.Get().(*bytes.Buffer)
	fields.Reset()
	defer utils.Bufpool.Put(fields)

	getFields(userFields, model)
	err = operatorComparison(filter, valuebuffer, userFields, fields)
	if err != nil {
		return
	}

	values := strings.Split(valuebuffer.String(), ".")
	values = values[:len(values)-1]
	filterProcessed = repository.IFilterProcessed{
		Field:  fields.String(),
		Values: values,
	}
	return
}

func operatorComparison(filter repository.IFilter, values, model, fields *bytes.Buffer) (err error) {

	switch filter.Operator {
	case constants.EQ:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')
			// values = append(values, filter.Value)
			fields.WriteString(filter.Field)
			fields.WriteString(" = ?")
			// fields += filter.Field + " = ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NE:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')
			// fields += filter.Field + " <> ?"
			fields.WriteString(filter.Field)
			fields.WriteString(" <> ?")
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LIKE:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", "%"+filter.Value.(string)+"%"))
			values.WriteByte('.')
			// values = append(values, "%"+filter.Value.(string)+"%")
			fields.WriteString(filter.Field)
			fields.WriteString(" LIKE ?")
			// fields += filter.Field + " LIKE ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.ILIKE:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", "%"+filter.Value.(string)+"%"))
			values.WriteByte('.')

			fields.WriteString("LOWER(")
			fields.WriteString(filter.Field)
			fields.WriteString(") LIKE LOWER(?)")
			// fields += "LOWER(" + filter.Field + ") LIKE LOWER(?)"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.GT:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')

			fields.WriteString(filter.Field)
			fields.WriteString(" > ?")
			// fields += filter.Field + " > ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.GTE:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')

			fields.WriteString(filter.Field)
			fields.WriteString(" >= ?")
			// fields += filter.Field + " >= ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LT:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')

			fields.WriteString(filter.Field)
			fields.WriteString(" < ?")
			// fields += filter.Field + " < ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LTE:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')

			fields.WriteString(filter.Field)
			fields.WriteString(" <= ?")
			// fields += filter.Field + " <= ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NIN:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')

			fields.WriteString(filter.Field)
			fields.WriteString(" NOT IN ?")
			// fields += filter.Field + " NOT IN ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.IN:
		if stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('.')

			fields.WriteString(filter.Field)
			fields.WriteString(" IN ?")
			// fields += filter.Field + " IN ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NOT:
		if len(filter.Items) > 0 {
			var fieldsNAND []string
			for _, item := range filter.Items {
				field := &bytes.Buffer{}
				field.Reset()
				errl := operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}

				if field.Len() > 0 {
					continue
				}

				fieldsNAND = append(fieldsNAND, field.String())
			}

			if err != nil {
				return
			}

			if len(fieldsNAND) > 0 {
				fields.WriteString("( NOT " + strings.Join(fieldsNAND, " AND NOT ") + ")")
			} else if len(fieldsNAND) > 0 {
				fields.WriteString(fieldsNAND[0])
			}
		}

	case constants.NOR:
		if len(filter.Items) > 0 {
			var fieldsNOR []string
			for _, item := range filter.Items {
				field := &bytes.Buffer{}
				field.Reset()
				errl := operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}

				if field.Len() > 0 {
					continue
				}
				fieldsNOR = append(fieldsNOR, field.String())
			}

			if err != nil {
				return
			}

			if len(fieldsNOR) > 0 {
				fields.WriteString("( NOT " + strings.Join(fieldsNOR, " OR NOT ") + ")")
			} else if len(fieldsNOR) > 0 {
				fields.WriteString(fieldsNOR[0])
			}
		}

	case constants.AND:
		if len(filter.Items) > 0 {
			var fieldsAND = &bytes.Buffer{}
			var saves = 0
			for _, item := range filter.Items {
				field := &bytes.Buffer{}
				field.Reset()
				errl := operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}

				if field.Len() > 0 {
					continue
				}

				if saves > 0 {
					fieldsAND.WriteString(" AND ")
				}

				fieldsAND.Write(field.Bytes())
				saves++
				// fieldsAND = append(fieldsAND, field.String())
			}

			if err != nil {
				return
			}

			if saves > 1 {
				fields.WriteByte('(')
				fields.Write(fieldsAND.Bytes())
				fields.WriteByte(')')
			} else if saves > 0 {
				fields.Write(fieldsAND.Bytes())
			}
		}

	case constants.OR:
		if len(filter.Items) > 0 {
			var fieldsOR []string
			for _, item := range filter.Items {
				field := &bytes.Buffer{}
				field.Reset()
				errl := operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}
				if field.Len() > 0 {
					continue
				}
				fieldsOR = append(fieldsOR, field.String())
			}

			if err != nil {
				return
			}

			if len(fieldsOR) > 0 {
				fields.WriteString("(" + strings.Join(fieldsOR, " OR ") + ")")
			} else if len(fieldsOR) > 0 {
				fields.WriteString(fieldsOR[0])
			}
		}

	default:
		err = errors.New("Operator \"" + string(filter.Operator) + "\" not available")
		if err != nil {
			return
		}
	}

	// fmt.Println(strings.ToLower(filter.Operator), fields)
	return
}

func getFields(buffer *bytes.Buffer, Interfacefield interface{}) {
	v := reflect.ValueOf(Interfacefield)
	for i := 0; i < v.Type().NumField(); i++ {
		buffer.WriteString(v.Type().Field(i).Tag.Get("json"))
		if i+1 < v.Type().NumField() {
			buffer.WriteByte('.')
		}
	}
}
