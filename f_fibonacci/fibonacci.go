package fibonacci

import (
	"fmt"
	"sync"
)

func fibonacci(n int, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()
	a, b := 0, 1
	for range n {
		ch <- a
		a, b = b, a+b
	}
	close(ch)
}

func fibonacciList(n int, fibs []int) []int {
	if len(fibs) >= n {
		//fmt.Println("Unnecessary.")
		return fibs
	}
	a, b := 0, 1
	for i := range n {
		if len(fibs) <= i {
			//fmt.Printf("Adding %d\n", a)
			fibs = append(fibs, a)
		}
		a, b = b, a+b
	}
	return fibs
}

func FibonacciChannel(size int) {
	fmt.Printf("Fibonacci #%d\n", size)
	wg := new(sync.WaitGroup)
	ch := make(chan int)

	wg.Add(1)
	go fibonacci(size, wg, ch)

	for num := range ch {
		fmt.Println(num)
	}

	wg.Wait()
}

func FibonacciLs(size int) {
	ls := []int{}

	ls = fibonacciList(size, ls)
	fmt.Println(ls)

	// will only add extra 5
	ls = fibonacciList(size+5, ls)
	fmt.Println(ls)

	// will skip function entirely
	ls = fibonacciList(size, ls)
	fmt.Println(ls)
}
