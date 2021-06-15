/*
Create new Array to store a prime number with multiples of five from less than n number. And sum it for return number.

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-3/primeNumber/main.go

*TBH I didn't get what the quiz requires -,-
So, I assume....
*/

package main

import "fmt"

// "Create new Array to store a prime number with multiples of five from less than n number. And sum it for return number.""
// & I just want to print these.
func main() {
	// Populate the maximum range of prime numbers. Must be greater than 2
	fmt.Println("Enter a maximum range to generate prime numbers: ")
	var num_range int
	fmt.Scan(&num_range)

	fmt.Println("\nThe prime numbers are: ")
	primeNums := genPrime(num_range)
	fmt.Println(primeNums)

	fmt.Println("\nSo, the multiples of five of the prime numbers are: ")
	multiPrimes := multiPly(primeNums)
	fmt.Println(multiPrimes)

	fmt.Println("\nAnd the sum of the prime numbers: ")
	fmt.Println(sumArr(primeNums))

	fmt.Println("\nWhile the sum of the multiple of five of the prime numbers are: ")
	fmt.Println(sumArr(multiPrimes))
}

// Determine whether a number is a prime
func isPrime(num int) bool {
	if num == 0 || num == 1 {
		return false
	}
	for i := 2; i < num; i++ {
		if (num % i) == 0 {
			return false
		}
	}
	return true
}

// Generate an array of primes from range
func genPrime(num_range int) []int {
	result := []int{}
	for i := 0; i <= num_range; i++ {
		if isPrime(i) {
			result = append(result, i)
		}
	}
	return result
}

// Return a "5x multiplied" array
func multiPly(array []int) []int {
	result := []int{}
	for _, v := range array {
		v *= 5
		result = append(result, v)
	}
	return result
}

// Return sum of an array
func sumArr(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}
