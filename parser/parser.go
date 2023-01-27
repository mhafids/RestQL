package parser

import "github.com/mhafids/RestQL/repository"

type Parser interface {
	QueryOne(data string, model interface{}) (repo repository.Repository, err error)
	Query(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error)
	Command(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error)
}

