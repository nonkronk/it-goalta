/*
Show me how to sum an integer array.

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-3/sumOfArr/main.go
*/

package main

import "fmt"

// Print the sum of int array
func main() {
	// The array sample
	array := []int{1, 2, 3, 4, 5}
	// Print the sum of the array
	fmt.Println(sumOfArr(array))
}

// Sum all the values inside of an array
func sumOfArr(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}
