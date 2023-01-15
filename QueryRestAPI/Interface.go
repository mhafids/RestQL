package queryrestapi

const (
	AND string = "$and"
	OR  string = "$or"
	NOT string = "$not"
	NOR string = "$nor"

	EQ    string = "$eq"
	NE    string = "$ne"
	LIKE  string = "$like"
	ILIKE string = "$ilike"
	GT    string = "$gt"
	GTE   string = "$gte"
	LT    string = "$lt"
	LTE   string = "$lte"
	NIN   string = "$nin"
	IN    string = "$in"
)

type IFilterSearch struct {
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
