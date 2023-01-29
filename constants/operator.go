package constants

type Operator string

const (
	AND Operator = "$and"
	OR  Operator = "$or"
	NOT Operator = "$not"
	NOR Operator = "$nor"

	EQ    Operator = "$eq"
	NE    Operator = "$ne"
	LIKE  Operator = "$like"
	ILIKE Operator = "$ilike"
	GT    Operator = "$gt"
	GTE   Operator = "$gte"
	LT    Operator = "$lt"
	LTE   Operator = "$lte"
	NIN   Operator = "$nin"
	IN    Operator = "$in"
)
