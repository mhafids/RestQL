package repo

import (
	"bytes"
	"errors"
)

func (query *Repo) selectFilterCheck(selects []string, model interface{}) error {
	userFields := &bytes.Buffer{}
	userFields.Reset()

	err := query.getFields(userFields, model)
	if err != nil {
		return err
	}
	for _, Select := range selects {
		if !query.stringInSlice(userFields, Select) {
			return errors.New(Select + " field not found")
		}
	}

	return nil
}
