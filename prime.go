package rsa

import (
	"math/rand"
	"time"
)

// Send a sequence of odd numbers to channel ch.
func Generate(ch chan<- int64) {
	for i := 3; ; i += 2 {
		ch <- int64(i)
	}
}

// Copy the values from channel in to channel out if not divisable by prime
// removing those divisible by 'prime'.
func Filter(in <-chan int64, out chan<- int64, prime int64) {
	for {
		i := <-in
		if i%prime != 0 {
			out <- i
		}
	}
}

//Return a random number between the given range
func randInRange(min, max int) int {
	return rand.Intn(max-min) + min
}

//Return a random prime number of size n
func CalculatePrime(n int) int64 {

	tab := []int{0, 3, 24, 167, 1228, 9591}
	rand.Seed(time.Now().UTC().UnixNano())
	primeIndex := randInRange(tab[n-1]+1, tab[n])
	var prime int64
	ch := make(chan int64) // Create a new channel.
	go Generate(ch)        // Launch Generate goroutine.
	for i := 0; i < primeIndex; i++ {
		prime = <-ch
		ch1 := make(chan int64)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
	return prime
}
