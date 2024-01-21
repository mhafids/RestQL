package Rsql

import (
	"bytes"
	"errors"

	"github.com/mhafids/RestQL/exception"
	"github.com/mhafids/RestQL/repository"
)

// Insert use Insert data to repository
func (command *Repo) Insert(Datas []repository.ICommand) repository.Repository {
	result := make(map[string]interface{}, 0)

	buffer := &bytes.Buffer{}
	buffer.Reset()

	err := command.getFields(buffer, command.model)
	if err != nil {
		command.err = err
		command.model = nil
		return command
	}

	for _, datainsert := range Datas {
		if command.stringInSlice(buffer, datainsert.Field) {
			result[datainsert.Field] = datainsert.Value
		} else {
			err = errors.New(datainsert.Field + exception.FieldUnknown)
			if err != nil {
				command.err = err
				return command
			}
		}
	}

	command.data.Command = repository.InsertSyntax
	command.data.Datas = result
	return command
}

// Update use Update data to repository
func (command *Repo) Update(Datas []repository.ICommand) repository.Repository {
	result := make(map[string]interface{}, 0)

	buffer := &bytes.Buffer{}
	buffer.Reset()

	err := command.getFields(buffer, command.model)
	if err != nil {
		command.err = err
		command.model = nil
		return command
	}

	for _, datainsert := range Datas {
		if command.stringInSlice(buffer, datainsert.Field) {
			result[datainsert.Field] = datainsert.Value
		} else {
			err = errors.New(datainsert.Field + exception.FieldUnknown)
			if err != nil {
				command.err = err
				return command
			}
		}
	}

	command.data.Command = repository.UpdateSyntax
	command.data.Datas = result
	return command
}

// Delete use Delete data to repository
func (command *Repo) Delete() repository.Repository {
	command.data.Command = repository.InsertSyntax
	return command
}
