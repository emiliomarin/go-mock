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

		mockCounter := &mockCounter{}

		u := user.User{
			Counter: mockCounter,
		}

		Convey("When it's successful", func() {
			mockCounter.countFn = func(s string) (int, error) { return expectedCount, nil }

			err := u.CountAndDo(s)

			Convey("Should return no error and the expected count", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the counter fails", func() {
			mockCounter.countFn = func(s string) (int, error) { return 0, errors.New("some-error") }

			err := u.CountAndDo(s)

			Convey("Should return error and 0", func() {
				So(err, ShouldNotBeNil)
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

		mockCounter := mocks.NewMockCounter(ctrl)
		u := user.User{
			Counter: mockCounter,
		}

		Convey("When it's successful", func() {
			// Instead of foo we could do gomock.Any() if we don't care about the input
			mockCounter.EXPECT().Count("foo").Return(expectedCount, nil).Times(1)

			err := u.CountAndDo(s)

			Convey("Should return no error and the expected count", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the counter fails", func() {
			mockCounter.EXPECT().Count("foo").Return(0, errors.New("some-error")).Times(1)

			err := u.CountAndDo(s)

			Convey("Should return error and 0", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When we want the mock to do something extra with the input", func() {
			mockCounter.
				EXPECT().
				Count(gomock.Not("bar")).
				Return(0, errors.New("some-error")).
				Times(1).
				Do(func(x string) {
					log.Println("Received:", x)
				})

			err := u.CountAndDo(s)

			Convey("Should return error and 0", func() {
				So(err, ShouldNotBeNil)
			})
		})

	})
}

type mockCounter struct {
	countFn func(s string) (int, error)
}

func (m mockCounter) Count(s string) (int, error) {
	return mockCounter.countFn(s)
}

type mockDoer struct {
	doFn func() error
}

func (m mockDoer) Do() error {
	return mockCounter.doFn()
}
