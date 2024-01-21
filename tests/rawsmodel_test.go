package tests

import (
	"testing"

	"github.com/mhafids/RestQL/parser/rawparser"
	"github.com/mhafids/RestQL/repository/Rsql"
)

func TestRawModelQueryOne(t *testing.T) {
	repoCfg := Rsql.NewRepoSql(Rsql.RepoConfig{})
	mts := rawparser.NewRawModelParser(repoCfg)
	var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}, "select":["by"]}`

	op, err := mts.Query([]byte(operatorJSON), Rawmodels{})
	if err != nil {
		t.Error(err)
	}

	orm, err := op.ToORM()
	if err != nil {
		t.Error(err)
	}
	t.Log(orm)
}
