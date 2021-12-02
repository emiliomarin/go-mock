package user

import "github.com/emiliomarin/go-mock/counter"

type User struct {
	Counter counter.Counter
}

func (u User) Count(s string) (int, error) {
	return u.Counter.Count(s)
}
