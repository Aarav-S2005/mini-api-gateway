package lb

type Backend interface {
	Identifier() string
	GetWeight() int
}
