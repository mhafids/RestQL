package repository

type ICommand struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type Insert struct {
	Datas []ICommand `json:"datas"`

	Select []string `json:"select"`
}

type Update struct {
	Datas   []ICommand `json:"datas"`
	Where   IFilter    `json:"where"`
	Orderby []ISortBy  `json:"orderby"`
	Limit   int        `json:"limit"`
	Skip    int        `json:"skip"`
	Select  []string   `json:"select"`
}

type Delete struct {
	Where   IFilter   `json:"where"`
	Orderby []ISortBy `json:"orderby"`
	Limit   int       `json:"limit"`
	Skip    int       `json:"skip"`
	Select  []string  `json:"select"`
}
