package benchmarks

import (
	"testing"

	"github.com/mhafids/RestQL/parser/rawparser"
	"github.com/mhafids/RestQL/repository/Rsql"
)

func BenchmarkRawModelQueryOne(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	repoCfg := Rsql.NewRepoSql(Rsql.RepoConfig{})
	for i := 0; i < b.N; i++ {
		mts := rawparser.NewRawModelParser(repoCfg)
		var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}`
		op, err := mts.Query([]byte(operatorJSON), Rawmodels{})
		if err != nil {
			b.Error(err)
		}

		_, err = op.ToORM()
		if err != nil {
			b.Error(err)
		}
	}

	// b.Log(orm)
}
