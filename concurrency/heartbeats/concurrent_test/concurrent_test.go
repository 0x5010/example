package concurrent_test

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
	heartbeat, results := DoWork(done, intSlice...)

	<-heartbeat
	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}
}
