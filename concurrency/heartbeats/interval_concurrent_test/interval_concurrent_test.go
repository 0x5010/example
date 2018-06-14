package interval_concurrent_test

import (
	"fmt"
	"testing"
	"time"
)

func DoWork(done <-chan interface{}, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeatStream)
		defer close(intStream)
		time.Sleep(2 * time.Second)
		pulse := time.Tick(pulseInterval)
	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeatStream <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop
				}
			}

		}
	}()
	return heartbeatStream, intStream
}

func TestDoWork_GeneratesAllnumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := DoWork(done, timeout/2, intSlice...)

	<-heartbeat
	i := 0
	for {
		select {
		case r, ok := <-results:
			if !ok {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
			i++
		case <-heartbeat:
			fmt.Println("get heartbeat")
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}

	}
}
