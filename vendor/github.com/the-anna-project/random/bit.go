package random

// Bit returns a pseudo random bit value that is either 0 or 1.
func Bit() int {
	c := make(chan int, 1)

	select {
	case c <- 0:
	case c <- 1:
	}

	return <-c
}
