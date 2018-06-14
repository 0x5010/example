package bad_concurrent_test

import (
	"testing"
	"time"
)

func DoWork(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeatStream)
		defer close(intStream)
		time.Sleep(2 * time.Second)

		for _, n := range nums {
			select {
			case heartbeatStream <- struct{}{}:
			default:
			}
			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()
	return heartbeatStream, intStream
}

func TestDoWork_GeneratesAllnumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := DoWork(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("test timed out")
		}
	}
}
