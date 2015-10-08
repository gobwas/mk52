package parser

type Expression interface {
	Evaluate() float64
	AddFunc(string, func(...interface{}) float64)
}

type Parser interface {
	Parse() (*Expression, error)
}
