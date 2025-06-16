package stoppingshort

import "fmt"

func Info() {
	fmt.Println(`
	- stages don't always receive ALL inbound values.
	- so need a way to stop waiting for channels when unneeded
	- (otherwise blocking happens)
	- can add a buffer to channel to indicate max # of expected/ needed inputs`)
}

func Buffered_SendIntsToChannel(nums ...int) <-chan int {
	// set max buffer size to simplify how vals are added to 'out'
	out := make(chan int, len(nums))
	// don't need to add them via another goroutine
	for _, n := range nums {
		out <- n
	}
	// still need to close out when done though!
	close(out)
	return out
}
