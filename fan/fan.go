// Package fan shows how to use fan pattern to parallelize intensive process.
// It fan out the randomInt generated to several goroutine that will compute primes in parallels.
// Then it will fan in all result in one stream.
package fan

import (
	"log/slog"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// repeatFunc will generate data thanks to fn func and send it to the read only chan returned.
func repeatFunc[T, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()

	return stream
}

// take takes n values from the stream and put in the taken steam.
// It allows to control the number of values taken from the generator.
func take[T, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)

	go func() {
		defer close(taken)

		for range n {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()

	return taken
}

// primeFinder will wait for int in the randIntSteam and check if the int is a prime.
// It will send all prime found in the returned chan.
func primeFinder(done <-chan bool, randIntSteam <-chan int) <-chan int {
	// simple algo to have a slow process
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}

		return true
	}

	primes := make(chan int)

	go func() {
		defer close(primes)

		for {
			select {
			case <-done:
				return
			case randomInt := <-randIntSteam:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()

	return primes
}

func fanIn[T any](done <-chan bool, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup

	fannedInStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()

		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}

	for _, c := range channels {
		wg.Add(1)

		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}

// Run runs the prime calculator by using the fan in fan out pattern.
func Run() {
	start := time.Now()

	done := make(chan bool)
	defer close(done)

	// Random int generator.
	randomNumFetcher := func() int {
		return rand.Intn(500000000)
	}

	// Stream of random int.
	randIntStream := repeatFunc(done, randomNumFetcher)

	// Stream of Primes.
	// Fan out. We will use several go routine to compute primes.
	CPUCount := runtime.NumCPU()
	primesFinderChannels := make([]<-chan int, CPUCount)

	for i := range CPUCount {
		primesFinderChannels[i] = primeFinder(done, randIntStream)
	}

	// Fan in
	fannedInStream := fanIn(done, primesFinderChannels...)

	for rando := range take(done, fannedInStream, 10) {
		slog.Info("Get prime from stream",
			slog.Int("value", rando),
		)
	}

	slog.Info("Run finished",
		slog.String("duration", time.Since(start).String()),
		slog.Int("PrimeFinderCount", CPUCount))
}
