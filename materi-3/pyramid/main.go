/*
Make a pyramid from your name from string variable with looping : *

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-3/pyramid/main.go
*/

package main

import "fmt"

// Print the the pyramid of a name
func main() {
	// Sample, my name "irvan"
	name := "irvan"
	printName(name)
}

func printName(name string) {
	// e.g. name = irvan; It will print i....irvan
	for i := range name {
		i += 1
		fmt.Println(name[:i])
	}
	// e.g. name = irvan; It will print irva....i
	for i := 0; i <= len(name)-2; i++ {
		fmt.Println(name[:len(name)-i-1])
	}
}
