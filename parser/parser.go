package parser

import (
	"github.com/mhafids/RestQL/repository"
)

type Parser interface {
	Query(data string, model interface{}) (repo repository.Repository, err error)
	QueryBatch(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error)
	Command(data string, model map[string]interface{}) (repo repository.Repository, err error)
	CommandBatch(data string, model map[string]interface{}) (repo map[string]repository.Repository, err error)
}
