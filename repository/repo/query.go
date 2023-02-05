package repo

import (
	"bytes"
	"errors"
	"reflect"
	"strings"

	"github.com/mhafids/RestQL/exception"
	"github.com/mhafids/RestQL/repository"
)

func (query *Repo) Filter(data repository.IFilter, model interface{}) (err error) {
	filtered, err := query.filterDB(data, model)

	if err != nil {
		return
	}

	query.data.Filter = filtered
	return
}

func (query *Repo) Limit(data int) (err error) {
	query.data.Limit = data
	return
}

func (query *Repo) Offset(data int) (err error) {
	query.data.Offset = data
	return
}

func (query *Repo) SortBy(sorts []repository.ISortBy, model interface{}) (err error) {
	var ssort []string
	for _, sort := range sorts {
		strings.Join(ssort, sort.Field+sort.Sort)
	}

	query.data.SortBy = strings.Join(ssort, ", ")
	return
}

func (query *Repo) Select(data []string, model interface{}) (err error) {
	buffer := &bytes.Buffer{}
	buffer.Reset()

	err = query.getFields(buffer, model)
	if err != nil {
		return err
	}

	for _, dt := range data {
		if !query.stringInSlice(buffer, dt) {
			err = errors.New(dt + exception.FieldUnknown)
			return
		}
	}

	query.data.Select = data
	return
}

func (query *Repo) ToORM() (orm repository.IORM, err error) {
	orm = query.data
	return
}

func (query *Repo) getFields(buffer *bytes.Buffer, Interfacefield interface{}) error {
	v := reflect.ValueOf(Interfacefield)
	for i := 0; i < v.Type().NumField(); i++ {
		buffer.WriteString(v.Type().Field(i).Tag.Get("json"))
		if i+1 < v.Type().NumField() {
			err := buffer.WriteByte(0b10101010)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (query *Repo) stringInSlice(bufferfield *bytes.Buffer, s string) bool {
	return bytes.Contains(bufferfield.Bytes(), []byte(s))
}
