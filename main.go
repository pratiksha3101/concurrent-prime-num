package main

import (
	"fmt"
	"sync"
)

// Write a concurrent Go program that prints the number of prime numbers between 1 and N, where N is a user-defined input
func main() {
	n := 100

	numWorkers := 4
	chunkSize := (n + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	results := make(chan map[int]bool, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i*chunkSize + 1
		end := (i + 1) * chunkSize
		if end > n {
			end = n
		}

		wg.Add(1)
		go worker(start, end, results, &wg)
	}

	// closing
	go func() {
		wg.Wait()
		close(results)
	}()

	allPrimes := make(map[int]bool)
	for result := range results {
		for num := range result {
			allPrimes[num] = true
		}
	}

	fmt.Println("Prime numbers in order:")
	for i := 1; i <= n; i++ {
		if allPrimes[i] {
			fmt.Printf("%d ", i)
		}
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func worker(start, end int, results chan<- map[int]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	primes := make(map[int]bool)

	for i := start; i <= end; i++ {
		if isPrime(i) {
			primes[i] = true
		}
	}
	results <- primes
}
