package rawparser

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/mhafids/RestQL/repository"

	"github.com/goccy/go-json"
)

func (mdl *RawModelParser) query(datamodel rawmodelactions, model interface{}) (repo repository.Repository) {
	mdl.mtx.Lock()
	mdl.txrepo = mdl.repo.Model(model)

	if datamodel.Find.Operator != "" || datamodel.Filter.Operator != "" || datamodel.Where.Operator != "" {
		mdl.filter(datamodel.filtermodel)
	}

	if len(datamodel.Sortby) > 0 || len(datamodel.Sort) > 0 || len(datamodel.Orderby) > 0 {
		mdl.sortby(datamodel.sortbymodel)
	}

	if datamodel.Limit > 0 {
		mdl.limit(datamodel.limitmodel)
	}

	if datamodel.Offset > 0 || datamodel.Skip > 0 {
		mdl.offset(datamodel.offsetmodel)
	}

	if len(datamodel.Select) > 0 {
		mdl.selects(datamodel.selectmodel)
	}

	if len(datamodel.GroupBy) > 0 {
		mdl.groupby(datamodel.groupbymodel)
	}

	mdl.mtx.Unlock()
	return
}

/* Query */
func (mdl *RawModelParser) filter(data filtermodel) {

	Filter := data.Filter

	if data.Where.Operator != "" {
		Filter = data.Where
	} else if data.Find.Operator != "" {
		Filter = data.Find
	}

	mdl.txrepo = mdl.txrepo.Filter(Filter)
	return
}

func (mdl *RawModelParser) sortby(data sortbymodel) {
	// Sort
	sortBy := []string{}

	if len(data.Sortby) > 0 {
		sortBy = data.Sortby
	}

	if len(data.Sort) > 0 {
		sortBy = data.Sort
	}

	if len(data.Orderby) > 0 {
		sortBy = data.Orderby
	}

	if len(sortBy) <= 0 {
		// id.asc is the default sort query
		sortBy = []string{"id asc"}
	}

	orders := []repository.ISortBy{}
	for _, sp := range sortBy {
		sp = strings.ReplaceAll(sp, ", ", ",")
		split := strings.Split(sp, " ")
		orders = append(orders, repository.ISortBy{
			Field: split[0],
			Sort:  split[1],
		})
	}

	mdl.txrepo = mdl.txrepo.SortBy(orders)
	return
}

func (mdl *RawModelParser) limit(data limitmodel) {
	// Limit
	var limit int = 10
	if data.Limit > 0 {
		limit = data.Limit
	}

	mdl.txrepo = mdl.txrepo.Limit(limit)
	return
}

func (mdl *RawModelParser) offset(data offsetmodel) {
	strOffset := data.Offset
	strSkip := data.Skip

	offset := 0
	if strOffset > 0 {
		offset = strOffset
	} else if strSkip > 0 {
		offset = strSkip
	}

	mdl.txrepo = mdl.txrepo.Offset(offset)
	return
}

func (mdl *RawModelParser) selects(data selectmodel) {
	mdl.txrepo = mdl.txrepo.Select(data.Select)
	return
}

func (mdl *RawModelParser) groupby(data groupbymodel) {
	mdl.txrepo = mdl.txrepo.GroupBy(data.GroupBy)
	return
}

// // Command
// func (mdl *RawModelParser) insert(data rawmodelactions, model interface{}) (repo repository.Repository, err error) {
// 	datainsert := repository.Insert{}
// 	if len(data.Create.Datas) > 0 {
// 		datainsert = data.Create
// 	} else if len(data.Save.Datas) > 0 {
// 		datainsert = repository.Insert{
// 			Datas: data.Save.Datas,
// 		}
// 	}

// 	return
// }

// func (mdl *RawModelParser) getFields(buffer *bytes.Buffer, Interfacefield interface{}) error {
// 	v := reflect.ValueOf(Interfacefield)
// 	for i := 0; i < v.Type().NumField(); i++ {
// 		buffer.WriteString(v.Type().Field(i).Tag.Get("json"))
// 		if i+1 < v.Type().NumField() {
// 			err := buffer.WriteByte(0b10101010)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

const spliter = byte(0b10101010)

// A function to extract the fields of a struct that have the tag restql with the value "enable" and return a slice of bytes with the column names and a separator 0b10101010
func (mdl *RawModelParser) getFields(buffer *bytes.Buffer, v interface{}) error {
	// Define a helper function to convert a string to snake case
	toSnakeCase := func(s string) string {
		// Split the string by uppercase letters
		parts := []string{}
		start := 0

		for i, c := range s {
			if i > 0 && c >= 'A' && c <= 'Z' {
				if s[i-1] >= 'A' && s[i-1] <= 'Z' {
					continue
				}
				parts = append(parts, s[start:i])
				start = i
			}
		}
		parts = append(parts, s[start:])

		// Join the parts with underscores and lowercase them
		return strings.ToLower(strings.Join(parts, "_"))
	}

	// Get the type of the struct
	t := reflect.TypeOf(v)

	// Define a helper function to get the column name from a slice of tag values or use snake case of the field name
	getColumnName := func(tagParts []string, fieldName string) (string, error) {
		count := 0       // A counter to track how many values have "indb:" prefix
		columnName := "" // A variable to store the column name

		for _, part := range tagParts {
			if strings.HasPrefix(part, "indb:") {
				count++
				if count > 1 {
					return "", errors.New("tag parts 'indb:' cannot have more than 1 value")
				}
				columnName = strings.TrimPrefix(part, "indb:")
			}
		}

		if count == 0 {
			columnName = toSnakeCase(fieldName)
		}

		return columnName, nil
	}

	// Loop through the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		// Get the field name and tag
		fieldName := t.Field(i).Name
		tag := t.Field(i).Tag.Get("restql")

		// Split the tag by semicolons
		tagParts := strings.Split(tag, ";")
		// Get the column name from the tag or use snake case of the field name
		columnName, err := getColumnName(tagParts, fieldName)
		if err != nil {
			return err
		}

		// Write the column name to the buffer
		buffer.WriteString(columnName)

		// Write the separator 0b10101010 to the buffer only if i+1 is less than t.NumField()
		if i+1 < t.NumField() {
			buffer.WriteByte(spliter)
		}

	}

	// Return the slice of bytes from the buffer
	return nil
}

func (mdl *RawModelParser) stringInSlice(bufferfield *bytes.Buffer, s string) bool {
	splits := bytes.Split(bufferfield.Bytes(), []byte{spliter})

	var same bool = false
	for _, split := range splits {
		if bytes.Equal(split, []byte(s)) {
			same = true
		}
	}

	return same
}

func (mdl *RawModelParser) modelpool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return rawmodelactions{}
		},
	}
}

// Query for single raw parser to Repository ORM
func (mdl *RawModelParser) Query(data []byte, model interface{}) (repo repository.Repository, err error) {
	if len(data) <= 0 {
		data = []byte("{}")
	}

	datamodel := rawmodelactions{}
	err = json.Unmarshal(data, &datamodel)
	if err != nil {
		return
	}

	mdl.query(datamodel, model)
	repo = mdl.txrepo
	mdl.txrepo = mdl.repo

	return
}
