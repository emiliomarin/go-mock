package doer

//go:generate mockgen --destination=../mocks/mock_doer.go --package=mocks . Doer

type Doer interface {
	Do() error
}
