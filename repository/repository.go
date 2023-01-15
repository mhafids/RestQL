package repository

type IORM struct {
	SortBy string
	Filter IFilterProcessed
	Limit  int
	Offset int
	Select []string
}

type IFilter struct {
	Operator string      `json:"op"`
	Field    string      `json:"field"`
	Items    []IFilter   `json:"items"`
	Value    interface{} `json:"value"`
}

type IFilterProcessed struct {
	Field  string
	Values []interface{}
}

type ISortBy struct {
	Field string
	Sort  string
}

type Repository interface {
	Filter(data IFilter, model interface{}) (err error)
	Limit(data int64) (err error)
	Offset(data int64) (err error)
	SortBy(data []ISortBy, model interface{}) (err error)
	Select(data []string, model interface{}) (err error)

	OutputRepository
}

type OutputRepository interface {
	ToORM() (orm IORM, err error)
}
