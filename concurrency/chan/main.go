package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	take := func(done, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	toInt := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	primeFinder := func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range intStream {
				integer--
				prime := true
				for divisor := integer - 1; divisor > 1; divisor-- {
					if integer%divisor == 0 {
						prime = false
						break
					}
				}
				if prime {
					select {
					case <-done:
						return
					case primeStream <- integer:
					}
				}
			}
		}()
		return primeStream
	}

	fanIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})
		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	tee := func(done, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})

		go func() {
			defer close(out1)
			defer close(out2)

			for val := range orDone(done, in) {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}

	bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				for val := range orDone(done, stream) {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	// fan-out-fan-in

	{
		fmt.Println("\nwithout fan-out-fan-in")
		randfn := func() interface{} { return rand.Intn(50000000) }
		done := make(chan interface{})
		defer close(done)
		start := time.Now()
		randIntStream := toInt(done, repeatFn(done, randfn))
		fmt.Println("Primes:")
		for prime := range take(done, primeFinder(done, randIntStream), 10) {
			fmt.Printf("\t%d\n", prime)
		}
		fmt.Printf("Search took: %v\n", time.Since(start))
	}
	{
		fmt.Println("\nfan-out-fan-in")
		randfn := func() interface{} { return rand.Intn(50000000) }
		done := make(chan interface{})
		defer close(done)
		start := time.Now()
		randIntStream := toInt(done, repeatFn(done, randfn))
		numFinders := runtime.NumCPU()
		fmt.Printf("Spinning up %d prime finders.\n", numFinders)
		fmt.Println("Primes:")
		finders := make([]<-chan interface{}, numFinders)
		for i := 0; i < numFinders; i++ {
			finders[i] = primeFinder(done, randIntStream)
		}
		for prime := range take(done, fanIn(done, finders...), 10) {
			fmt.Printf("\t%d\n", prime)
		}
		fmt.Printf("Search took: %v\n", time.Since(start))
	}

	// tee-channel
	{
		fmt.Println("\ntee-channel")
		done := make(chan interface{})
		defer close(done)

		out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))
		for val1 := range out1 {
			fmt.Printf("out1: %v, out2:%v\n", val1, <-out2)
		}
	}

	// bridge-channel
	{
		fmt.Println("\nbridge-channel")
		genVals := func() <-chan <-chan interface{} {
			chanStream := make(chan (<-chan interface{}))
			go func() {
				defer close(chanStream)
				for i := 0; i < 10; i++ {
					stream := make(chan interface{}, 1)
					stream <- i
					close(stream)
					chanStream <- stream
				}
			}()
			return chanStream
		}
		for v := range bridge(nil, genVals()) {
			fmt.Printf("%v ", v)
		}
	}
}
