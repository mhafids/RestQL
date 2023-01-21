package repo

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateAndReturnSortQuery(sortBy string, Interfacefield interface{}) (string, error) {
	var userFields = getFields(Interfacefield)
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

		if !stringInSlice(userFields, field) && field != "id" {
			return "", errors.New("unknown field in sortBy query parameter")
		}

		orderby = append(orderby, fmt.Sprintf("%s %s", field, strings.ToUpper(order)))
	}

	return strings.Join(orderby, ", "), nil
}

func stringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}
