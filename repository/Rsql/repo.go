package Rsql

import (
	"sync"

	"github.com/mhafids/RestQL/repository"
)

type Repo struct {
	Setting RepoConfig

	// Filter processed
	data  repository.IORM
	model interface{}
	mtx   sync.Mutex
	wg    sync.WaitGroup
	err   error
}

type RepoConfig struct {
	MaxforLimit  int
	MinforLimit  int
	MaxforOffset int
}

func NewRepoSql(config RepoConfig) repository.Repository {
	if config.MaxforLimit <= 0 {
		config.MaxforLimit = 700
	}

	if config.MinforLimit <= 0 {
		config.MinforLimit = 10
	}

	if config.MaxforOffset <= 0 {
		config.MaxforOffset = 200
	}

	repo := &Repo{
		Setting: config,
		data:    repository.IORM{Limit: config.MinforLimit},
	}

	return repo
}
