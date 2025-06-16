package basic

import "fmt"

func Info() {
	fmt.Println(`Pipelines (informal concept):
	- series of stages connected by channels: stage = group of goroutines running the same function.
	- goroutines: receive values from upstream (from inbound channels) and output via outbound channels.
	- first stage only has outbound channels (source/ producer stage),
	- last stage only has inbound (sink/ consumer stage).`)
}

/* simple example: squaring numbers */
// first stage
func SendIntsToChannel(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		// ignore index, add JUST the nums to the channel defined above
		for _, n := range nums {
			out <- n
		}
		// close channel once all numbers have been added
		close(out)
	}()
	// return channel
	return out
}

// second stage
func SqrIntsInChannel(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// for each number in input channel (in), add the sqr to output chan (out)
		for n := range in {
			out <- n * n
		}
		// close channel when all values are added
		close(out)
	}()
	// return channel
	return out
}

// third stage
func SetupPipeline() {
	// setup first stage, send values to second
	c := SendIntsToChannel(2, 3, 4)
	output := SqrIntsInChannel(c)

	// output results by consuming output channel
	for o := range output {
		fmt.Println(o)
	}

	// can also chain a stage if its input + output chan types match:
	sqrOutput := SqrIntsInChannel(SqrIntsInChannel(SendIntsToChannel(1, 2, 3)))
	for s := range sqrOutput {
		fmt.Println(s)
	}
}
