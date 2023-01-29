package repo

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/mhafids/RestQL/utils"
)

func ValidateAndReturnSortQuery(sortBy string, Interfacefield interface{}) (string, error) {
	userFields := utils.Bufpool.Get().(*bytes.Buffer)
	userFields.Reset()
	defer utils.Bufpool.Put(userFields)

	getFields(userFields, Interfacefield)
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

func stringInSlice(bufferfield *bytes.Buffer, s string) bool {
	return bytes.Contains(bufferfield.Bytes(), []byte(s))
}
