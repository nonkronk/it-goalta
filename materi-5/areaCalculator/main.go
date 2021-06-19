/*
Create an API that parse two int, type, then return area of it. type: square, rectangle, triangle, parallelogram (use a method and a function)

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-5/areaCalculator/main.go
*/

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Shape struct {
	Base   float64
	Height float64
	Type   string
}

// Run the web server
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/areas", areas)
	log.Println("Starting RESTFul API endpoint at http://localhost:8080/areas") // try curl http://localhost:8080/areas?base=10&height=10&type=rectangle
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func areas(w http.ResponseWriter, r *http.Request) {
	// Setup decoder and return error for any ignored or not matching key-value
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode json into struct
	var s Shape
	err := dec.Decode(&s)
	if err != nil {
		// Enable request by URL
		// If API accesed without body
		if err.Error() == "EOF" {
			// Check "base" parameter
			base, err := strconv.ParseFloat(r.FormValue("base"), 64)
			if err != nil {
				http.Error(w, `Invalid "base" parameter`, http.StatusBadRequest)
				return
			}
			// Check "height" parameter
			height, err := strconv.ParseFloat(r.FormValue("height"), 64)
			if err != nil {
				http.Error(w, `Invalid "height" parameter`, http.StatusBadRequest)
				return
			}
			// Check "type" parameter
			type_ := r.FormValue("type")
			if !(type_ == "square" || type_ == "rectangle" || type_ == "triangle" || type_ == "parallelogram") {
				http.Error(w, `Invalid "type" parameter`, http.StatusBadRequest)
				return
			}
			s.Base = base
			s.Height = height
			s.Type = type_
		} else {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if r.Method == "GET" {
		// Identify unrecognized value of key "type"
		if s.area() < 0 {
			log.Println(`Invalid "type" parameter`)
			http.Error(w, `Invalid "type" parameter`, http.StatusBadRequest)
			return
		} else {
			// Populate calculated area into a map
			area := map[string]float64{"result": s.area()}
			// Encode area result into json
			result, _ := json.Marshal(area)
			w.Write(result)
			return
		}
	}
	http.Error(w, "", http.StatusBadRequest)
}

// Calculate area method
func (s Shape) area() float64 {
	if s.Type == "square" || s.Type == "rectangle" || s.Type == "parallelogram" {
		return s.Base * s.Height
	} else if s.Type == "triangle" {
		return s.Base * s.Height / 2
	}
	return -1
}
