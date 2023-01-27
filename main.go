package main

import (
	"fmt"

	"github.com/mhafids/RestQL/parser"
	"github.com/mhafids/RestQL/repository"
	"github.com/mhafids/RestQL/repository/repo"
)

type Rawmodels struct {
	By        string `json:"by"`
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

func main() {
	repocfg := repo.NewRepo(&repo.RepoConfig{})
	Rawmodel(repocfg)
	Mongomodel(repocfg)
}

func Rawmodel(repocfg repository.Repository) {
	mts := parser.NewRawModel(repocfg)
	var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}`

	op, err := mts.QueryOne(operatorJSON, Rawmodels{})
	if err != nil {
		fmt.Println(err)
	}

	orm, err := op.ToORM()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Rawmodel:", orm)
}

func Mongomodel(repocfg repository.Repository) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewMongoModel(repoCfg)
	var operatorJSON string = `{"find":{"phone":{"$not":{"$gt":"25"}}}}`

	op, err := mts.QueryOne(operatorJSON, Rawmodels{})
	if err != nil {
		fmt.Println(err)
	}

	orm, err := op.ToORM()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Mongo Model:", orm)
}
