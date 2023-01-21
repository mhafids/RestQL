package benchmarks

import (
	"encoding/json"
	"testing"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository/repo"
)

func BenchmarkMongoModelQueryOne(b *testing.B) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewMongoModel(repoCfg)
	var operatorJSON string = `{"find":{"phone":{"$not":{"$gt":"25"}}}}`

	var operatorMap parser.ModelActions
	json.Unmarshal([]byte(operatorJSON), &operatorMap)
	operatorJSON = ""

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
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

func BenchmarkMongoModelBenchQuery(b *testing.B) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewMongoModel(repoCfg)
	var operatorJSON string = `{"test":{"find":{"phone":{"$not":{"$gt":"25"}}}}}`

	var operatorMap parser.ModelColumn
	json.Unmarshal([]byte(operatorJSON), &operatorMap)
	operatorJSON = ""

	var models map[string]interface{} = make(map[string]interface{}, 0)
	models["test"] = Rawmodels{}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			op, err := mts.Query(operatorMap, models)
			if err != nil {
				b.Error(err)
			}

			for _, v := range op {
				_, err := v.ToORM()
				if err != nil {
					b.Error(err)
				}
				// b.Log(k, orm)
			}
		}
	})
}
