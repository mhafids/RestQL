package benchmarks

import (
	"encoding/json"
	"restql/parser"
	"restql/repository/repo"
	"testing"
)

func BenchmarkRawModelQueryOne(b *testing.B) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewRawModel(repoCfg)
	var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}`

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

func BenchmarkRawModelBenchQuery(b *testing.B) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewRawModel(repoCfg)
	var operatorJSON string = `{"test":{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}}`

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
