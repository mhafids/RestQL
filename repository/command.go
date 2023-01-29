package repository

type ICommand struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type Insert struct {
	Datas []ICommand `json:"datas"`
}

type Update struct {
	Datas []ICommand `json:"datas"`

	Find   IFilter `json:"find"`
	Filter IFilter `json:"filter"`
	Where  IFilter `json:"where"`

	Sortby  string `json:"sortby"`
	Sort    string `json:"sort"`
	Orderby string `json:"orderby"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`
	Skip   int `json:"skip"`
}

type Delete struct {
	Find    IFilter `json:"find"`
	Filter  IFilter `json:"filter"`
	Where   IFilter `json:"where"`
	Sortby  string  `json:"sortby"`
	Sort    string  `json:"sort"`
	Orderby string  `json:"orderby"`

	Limit int `json:"limit"`

	Offset int `json:"offset"`
	Skip   int `json:"skip"`
}
