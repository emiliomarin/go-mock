package counter

//go:generate mockgen --destination=../mocks/mock_counter.go --package=mocks . Counter

type Counter interface {
	Count(s string) (int, error)
}
