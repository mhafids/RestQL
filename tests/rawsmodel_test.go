package tests

import (
	"encoding/json"
	"restql/parser"
	"restql/repository/repo"
	"testing"
)

func TestRawModelQueryOne(t *testing.T) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewRawModel(repoCfg)
	var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}`

	var operatorMap parser.ModelActions
	json.Unmarshal([]byte(operatorJSON), &operatorMap)
	operatorJSON = ""

	t.Log(operatorMap.Find)
	op, err := mts.QueryOne(operatorMap, Rawmodels{})
	if err != nil {
		t.Error(err)
	}

	orm, err := op.ToORM()
	if err != nil {
		t.Error(err)
	}
	t.Log(orm)
}

func TestRawModelQuery(t *testing.T) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewRawModel(repoCfg)
	var operatorJSON string = `{"test":{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}}`

	var operatorMap parser.ModelColumn
	json.Unmarshal([]byte(operatorJSON), &operatorMap)
	operatorJSON = ""

	var models map[string]interface{} = make(map[string]interface{}, 0)
	models["test"] = Rawmodels{}

	op, err := mts.Query(operatorMap, models)
	if err != nil {
		t.Error(err)
	}

	for k, v := range op {
		orm, err := v.ToORM()
		if err != nil {
			t.Error(err)
		}
		t.Log(k, orm)
	}
}
