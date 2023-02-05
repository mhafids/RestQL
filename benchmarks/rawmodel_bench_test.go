package benchmarks

import (
	"testing"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository/repo"
)

func BenchmarkRawModelQueryOne(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	repoCfg := repo.NewRepo(repo.RepoConfig{})
	mts := parser.NewRawModel(repoCfg)
	var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}`
	op, err := mts.Query([]byte(operatorJSON), Rawmodels{})
	if err != nil {
		b.Error(err)
	}

	_, err = op.ToORM()
	if err != nil {
		b.Error(err)
	}
	// b.Log(orm)
}

func BenchmarkRawModelBenchQuery(b *testing.B) {
	repoCfg := repo.NewRepo(repo.RepoConfig{})
	mts := parser.NewRawModel(repoCfg)
	var operatorJSON string = `{"test":{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}}`

	var models map[string]interface{} = make(map[string]interface{}, 0)
	models["test"] = Rawmodels{}

	op, err := mts.QueryBatch([]byte(operatorJSON), models)
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
