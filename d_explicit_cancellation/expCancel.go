package explicitcancellation

import (
	"fmt"
	"sync"
)

func Info() {
	fmt.Println(`
	- if main func exits without receiving all inbound vals,
	- send msg via a 'done' channel so upstream stages abandon their sends.
	- value/ type of done is irrelevant, just that it receives something.`)
}

func Done_Fanning() {
	// this done channel is shared by entire pipeline
	// close it when pipeline exits, thus closing any goroutines
	// still running (see below for more details)
	done := make(chan struct{})
	defer close(done)

	input := Done_SendIntsToChannel(done, 2, 3)
	c1 := Done_SqrIntsInChannel(done, input)
	c2 := Done_SqrIntsInChannel(done, input)

	out := DoneChannel_Merge(done, c1, c2)
	fmt.Println(<-out)

	// done will be closed by 'defer' call at end of function
}

/*
note: defer close(done) needs to be called before done is passed to this function.
doing so will ensure done is closed at the end of all poss. execution paths in program.
*/
func DoneChannel_Merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done() // this is done instead of calling Done later if defer is implemented as mentioned above
		for n := range c {
			select {
			// copy values from c to out UNLESS done receives a value
			// in which case break the loop and skip to wg.Done()
			case out <- n:
			case <-done:
				return // this line is used if defer is used above
			}
		}
	}
	// wg.Done()

	// the rest is unchanged
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

/* sending a call to done only works if ALL funcs are tweaked to include it */
func Done_SqrIntsInChannel(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

/* own attempt at adding done to a function */
func Done_SendIntsToChannel(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}

		}
	}()
	return out
}
