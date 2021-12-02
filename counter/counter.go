package counter

type Counter interface {
	Count(s string) (int, error)
}
