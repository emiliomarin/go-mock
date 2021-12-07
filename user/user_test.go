package user_test

import (
	"errors"
	"log"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/emiliomarin/go-mock/mocks"
	"github.com/emiliomarin/go-mock/user"
)

func TestCountAndDoWithManualMock(t *testing.T) {
	Convey("Given we want to count characters in a string", t, func() {
		s := "foo"
		expectedCount := 3

		mockCounter := &mockCounter{}
		mockDoer := &mockDoer{}

		u := user.User{
			Counter: mockCounter,
			Doer:    mockDoer,
		}

		Convey("When it's successful", func() {
			mockCounter.countFn = func(s string) (int, error) { return expectedCount, nil }
			mockDoer.doFn = func() error { return nil }

			err := u.CountAndDo(s)

			Convey("Should return no error and the expected count", func() {
				So(err, ShouldBeNil)
				So(mockDoer.calls, ShouldEqual, expectedCount)
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

func TestCountAndDoWithMockGen(t *testing.T) {
	Convey("Given we want to count characters in a string", t, func() {
		s := "foo"
		expectedCount := 3

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCounter := mocks.NewMockCounter(ctrl)
		mockDoer := mocks.NewMockDoer(ctrl)
		u := user.User{
			Counter: mockCounter,
			Doer:    mockDoer,
		}

		Convey("When it's successful", func() {
			// Instead of foo we could do gomock.Any() if we don't care about the input
			mockCounter.EXPECT().Count("foo").Return(expectedCount, nil).Times(1)
			mockDoer.EXPECT().Do().Return(nil).Times(expectedCount)

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
			mockDoer.EXPECT().Do().Return(nil).Times(0)

			err := u.CountAndDo(s)

			Convey("Should return error and 0", func() {
				So(err, ShouldNotBeNil)
			})
		})

	})
}

func TestCountAndDoAsync(t *testing.T) {
	Convey("Given we want to count and do async", t, func() {
		s := "foo"
		expectedCount := 3

		Convey("When we use manual mocks", func() {
			mockCounter := &mockCounter{
				countFn: func(s string) (int, error) { return expectedCount, nil },
			}
			mockDoer := &mockDoer{
				doFn: func() error { return nil },
				wg:   &sync.WaitGroup{},
			}

			u := user.User{
				Counter: mockCounter,
				Doer:    mockDoer,
			}

			mockDoer.wg.Add(expectedCount)
			err := u.CountAndDoAsync(s)
			mockDoer.wg.Wait()

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
				So(mockDoer.calls, ShouldEqual, expectedCount)
			})

		})

		Convey("When we use mockgen", func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			var wg sync.WaitGroup

			mockCounter := mocks.NewMockCounter(ctrl)
			mockDoer := mocks.NewMockDoer(ctrl)

			u := user.User{
				Counter: mockCounter,
				Doer:    mockDoer,
			}

			mockCounter.EXPECT().Count(gomock.Any()).Return(expectedCount, nil).Times(1)
			mockDoer.EXPECT().Do().Return(nil).Times(expectedCount).Do(func() { wg.Done() })

			wg.Add(expectedCount)
			err := u.CountAndDoAsync(s)
			wg.Wait()

			Convey("Should return no error", func() {
				So(err, ShouldBeNil)
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

type mockDoer struct {
	calls int
	doFn  func() error
	wg    *sync.WaitGroup
}

func (m *mockDoer) Do() error {
	if m.wg != nil {
		defer m.wg.Done()
	}
	m.calls++
	return m.doFn()
}
