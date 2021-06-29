/*
Task description from the quiz:
Using app RentCar before, add middleware to it. Specification for Auth : Employees just can access Order API and Read Cars, SuperAdmin can access all

Also available at: https://github.com/nonkronk/it-goalta/tree/master/materi-8
*/
package main

import (
	"log"
	"project/config"
	"project/routes"
)

func main() {
	config.SetConfig()
	config.InitDB()
	e := routes.New()
	log.Fatal(e.Start(":" + config.Config.Port))
}
