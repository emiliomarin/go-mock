package user

import (
	"github.com/emiliomarin/go-mock/counter"
	"github.com/emiliomarin/go-mock/doer"
)

type User struct {
	Counter counter.Counter
	Doer    doer.Doer
}

func (u User) CountAndDo(s string) error {
	count, err := u.Counter.Count(s)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err = u.Doer.Do()
		if err != nil {
			return err
		}
	}
	return nil
}

func (u User) CountAndDoAsync(s string) error {
	count, err := u.Counter.Count(s)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		go func() {
			u.Doer.Do()
		}()
	}
	return nil
}
