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
	Last      int    `json:"last"`
	Oke       bool   `json:"oke"`
}

func main() {
	repocfg := repo.NewRepo(repo.RepoConfig{})
	var op repository.Repository
	var err error
	mts := parser.NewRawModel(repocfg)
	var operatorJSON string = `{"find":{"op":"$and","items":[{"op":"$eq","field":"last","value":501},{"op":"$eq","field":"oke","value":true},{"op":"$eq","field":"first_name","value":"Jawa Barat"},{"op":"$eq","field":"first_name","value":"Jawa Tengah"},{"op":"$eq","field":"first_name","value":"Banten"},{"op":"$eq","field":"first_name","value":"Yogyakarta"},{"op":"$eq","field":"first_name","value":"Jakarta"},{"op":"$eq","field":"first_name","value":"kalimantan Barat"},{"op":"$eq","field":"first_name","value":"sumatra Barat"}]}}`
	op, err = mts.Query([]byte(operatorJSON), Rawmodels{})
	if err != nil {
		fmt.Println(err)
	}

	logs, err := op.ToORM()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(logs)
}

// func Rawmodel(repocfg repository.Repository) {

// 	// b.Log(orm)
// }

// func Mongomodel(repocfg repository.Repository) {
// 	repoCfg := repo.NewRepo(&repo.RepoConfig{})
// 	mts := parser.NewMongoModel(repoCfg)
// 	var operatorJSON string = `{"find":{"phone":{"$not":{"$gt":"25"}}}}`

// 	op, err := mts.Query([]byte(operatorJSON), Rawmodels{})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	orm, err := op.ToORM()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("Mongo Model:", orm)
// }
