package user_test

import (
	"errors"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/emiliomarin/go-mock/mocks"
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

func TestCountWithMockGen(t *testing.T) {
	Convey("Given we want to count characters in a string", t, func() {
		s := "foo"
		expectedCount := 3

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockCounter(ctrl)
		u := user.User{
			Counter: m,
		}

		Convey("When it's successful", func() {
			// Instead of foo we could do gomock.Any() if we don't care about the input
			m.EXPECT().Count("foo").Return(expectedCount, nil).Times(1)

			res, err := u.Count(s)

			Convey("Should return no error and the expected count", func() {
				So(err, ShouldBeNil)
				So(res, ShouldEqual, expectedCount)
			})
		})

		Convey("When the counter fails", func() {
			m.EXPECT().Count("foo").Return(0, errors.New("some-error")).Times(1)

			res, err := u.Count(s)

			Convey("Should return error and 0", func() {
				So(err, ShouldNotBeNil)
				So(res, ShouldEqual, 0)
			})
		})

		Convey("When we want the mock to do something extra with the input", func() {
			m.
				EXPECT().
				Count(gomock.Not("bar")).
				Return(0, errors.New("some-error")).
				Times(1).
				Do(func(x string) {
					log.Println("Received:", x)
				})

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
