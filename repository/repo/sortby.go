package repo

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

func (query *Repo) ValidateAndReturnSortQuery(sortBy string, Interfacefield interface{}) (string, error) {
	userFields := &bytes.Buffer{}
	userFields.Reset()

	err := query.getFields(userFields, Interfacefield)
	if err != nil {
		return "", err
	}
	sortBy = strings.ReplaceAll(sortBy, ", ", ",")
	commasplit := strings.Split(sortBy, ",")
	var orderby []string

	for _, cs := range commasplit {
		splits := strings.Split(cs, " ")
		if len(splits) != 2 {
			splits = append(splits, "asc")
		}
		field, order := splits[0], splits[1]
		order = strings.ToLower(order)

		if order != "desc" && order != "asc" {
			return "", errors.New("malformed orderdirection in sortBy query parameter, should be asc or desc")
		}

		if !query.stringInSlice(userFields, field) && field != "id" {
			return "", errors.New("unknown field in sortBy query parameter")
		}

		orderby = append(orderby, fmt.Sprintf("%s %s", field, strings.ToUpper(order)))
	}

	return strings.Join(orderby, ", "), nil
}
