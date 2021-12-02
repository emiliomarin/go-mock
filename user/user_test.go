package user_test

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/emiliomarin/go-mock/user"
)

func TestCountWithManualMock(t *testing.T) {
	Convey("Given we want to count characters in a string", t, func() {
		s := "foo"
		expectedCount := 3

		m := &mockCounter{}
		u := user.User{
			Counter: m,
		}
		Convey("When it's successful", func() {
			m.countFn = func(s string) (int, error) { return expectedCount, nil }

			res, err := u.Count(s)

			Convey("Should return no error and the expected count", func() {
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expectedCount)
			})
		})

		Convey("When the counter fails", func() {
			m.countFn = func(s string) (int, error) { return 0, errors.New("some-error") }

			res, err := u.Count(s)

			Convey("Should return error and 0", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldEqual, 0)
			})
		})
	})
}

type mockCounter struct {
	countFn func(s string) (int, error)
}

func (m mockCounter) Count(s string) (int, error) {
	return m.countFn(s)
}
