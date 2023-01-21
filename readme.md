# RestQL
Query parser for `REST` to be easy building and connect with ORM.

# Instalation
```
go get -u github.com/mhafids/RestQL
```
# Get Started
## Raw Parser
```
package main

import (
	"encoding/json"
	"fmt"
	"restql/parser"
	"restql/repository"
	"restql/repository/repo"
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
}

func Rawmodel(repocfg repository.Repository) {
	mts := parser.NewRawModel(repocfg)
	var operatorJSON string = `{"find":{"op":"$eq","field":"first_name","value":"Jawa Timur"}}`

	var operatorMap parser.ModelActions
	json.Unmarshal([]byte(operatorJSON), &operatorMap)
	operatorJSON = ""

	op, err := mts.QueryOne(operatorMap, Rawmodels{})
	if err != nil {
		fmt.Println(err)
	}

	orm, err := op.ToORM()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Rawmodel:", orm)
}

```

## Mongo Parser
```
package main

import (
	"encoding/json"
	"fmt"
	"restql/parser"
	"restql/repository"
	"restql/repository/repo"
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
	Mongomodel(repocfg)
}

func Mongomodel(repocfg repository.Repository) {
	repoCfg := repo.NewRepo(&repo.RepoConfig{})
	mts := parser.NewMongoModel(repoCfg)
	var operatorJSON string = `{"find":{"phone":{"$not":{"$gt":"25"}}}}`

	var operatorMap parser.ModelActions
	json.Unmarshal([]byte(operatorJSON), &operatorMap)
	operatorJSON = ""

	op, err := mts.QueryOne(operatorMap, Rawmodels{})
	if err != nil {
		fmt.Println(err)
	}

	orm, err := op.ToORM()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Mongo Model:", orm)
}

```
