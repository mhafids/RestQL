package queryrestapi

import "errors"

func (query *QueryJSON) selectFilterCheck(selects []string, model interface{}) error {
	var userFields = getFields(model)
	for _, Select := range selects {
		if !stringInSlice(userFields, Select) {
			return errors.New(Select + " field not found")
		}
	}

	return nil
}
