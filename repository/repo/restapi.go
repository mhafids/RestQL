package repo

import (
	"RestQL/repository"
	"encoding/json"
	"strings"
)

type Repo struct {
	Setting *RepoConfig

	// Filter processed
	data repository.IORM
}

type RepoConfig struct {
}

func NewRepo(config *RepoConfig) repository.Repository {
	return &Repo{
		Setting: config,
		data:    repository.IORM{},
	}
}

func (query *Repo) Filter(data repository.IFilter, model interface{}) (err error) {
	filtered, err := query.filterDB(data, model)

	if err != nil {
		return
	}

	query.data.Filter = filtered
	return
}

func (query *Repo) Limit(data int) (err error) {
	query.data.Limit = data
	return
}

func (query *Repo) Offset(data int) (err error) {
	query.data.Offset = data
	return
}

func (query *Repo) SortBy(sorts []repository.ISortBy, model interface{}) (err error) {
	var ssort []string
	for _, sort := range sorts {
		strings.Join(ssort, sort.Field+sort.Sort)
	}

	query.data.SortBy = strings.Join(ssort, ", ")
	return
}

func (query *Repo) Select(data []string, model interface{}) (err error) {
	query.data.Select = data
	return
}

func (query *Repo) ToORM() (orm repository.IORM, err error) {
	orm = query.data
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
