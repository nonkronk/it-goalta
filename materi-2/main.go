/*
  ______                              ______              _                 
 / _____)                        /\  |  ___ \            | |                
| /  ___ _   _  ____  ___  ___  /  \ | |   | |_   _ ____ | | _   ____  ____ 
| | (___) | | |/ _  )/___)/___)/ /\ \| |   | | | | |    \| || \ / _  )/ ___)
| \____/| |_| ( (/ /|___ |___ | |__| | |   | | |_| | | | | |_) | (/ /| |    
 \_____/ \____|\____|___/(___/|______|_|   |_|\____|_|_|_|____/ \____)_|    
                                                                            

GuessANumber is a less fun game that challenges you to find a number based on greater than or less than feedback.
If you somehow keep guessing a wrong number, I guarantee you that it won't stop asking.
 
The flowchart can be found at https://github.com/nonkronk/it-goalta/blob/master/materi-1/GuessANumber%20Game%20algorithm%20flowchart.pdf
 */

package main

import (
	"fmt"
	"math/rand" // to generate random number
	"time" // as a seed for rand to generate random every time
)

func main(){	
	// Create a seed for random
	rand.Seed(time.Now().UnixNano())
	// Generate random number
	var r_num int = 1 + rand.Intn(10) // Total noobz level 1-10
	// Prompt the user
	for {
		guess := prompt()
		if guess == r_num {
			break
		} else if guess > r_num {
			println("Your guess is too high")
		} else {
			println("Your guess is too low")
		}
	}
	println("You win! Awesome!")
}

func prompt() int {
	// Populate input (int)
	fmt.Printf("Guess a number: ")
	var num int
	fmt.Scan(&num)
	return num
}

func init(){
	// Init function to make the game have its personalisation
	fmt.Println(`  ______                              ______              _                 
 / _____)                        /\  |  ___ \            | |                
| /  ___ _   _  ____  ___  ___  /  \ | |   | |_   _ ____ | | _   ____  ____ 
| | (___) | | |/ _  )/___)/___)/ /\ \| |   | | | | |    \| || \ / _  )/ ___)
| \____/| |_| ( (/ /|___ |___ | |__| | |   | | |_| | | | | |_) | (/ /| |    
 \_____/ \____|\____|___/(___/|______|_|   |_|\____|_|_|_|____/ \____)_|   

GuessANumber is a less fun game that challenges you to find a number based on greater than or less than feedback.
=================================================================================================================
`)
}