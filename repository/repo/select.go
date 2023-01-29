package repo

import (
	"bytes"
	"errors"

	"github.com/mhafids/RestQL/utils"
)

func (query *Repo) selectFilterCheck(selects []string, model interface{}) error {
	userFields := utils.Bufpool.Get().(*bytes.Buffer)
	userFields.Reset()
	defer utils.Bufpool.Put(userFields)

	getFields(userFields, model)
	for _, Select := range selects {
		if !stringInSlice(userFields, Select) {
			return errors.New(Select + " field not found")
		}
	}

	return nil
}
