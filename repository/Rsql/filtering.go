package Rsql

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/mhafids/RestQL/repository"
)

func (query *Repo) Limit(data int) repository.Repository {
	query.mtx.Lock()
	defer query.mtx.Unlock()

	if query.model == nil {
		query.err = errors.New("Model not Nil")

		return query
	}

	if query.err != nil {
		return query
	}

	if data > query.Setting.MaxforLimit {
		data = query.Setting.MaxforLimit
	}

	query.data.Limit = data
	return query
}

func (query *Repo) Offset(data int) repository.Repository {
	query.mtx.Lock()
	defer query.mtx.Unlock()

	if query.model == nil {
		query.err = errors.New("Model not Nil")
		return query
	}

	if query.err != nil {
		return query
	}

	// if data > query.Setting.MaxforOffset {
	// 	data = query.Setting.MaxforOffset
	// }

	query.data.Offset = data
	return query
}

func (query *Repo) SortBy(sorts []repository.ISortBy) repository.Repository {
	query.mtx.Lock()
	defer query.mtx.Unlock()

	if query.model == nil {
		query.err = errors.New("Model not Nil")
		return query
	}

	if query.err != nil {
		return query
	}

	userFields := &bytes.Buffer{}
	userFields.Reset()

	err := query.getFields(userFields, query.model)
	if err != nil {
		query.err = err
		return query
	}

	fieldsort := func(commasplit []string, direction string) (string, string) {
		commasplit[1] = strings.ReplaceAll(commasplit[1], ":)))", ",")
		split := strings.Split(commasplit[1], ",")
		field := split[0]
		values := split[1:]

		orderby := "CASE "
		if direction == "asc" {
			for index, value := range values {
				orderby += "WHEN " + field + "='" + value + "' THEN " + strconv.Itoa((index + 1)) + " "
			}
		} else {
			for index := range values {
				orderby += "WHEN " + field + "='" + values[len(values)-(index+1)] + "' THEN " + strconv.Itoa((index + 1)) + " "
			}
		}
		orderby += "END"
		return orderby, field
	}

	for index, sort := range sorts {
		if strings.ToLower(sort.Sort) != "desc" && strings.ToLower(sort.Sort) != "asc" {
			query.err = errors.New("malformed order direction in sortBy query parameter, should be asc or desc")
			return query
		}

		regex := regexp.MustCompile(`(?i)field\(([^(]*)\)`)
		regexsubmatch := regex.FindStringSubmatch(sort.Field)
		if len(regexsubmatch) > 0 {
			result, field := fieldsort(regexsubmatch, strings.ToLower(sort.Sort))
			if !query.stringInSlice(userFields, field) && field != "id" {
				query.err = errors.New("unknown field in sortBy query parameter")
				return query
			}

			sorts[index].Field = result
			sorts[index].Sort = ""
		}
	}

	var ssort []string
	for _, sort := range sorts {
		if sort.Sort == "" {
			ssort = append(ssort, sort.Field)
		} else {
			ssort = append(ssort, sort.Field+" "+sort.Sort)
		}
	}

	query.data.SortBy = strings.Join(ssort, ", ")
	return query
}
