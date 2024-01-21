package Rsql

import (
	"errors"

	"github.com/mhafids/RestQL/repository"
)

func (query *Repo) ToORM() (repository.IORM, error) {

	query.mtx.Lock()
	defer query.mtx.Unlock()

	defer query.wg.Done()

	if query.model == nil {
		query.err = errors.New("Model not Nil")
		query.model = nil
		return repository.IORM{}, query.err
	}

	if query.err != nil {
		query.model = nil
		return repository.IORM{}, query.err
	}

	orm := query.data
	query.data = repository.IORM{
		Limit: query.Setting.MinforLimit,
	}

	query.initialselect()
	query.model = nil
	return orm, query.err
}
