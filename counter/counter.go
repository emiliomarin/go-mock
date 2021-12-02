package counter

//go:generate mockgen --destination=../mocks/mock_counter.go --package=mocks . Counter,Doer

type Counter interface {
	Count(s string) (int, error)
}

type Doer interface {
	Do() error
}
