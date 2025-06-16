package fanning

import (
	"fmt"
	b "pipes/a_basic"
	"sync"
)

func Info() {
	fmt.Println(
		`- fan-out: multiple funcs reading from the same channel until it's closed.
		- fan-in: func reads from mult. channels. Builds one channel that closes when inputs are.`)
}

func Fanning() {
	nums := b.SendIntsToChannel(2, 3)

	// each one consumes one of the values passed to nums
	c1 := b.SqrIntsInChannel(nums)
	c2 := b.SqrIntsInChannel(nums)

	// get values from fanned channels
	for n := range Merge(c1, c2) {
		fmt.Println(n)
	}
}

/*
- start a goroutine for each inbound channel that copies its value to sole outbound channel
- use an additional goroutine to close outbound channel when inbounds are done sending
- use sync.Waitgroup to make sure all sends are done.
*/
func Merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// takes value from input channel and adds it to 'out'
	// signals "done" to wg when done
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	// set capacity of wait group (= # of input channels)
	wg.Add(len(channels))
	// pass each input channel to output function -> channel is consumed
	// and value is added to sole output channel
	for _, c := range channels {
		go output(c)
	}

	// start a goroutine to Wait until wg is empty before closing out
	// must be called AFTER .Add call:
	go func() {
		wg.Wait()
		close(out)
	}() // no params passed in so leave this as ()

	return out
}
