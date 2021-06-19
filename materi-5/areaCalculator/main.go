/*
Create an API that parse two int, type, then return area of it. type: square, rectangle, triangle, parallelogram (use a method and a function)

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-5/areaCalculator/main.go
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type shape struct {
	Base, Height float64
	Type         string
}

// Run the web server
func main() {
	http.HandleFunc("/", areas)
	fmt.Println("Starting RESTFul API endpoint at http://localhost:8080/") // try curl http://localhost:8080/?base=10&height=10&type=rectangle
	http.ListenAndServe(":8080", nil)
}

func areas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Check "base" parameter
	base, err := strconv.ParseFloat(r.FormValue("base"), 64)
	if err != nil {
		http.Error(w, `Invalid "base" parameter`, http.StatusInternalServerError)
		return
	}
	// Check "height" parameter
	height, err := strconv.ParseFloat(r.FormValue("height"), 64)
	if err != nil {
		http.Error(w, `Invalid "height" parameter`, http.StatusInternalServerError)
		return
	}
	// Check "type" parameter
	type_ := r.FormValue("type")
	if !(type_ == "square" || type_ == "rectangle" || type_ == "triangle" || type_ == "parallelogram") {
		http.Error(w, `Invalid "type" parameter`, http.StatusInternalServerError)
		return
	}
	// An idiomatic way to encapsulate new struct creation
	s := newShape(base, height, type_)

	if r.Method == "GET" {
		// Parse calculated area into a map
		area := map[string]float64{"result": s.area()}
		// Encode area result into json
		result, _ := json.Marshal(area)
		w.Write(result)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

// Construct a new shape struct with the given name
func newShape(base float64, height float64, type_ string) *shape {
	s := shape{Type: type_}
	s.Base = base
	s.Height = height
	return &s
}

// Calculate area method
func (s shape) area() float64 {
	if s.Type == "square" || s.Type == "rectangle" || s.Type == "parallelogram" {
		return s.Base * s.Height
	} else {
		return s.Base * s.Height / 2
	}
}
