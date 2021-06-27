/*
API Spesification
1. Create API using this(https://dbdiagram.io/d/608dd01eb29a09603d12f9e9) DB Diagram.
2. Please build API with Create, Read(List with filter by ID), Update, Delete, for Customer Table, Garage Table, and Car Table. (Pagination on Read List will be an additional value)
3. Build API Transaction Order. (explain down below)

Also available at: https://github.com/nonkronk/it-goalta/tree/master/materi-7
*/
package main

import (
	"project/config"
	"project/routes"
)

func main() {
	config.InitDB()
	e := routes.New()
	e.Start(":8000")
}
