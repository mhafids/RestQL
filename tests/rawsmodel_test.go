package tests

import (
	"testing"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository/repo"
)

func TestRawModelQueryOne(t *testing.T) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewRawModel(repoCfg)
	var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}`

	op, err := mts.QueryOne(operatorJSON, Rawmodels{})
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

	var models map[string]interface{} = make(map[string]interface{}, 0)
	models["test"] = Rawmodels{}

	op, err := mts.Query(operatorJSON, models)
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
