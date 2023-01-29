package repository

type IORM struct {
	SortBy string
	Filter IFilterProcessed
	Limit  int
	Offset int
	Select []string
}

type Repository interface {
	QueryRepository

	OutputRepository
}

type OutputRepository interface {
	ToORM() (orm IORM, err error)
}

// QueryRepository is Interface for QueryBatch
type QueryRepository interface {
	Filter(data IFilter, model interface{}) (err error)
	Limit(data int) (err error)
	Offset(data int) (err error)
	SortBy(sorts []ISortBy, model interface{}) (err error)
	Select(data []string, model interface{}) (err error)
}

type CommandRepository interface {
	// Filter(data IFilter, model interface{}) (err error)
	// Limit(data int64) (err error)
	// Offset(data int64) (err error)
	// SortBy(data []ISortBy, model interface{}) (err error)
	// Select(data []string, model interface{}) (err error)
}
