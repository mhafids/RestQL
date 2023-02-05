package repo

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/mhafids/RestQL/constants"
	"github.com/mhafids/RestQL/repository"
	"github.com/mhafids/RestQL/utils"
)

func (query *Repo) filterDB(filter repository.IFilter, model interface{}) (filterProcessed repository.IFilterProcessed, err error) {
	userFields := &bytes.Buffer{}
	userFields.Reset()

	valuebuffer := &bytes.Buffer{}
	valuebuffer.Reset()

	fields := utils.Bufpool.Get().(*bytes.Buffer)
	// fields := new(bytes.Buffer)
	fields.Reset()
	defer utils.Bufpool.Put(fields)

	err = query.getFields(userFields, model)
	if err != nil {
		return
	}

	err = query.operatorComparison(filter, valuebuffer, userFields, fields)
	if err != nil {
		return
	}

	valuebytes := bytes.Split(valuebuffer.Bytes(), []byte{0x90})
	values := make([]string, len(valuebytes))
	for index, valuebyte := range valuebytes {
		values[index] = string(valuebyte)
	}

	// values := strings.Split(valuebuffer.String(), "|")
	values = values[:len(values)-1]
	filterProcessed = repository.IFilterProcessed{
		Field:  fields.String(),
		Values: values,
	}
	return
}

func (query *Repo) operatorComparison(filter repository.IFilter, values, model, fields *bytes.Buffer) (err error) {

	switch filter.Operator {
	case constants.EQ:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}
			// values = append(values, filter.Value)
			fields.WriteString(filter.Field)
			fields.WriteString(" = ?")
			// fields += filter.Field + " = ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NE:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}
			// fields += filter.Field + " <> ?"
			fields.WriteString(filter.Field)
			fields.WriteString(" <> ?")
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LIKE:
		if query.stringInSlice(model, filter.Field) {
			values.WriteByte('%')
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('%')
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}
			// values = append(values, "%"+fmt.Sprintf("%v", filter.Value)+"%")
			fields.WriteString(filter.Field)
			fields.WriteString(" LIKE ?")
			// fields += filter.Field + " LIKE ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.ILIKE:
		if query.stringInSlice(model, filter.Field) {
			values.WriteByte('%')
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			values.WriteByte('%')
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}

			fields.WriteString("LOWER(")
			fields.WriteString(filter.Field)
			fields.WriteString(") LIKE LOWER(?)")
			// fields += "LOWER(" + filter.Field + ") LIKE LOWER(?)"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.GT:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}

			fields.WriteString(filter.Field)
			fields.WriteString(" > ?")
			// fields += filter.Field + " > ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.GTE:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}

			fields.WriteString(filter.Field)
			fields.WriteString(" >= ?")
			// fields += filter.Field + " >= ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LT:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}

			fields.WriteString(filter.Field)
			fields.WriteString(" < ?")
			// fields += filter.Field + " < ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.LTE:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}

			fields.WriteString(filter.Field)
			fields.WriteString(" <= ?")
			// fields += filter.Field + " <= ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NIN:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}

			fields.WriteString(filter.Field)
			fields.WriteString(" NOT IN ?")
			// fields += filter.Field + " NOT IN ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.IN:
		if query.stringInSlice(model, filter.Field) {
			values.WriteString(fmt.Sprintf("%v", filter.Value))
			err = values.WriteByte(0x90)
			if err != nil {
				return
			}

			fields.WriteString(filter.Field)
			fields.WriteString(" IN ?")
			// fields += filter.Field + " IN ?"
		} else {
			err = errors.New("\"" + filter.Field + "\" Not Found in model")
			return
		}

	case constants.NOT:
		if len(filter.Items) > 0 {
			// var fieldsNAND []string
			var fieldsNAND = &bytes.Buffer{}
			fieldsNAND.Reset()
			var saves = 0
			for _, item := range filter.Items {
				// var field = &bytes.Buffer{}
				// field.Reset()
				field := utils.Bufpool.Get().(*bytes.Buffer)
				field.Reset()
				errl := query.operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}

				if field.Len() == 0 {
					continue
				}

				if saves > 0 {
					fieldsNAND.WriteString(" AND NOT ")
				}

				fieldsNAND.Write(field.Bytes())
				saves++
				utils.Bufpool.Put(field)
			}

			if err != nil {
				return
			}

			if saves > 1 {
				fields.WriteByte('(')
				fields.Write(fieldsNAND.Bytes())
				fields.WriteByte(')')
			} else if saves > 0 {
				fields.Write(fieldsNAND.Bytes())
			}
		}

	case constants.NOR:
		if len(filter.Items) > 0 {
			// var fieldsNOR []string
			var fieldsNOR = &bytes.Buffer{}
			fieldsNOR.Reset()
			var saves = 0
			for _, item := range filter.Items {
				// var field = &bytes.Buffer{}
				// field.Reset()
				field := utils.Bufpool.Get().(*bytes.Buffer)
				field.Reset()
				errl := query.operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}

				if field.Len() == 0 {
					continue
				}

				if saves > 0 {
					fieldsNOR.WriteString(" OR NOT ")
				}

				fieldsNOR.Write(field.Bytes())
				saves++
				utils.Bufpool.Put(field)
				// fieldsNOR = append(fieldsNOR, field.String())
			}

			if err != nil {
				return
			}

			if saves > 1 {
				fields.WriteByte('(')
				fields.Write(fieldsNOR.Bytes())
				fields.WriteByte(')')
			} else if saves > 0 {
				fields.Write(fieldsNOR.Bytes())
			}
			// if len(fieldsNOR) > 0 {
			// 	fields.WriteString("( NOT " + strings.Join(fieldsNOR, " OR NOT ") + ")")
			// } else if len(fieldsNOR) > 0 {
			// 	fields.WriteString(fieldsNOR[0])
			// }
		}

	case constants.AND:
		if len(filter.Items) > 0 {
			var fieldsAND = &bytes.Buffer{}
			fieldsAND.Reset()
			// var fieldsAND = &bytes.Buffer{}
			var saves = 0
			for _, item := range filter.Items {
				// var field = &bytes.Buffer{}
				field := utils.Bufpool.Get().(*bytes.Buffer)
				field.Reset()
				errl := query.operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}

				if field.Len() == 0 {
					continue
				}

				if saves > 0 {
					fieldsAND.WriteString(" AND ")
				}

				fieldsAND.Write(field.Bytes())
				saves++
				utils.Bufpool.Put(field)
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
			// var fieldsOR []string
			var fieldsOR = &bytes.Buffer{}
			fieldsOR.Reset()

			var saves = 0
			for _, item := range filter.Items {
				// var field = &bytes.Buffer{}
				// field.Reset()
				field := utils.Bufpool.Get().(*bytes.Buffer)
				field.Reset()
				errl := query.operatorComparison(item, values, model, field)

				if errl != nil {
					err = errl
					break
				}
				if field.Len() == 0 {
					continue
				}

				if saves > 0 {
					fieldsOR.WriteString(" OR ")
				}

				fieldsOR.Write(field.Bytes())
				saves++
				utils.Bufpool.Put(field)
				// fieldsOR = append(fieldsOR, field.String())
			}

			if err != nil {
				return
			}

			if saves > 1 {
				fields.WriteByte('(')
				fields.Write(fieldsOR.Bytes())
				fields.WriteByte(')')
			} else if saves > 0 {
				fields.Write(fieldsOR.Bytes())
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
