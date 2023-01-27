package tests

import (
	"testing"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository/repo"
)

func TestMongoModelQueryOne(t *testing.T) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewMongoModel(repoCfg)
	var operatorJSON string = `{"find":{"phone":{"$not":{"$gt":"25"}}}}`

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

func TestMongoModelQuery(t *testing.T) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewMongoModel(repoCfg)
	var operatorJSON string = `{"test":{"find":{"phone":{"$not":{"$gt":"25"}}}}}`

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
