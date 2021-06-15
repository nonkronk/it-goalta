/*
Create one function that calculate factorial (n!)

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-4/theFactorial/main.go
*/

package main

import "fmt"

// The program for factorial of number given
func main() {
	// Input number
	fmt.Print("Enter a number: ")
	var num int
	fmt.Scan(&num)

	// Return the factorial of the number
	result := factorN(num)
	fmt.Println("The factorial of", num, "is", result)
}

// Calculate factorial number
func factorN(num int) int {
	if num == 0 || num == 1 {
		return num
	}
	return num * factorN(num-1)
}
