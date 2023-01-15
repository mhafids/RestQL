package parser

type Parser interface {
	Query()
	Insert()
	Update()
	Delete()
}
