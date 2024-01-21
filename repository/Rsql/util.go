package Rsql

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
)

func (query *Repo) initialselect() error {
	var data []string
	buffer := &bytes.Buffer{}
	buffer.Reset()
	err := query.getFields(buffer, query.model)
	if err != nil {
		return err
	}

	databytes := bytes.Split(buffer.Bytes(), []byte{0b10101010})

	for _, databyte := range databytes {
		data = append(data, string(databyte))
	}

	query.data.Select = data
	return nil
}

const spliter = byte(0b10101010)

// A function to extract the fields of a struct that have the tag restql with the value "enable" and return a slice of bytes with the column names and a separator 0b10101010
func (query *Repo) getFields(buffer *bytes.Buffer, v interface{}) error {
	// Define a helper function to convert a string to snake case
	toSnakeCase := func(s string) string {
		// Split the string by uppercase letters
		parts := []string{}
		start := 0

		for i, c := range s {
			if i > 0 && c >= 'A' && c <= 'Z' {
				if s[i-1] >= 'A' && s[i-1] <= 'Z' {
					continue
				}
				parts = append(parts, s[start:i])
				start = i
			}
		}
		parts = append(parts, s[start:])

		// Join the parts with underscores and lowercase them
		return strings.ToLower(strings.Join(parts, "_"))
	}

	// Get the type of the struct
	t := reflect.TypeOf(v)

	// Define a helper function to get the column name from a slice of tag values or use snake case of the field name
	getColumnName := func(tagParts []string, fieldName string) (string, error) {
		count := 0       // A counter to track how many values have "indb:" prefix
		columnName := "" // A variable to store the column name

		for _, part := range tagParts {
			if strings.HasPrefix(part, "indb:") {
				count++
				if count > 1 {
					return "", errors.New("tag parts 'indb:' cannot have more than 1 value")
				}
				columnName = strings.TrimPrefix(part, "indb:")
			}
		}

		if count == 0 {
			columnName = toSnakeCase(fieldName)
		}

		return columnName, nil
	}

	// Loop through the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		// Get the field name and tag
		fieldName := t.Field(i).Name
		tag := t.Field(i).Tag.Get("restql")

		// Split the tag by semicolons
		tagParts := strings.Split(tag, ";")
		// Get the column name from the tag or use snake case of the field name
		columnName, err := getColumnName(tagParts, fieldName)
		if err != nil {
			return err
		}

		// Write the column name to the buffer
		buffer.WriteString(columnName)

		// Write the separator 0b10101010 to the buffer only if i+1 is less than t.NumField()
		if i+1 < t.NumField() {
			buffer.WriteByte(spliter)
		}

	}

	// Return the slice of bytes from the buffer
	return nil
}

func (query *Repo) stringInSlice(bufferfield *bytes.Buffer, s string) bool {
	splits := bytes.Split(bufferfield.Bytes(), []byte{spliter})

	var same bool = false
	for _, split := range splits {
		if bytes.Equal(split, []byte(s)) {
			same = true
		}
	}
	return same
}
