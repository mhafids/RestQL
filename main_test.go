package main

import (
	"testing"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository"
	"github.com/mhafids/RestQL/repository/repo"
)

func BenchmarkRawModelQueryOne(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var op repository.Repository
		var err error
		repoCfg := repo.NewRepo(repo.RepoConfig{})
		mts := parser.NewRawModel(repoCfg)
		var operatorJSON string = `{"find":{"op":"$and","items":[{"op":"$eq","field":"last","value":501},{"op":"$eq","field":"first_name","value":"Jawa Barat"},{"op":"$eq","field":"first_name","value":"Jawa Tengah"},{"op":"$eq","field":"first_name","value":"Banten"},{"op":"$eq","field":"first_name","value":"Yogyakarta"},{"op":"$eq","field":"first_name","value":"Jakarta"},{"op":"$eq","field":"first_name","value":"kalimantan Barat"},{"op":"$eq","field":"first_name","value":"sumatra Barat"}]}}`

		op, err = mts.Query([]byte(operatorJSON), Rawmodels{})
		if err != nil {
			b.Error(err)
		}

		_, err = op.ToORM()
		if err != nil {
			b.Error(err)
		}
		// b.Log(orm)
	}
}
