package exception

const (
	Orderby        string = "malformed order direction in sortBy query parameter, should be asc or desc"
	OrderbyUnknown string = "unknown field in sortBy query parameter"
	FieldUnknown   string = " field not found"
)
