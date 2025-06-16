package main

import (
	"fmt"
	b "pipes/a_basic"
	f "pipes/b_fanning"
	s "pipes/c_stopping_short"
	e "pipes/d_explicit_cancellation"
	d "pipes/e_digest_a_tree"
	fi "pipes/f_fibonacci"
)

func main() {
	Fibonacci()
}

func Basic() {
	fmt.Println("Basic Example:")
	b.Info()
	b.SetupPipeline()
}

func Fanning() {
	fmt.Println("Fanning Example:")
	f.Info()
	f.Fanning()
}

func StopShort() {
	fmt.Println("Stopping Short (with Buffers):")
	s.Info()
}

func ExpCancel() {
	fmt.Println("Explicit Cancellation with Done Channel:")
	e.Info()
	e.Done_Fanning()
}

func TreeDigest() {
	// run with "go run . ."
	d.DefaultTreeDigest()
	d.ParallelTreeDigest()
	d.BoundedParallelTreeDigest()
}

func Fibonacci() {
	fi.FibonacciChannel(10)
	fi.FibonacciLs(10)
}
