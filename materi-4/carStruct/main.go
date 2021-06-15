/*
Create method with Struct and Pointer that represent a Car

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-4/carStruct/main.go
*/

package main

import "fmt"

type Car struct {
	Name, Model      string
	Tank_Volume      float32 // liter
	Fuel_Consumption float32 // km/L
}

// This program demonstrate the struct and the pointer feature of Go
func main() {
	//
	c := Car{
		Name:             "Suzuki",
		Model:            "Swift",
		Tank_Volume:      42,
		Fuel_Consumption: 16,
	}
	fmt.Println("Name             :", c.Name)
	fmt.Println("Model            :", c.Model)
	fmt.Println("Tank Volume      :", c.Tank_Volume, "Liters")
	fmt.Println("Fuel Consumption :", c.Fuel_Consumption, "k/L")
	fmt.Println("Maximum Mileage  :", c.mileAge(), "Kilometers")
	fmt.Println("Efficiency Level :", c.fuelEff(c.Fuel_Consumption))
}

// Method to calculate the mileage of the car
func (c *Car) mileAge() float32 {
	return (c.Tank_Volume) * c.Fuel_Consumption
}

// Method to classify the car's fuel efficiency level
func (c *Car) fuelEff(fuel_cons float32) string {
	if fuel_cons < 10 {
		return "Poor"
	} else if fuel_cons < 20 {
		return "Moderate"
	} else {
		return "Superior"
	}
}
