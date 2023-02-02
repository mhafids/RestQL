package main_test

import (
	"encoding/json"
	"testing"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository/repo"
)

func BenchmarkRawModelQueryOne(b *testing.B) {

	type Rawmodels struct {
		By        string `json:"by"`
		Title     string `json:"title"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			repoCfg := repo.NewRepo(&repo.RepoConfig{})
			mts := parser.NewRawModel(repoCfg)
			var operatorJSON string = `{"find":{"op":"$and","items":[{"op":"$eq","field":"first_name","value":"Jawa Timur"},{"op":"$eq","field":"first_name","value":"Jawa Barat"},{"op":"$eq","field":"first_name","value":"Jawa Tengah"},{"op":"$eq","field":"first_name","value":"Banten"},{"op":"$eq","field":"first_name","value":"Yogyakarta"},{"op":"$eq","field":"first_name","value":"Jakarta"},{"op":"$eq","field":"first_name","value":"kalimantan Barat"},{"op":"$eq","field":"first_name","value":"sumatra Barat"}]}}`

			var operatorMap parser.ModelActions
			json.Unmarshal([]byte(operatorJSON), &operatorMap)
			operatorJSON = ""

			op, err := mts.QueryOne(operatorMap, Rawmodels{})
			if err != nil {
				b.Error(err)
			}

			_, err = op.ToORM()
			if err != nil {
				b.Error(err)
			}
			// b.Log(orm)
		}
	})
}
