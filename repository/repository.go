package repository

type CommandSyntax string

const (
	// InsertSyntax statement is used for insert
	InsertSyntax CommandSyntax = "insert"

	// UpdateSyntax statement is used for update
	UpdateSyntax CommandSyntax = "update"

	// DeleteSyntax statement is used for delete
	DeleteSyntax CommandSyntax = "delete"
)

type IORM struct {
	Datas   map[string]interface{}
	Command CommandSyntax

	SortBy  string
	Filter  IFilterProcessed
	Limit   int
	Offset  int
	Select  []string
	GroupBy string
}

type Repository interface {
	Model(model interface{}) Repository

	QueryRepository
	CommandRepository
	OutputRepository
	FilteringRepository
}

type OutputRepository interface {
	ToORM() (IORM, error)
}

// QueryRepository is Interface for QueryBatch
type QueryRepository interface {
	Select(data []string) Repository
}

// FilteringRepository is interface for filter repository
type FilteringRepository interface {
	Filter(data IFilter) Repository
	Limit(data int) Repository
	Offset(data int) Repository
	SortBy(sorts []ISortBy) Repository
	GroupBy(data []string) Repository
}

type CommandRepository interface {
	Insert(Datas []ICommand) Repository
	Update(Datas []ICommand) Repository
	Delete() Repository
}
