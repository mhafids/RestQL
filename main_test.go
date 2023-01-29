package main

import (
	"testing"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository/repo"
)

func BenchmarkRawModelQueryOne(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			repoCfg := repo.NewRepo(&repo.RepoConfig{})
			mts := parser.NewRawModel(repoCfg)
			var operatorJSON string = `{"find":{"op":"$and","items":[{"op":"$eq","field":"first_name","value":"Jawa Timur"},{"op":"$eq","field":"first_name","value":"Jawa Barat"},{"op":"$eq","field":"first_name","value":"Jawa Tengah"},{"op":"$eq","field":"first_name","value":"Banten"},{"op":"$eq","field":"first_name","value":"Yogyakarta"},{"op":"$eq","field":"first_name","value":"Jakarta"},{"op":"$eq","field":"first_name","value":"kalimantan Barat"},{"op":"$eq","field":"first_name","value":"sumatra Barat"}]}}`
			op, err := mts.Query(operatorJSON, Rawmodels{})
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
