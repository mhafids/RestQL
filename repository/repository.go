package repository

type IORM struct {
	Insert map[string]interface{}
	Update map[string]interface{}

	SortBy string
	Filter IFilterProcessed
	Limit  int
	Offset int
	Select []string
}

type Repository interface {
	QueryRepository
	CommandRepository

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
	Insert(data Insert, model interface{}) (err error)
	Update(data Update, model interface{}) (err error)
	Delete(data Delete, model interface{}) (err error)
}
