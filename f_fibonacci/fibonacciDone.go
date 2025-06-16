package fibonacci

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

func fibonacciChDefault(n int, wg *sync.WaitGroup, ch chan<- int, done <-chan struct{}) {
	fibonacciCh(n, wg, ch, done, 0, 1)
}

func fibonacciCh(n int, wg *sync.WaitGroup, ch chan<- int, done <-chan struct{}, startA, startB int) {
	defer wg.Done()
	a, b := startA, startB

	for range n {
		select {
		case ch <- a:
			a, b = b, a+b
		case <-done:
			return
		}
	}
	close(ch)
}

func fibonacciListCh(n int, fibs []int, done <-chan struct{}) []int {
	if len(fibs) >= n {
		return fibs
	}

	neededCount := n - len(fibs)
	ch := make(chan int)
	var wg sync.WaitGroup

	// Reuse the fibonacci function
	wg.Add(1)
	go func() {
		// If continuing an existing sequence, adjust starting values
		a, b := 0, 1
		if len(fibs) > 2 {
			// if fibs has more than just 0,1 -> get 2 latest numbers as starting values. else use default 0,1
			a, b = fibs[len(fibs)-1], fibs[len(fibs)-1]+fibs[len(fibs)-2]
		}
		fibonacciCh(neededCount, &wg, ch, done, a, b)
	}()

	result := slices.Clone(fibs)

collect:
	for {
		select {
		case num, ok := <-ch:
			if !ok {
				break collect
			}
			result = append(result, num)
		case <-done:
			wg.Wait()
			return result
		}
	}

	wg.Wait()
	return result
}

func FibonacciWithDone() {
	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	ch := make(chan int)

	// Start fibonacci generator
	wg.Add(1)
	go fibonacciChDefault(100, wg, ch, done) // CAN generate up to 100 numbers...

	// but will be cancelled after 5
	go func() {
		for range 5 {
			fmt.Println(<-ch)
		}
		close(done)
	}()

	wg.Wait()
	fmt.Println("Done")
}

func FibonacciListWithDone() {
	done := make(chan struct{})
	ls := []int{}

	// First generate 10 numbers
	ls = fibonacciListCh(10, ls, done)
	fmt.Println(ls)

	// this will close done channel after given time
	// thus terminating the following 'calculate 100 numbers' function.
	go func() {
		<-time.After(1 * time.Millisecond)
		close(done)
	}()

	ls = fibonacciListCh(100, ls, done)
	fmt.Println(ls)
}
