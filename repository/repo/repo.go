package repo

import "github.com/mhafids/RestQL/repository"

type Repo struct {
	Setting RepoConfig

	// Filter processed
	data repository.IORM
}

type RepoConfig struct {
}

func NewRepo(config RepoConfig) repository.Repository {
	return &Repo{
		Setting: config,
		data:    repository.IORM{},
	}
}
