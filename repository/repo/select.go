package repo

import "errors"

func (query *Repo) selectFilterCheck(selects []string, model interface{}) error {
	var userFields = getFields(model)
	for _, Select := range selects {
		if !stringInSlice(userFields, Select) {
			return errors.New(Select + " field not found")
		}
	}

	return nil
}
