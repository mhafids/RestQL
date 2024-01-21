package Rsql

import (
	"bytes"
	"errors"
	"strings"

	"github.com/mhafids/RestQL/exception"
	"github.com/mhafids/RestQL/repository"
)

func (query *Repo) Select(data []string) repository.Repository {
	query.mtx.Lock()
	defer query.mtx.Unlock()

	if query.model == nil {
		query.err = errors.New("Model not Nil")
		return query
	}

	if query.err != nil {
		return query
	}

	buffer := &bytes.Buffer{}
	buffer.Reset()

	err := query.getFields(buffer, query.model)
	if err != nil {
		query.err = err
		return query
	}

	for _, dt := range data {
		if !query.stringInSlice(buffer, dt) {
			query.err = errors.New(dt + exception.FieldUnknown)
			return query
		}
	}

	query.data.Select = data
	return query
}

func (query *Repo) GroupBy(data []string) repository.Repository {
	query.mtx.Lock()
	defer query.mtx.Unlock()

	if query.model == nil {
		query.err = errors.New("Model not Nil")
		return query
	}

	if query.err != nil {
		return query
	}

	buffer := &bytes.Buffer{}
	buffer.Reset()

	err := query.getFields(buffer, query.model)
	if err != nil {
		query.err = err
		return query
	}

	for _, dt := range data {
		if !query.stringInSlice(buffer, dt) {
			query.err = errors.New(dt + exception.FieldUnknown)
			return query
		}
	}

	query.data.GroupBy = strings.Join(data, ", ")
	return query
}

func (query *Repo) Model(model interface{}) repository.Repository {
	query.wg.Add(1)
	query.mtx.Lock()
	defer query.mtx.Unlock()
	query.err = nil
	query.model = model

	err := query.initialselect()
	if err != nil {
		query.err = err
		return query
	}
	return query
}
