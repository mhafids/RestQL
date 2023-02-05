package repo

import (
	"bytes"
	"errors"

	"github.com/mhafids/RestQL/exception"
	"github.com/mhafids/RestQL/repository"
)

// Insert use Insert data to repository
func (command *Repo) Insert(data repository.Insert, model interface{}) (err error) {
	result := make(map[string]interface{}, 0)
	err = command.Select(data.Select, model)
	if err != nil {
		return
	}

	buffer := &bytes.Buffer{}
	buffer.Reset()

	err = command.getFields(buffer, model)
	if err != nil {
		return
	}

	for _, datainsert := range data.Datas {
		if command.stringInSlice(buffer, datainsert.Field) {
			result[datainsert.Field] = datainsert.Value
		} else {
			err = errors.New(datainsert.Field + exception.FieldUnknown)
		}
	}

	command.data.Insert = result
	return
}

// Update use Update data to repository
func (command *Repo) Update(data repository.Update, model interface{}) (err error) {
	result := make(map[string]interface{}, 0)

	// Select
	err = command.Select(data.Select, model)
	if err != nil {
		return
	}

	// orderby
	err = command.SortBy(data.Orderby, model)
	if err != nil {
		return err
	}

	// Where
	err = command.Filter(data.Where, model)
	if err != nil {
		return err
	}

	// Limit
	err = command.Limit(data.Limit)
	if err != nil {
		return err
	}

	// Skip
	err = command.Offset(data.Skip)
	if err != nil {
		return err
	}

	buffer := &bytes.Buffer{}
	buffer.Reset()

	err = command.getFields(buffer, model)
	if err != nil {
		return
	}

	for _, datainsert := range data.Datas {
		if command.stringInSlice(buffer, datainsert.Field) {
			result[datainsert.Field] = datainsert.Value
		} else {
			err = errors.New(datainsert.Field + exception.FieldUnknown)
		}
	}

	command.data.Update = result

	return
}

// Delete use Delete data to repository
func (command *Repo) Delete(data repository.Delete, model interface{}) (err error) {

	// Select
	err = command.Select(data.Select, model)
	if err != nil {
		return
	}

	// orderby
	err = command.SortBy(data.Orderby, model)
	if err != nil {
		return err
	}

	// Where
	err = command.Filter(data.Where, model)
	if err != nil {
		return err
	}

	// Limit
	err = command.Limit(data.Limit)
	if err != nil {
		return err
	}

	// Skip
	err = command.Offset(data.Skip)
	if err != nil {
		return err
	}

	return
}
