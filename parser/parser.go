package parser

import (
	"github.com/mhafids/RestQL/repository"
)

type Parser interface {
	Query(data []byte, model interface{}) (repo repository.Repository, err error)
	Insert(data []byte, model interface{}) (repo repository.Repository, err error)
	Delete(data []byte, model interface{}) (repo repository.Repository, err error)
	Update(data []byte, model interface{}) (repo repository.Repository, err error)
	// QueryBatch(data []byte, model map[string]interface{}) (repo map[string]repository.Repository, err error)
	// Command(data []byte, model interface{}) (repo repository.Repository, err error)
	// CommandBatch(data []byte, model map[string]interface{}) (repo map[string]repository.Repository, err error)
}
