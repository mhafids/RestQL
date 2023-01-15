package repository

import (
	"encoding/json"
	"errors"
)

type QueryJSON struct {
	Setting *QueryConfig
}

type QueryConfig struct {
	MaxQuery       int
	LimitDefault   int
	RequiredSelect bool
}

func NewQueryJSON(config *QueryConfig) Repository {
	return &QueryJSON{
		Setting: config,
	}
}

func (query *QueryJSON) flimit(listPayload ListPayload) (limit int, err error) {
	// Limit
	if query.Setting.LimitDefault <= 0 {
		limit = 10
	} else if listPayload.Limit > 0 {
		limit = listPayload.Limit
	} else {
		limit = query.Setting.LimitDefault
	}

	return
}

func (query *QueryJSON) foffset(listPayload ListPayload) (offset int, err error) {
	strOffset := listPayload.Offset
	strSkip := listPayload.Skip

	offset = 0
	if strOffset > 0 {
		offset = strOffset
	} else if strSkip > 0 {
		offset = strSkip
	}

	return
}

func (query *QueryJSON) List(listPayload ListPayload, model interface{}) (restapiFilter IFilterSearch, err error) {

	restapiFilter.SortBy, err = query.sortBy(listPayload, model)
	if err != nil {
		return
	}

	restapiFilter.Limit, err = query.limit(listPayload)
	if err != nil {
		return
	}

	restapiFilter.Offset, err = query.offset(listPayload)
	if err != nil {
		return
	}

	restapiFilter.Filter, err = query.Filter(listPayload, model)
	if err != nil {
		return
	}

	restapiFilter.Select, err = query.selects(listPayload, model)
	if err != nil {
		return
	}

	return
}

func (query *QueryJSON) SortBy(listPayload ListPayload, model interface{}) (Query string, err error) {
	// Sort
	request := listPayload
	sortBy := ""

	if request.Sort != "" {
		sortBy = request.Sort
	}

	if request.Orderby != "" {
		sortBy = request.Orderby
	}

	if sortBy == "" {
		// id.asc is the default sort query
		sortBy = "id asc"
	}

	Query, err = ValidateAndReturnSortQuery(sortBy, model)
	return
}

func (query *QueryJSON) Limit(listPayload ListPayload) (limit int, err error) {
	offset, err := query.foffset(listPayload)
	if err != nil {
		return
	}

	limit, err = query.flimit(listPayload)
	if err != nil {
		return
	}

	var MaxQuery = 0

	if query.Setting.MaxQuery > 0 {
		MaxQuery = query.Setting.MaxQuery
	}

	if MaxQuery > 0 {
		if offset >= (MaxQuery) {
			limit = offset - MaxQuery
			return
		}
	}

	return
}

func (query *QueryJSON) Offset(listPayload ListPayload) (offset int, err error) {
	offset, err = query.foffset(listPayload)
	if err != nil {
		return
	}

	limits, err := query.flimit(listPayload)
	if err != nil {
		return
	}

	if offset >= (query.Setting.MaxQuery) {
		offset = query.Setting.MaxQuery - limits
	}

	return
}

func (query *QueryJSON) Filter(listPayload ListPayload, model interface{}) (filterProcessed IFilterProcessed, err error) {

	filter := listPayload.Filter
	where := listPayload.Where
	find := listPayload.Find

	if where.Operator != "" {
		filter = where
	} else if find.Operator != "" {
		filter = find
	}

	filterProcessed, err = query.filterDB(filter, model)

	if err != nil {
		return
	}
	return
}

func (query *QueryJSON) Select(listPayload ListPayload, model interface{}) (selects []string, err error) {
	request := listPayload

	if len(request.Select) == 0 && query.Setting.RequiredSelect {
		err = errors.New("Select request Not null")
		return
	}

	err = query.selectFilterCheck(request.Select, model)
	if err != nil {
		return
	}
	selects = request.Select

	return
}

func convertMapToStruct(datamap interface{}, marshal interface{}) (err error) {
	// convert map to json
	jsonString, err := json.Marshal(datamap)
	if err != nil {
		return
	}
	// convert json to struct
	err = json.Unmarshal(jsonString, marshal)
	if err != nil {
		return
	}
	return
}
