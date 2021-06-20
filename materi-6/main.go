/*
Create an API with this specification (don’t forget to using all we learn before)

Also available at: https://github.com/nonkronk/it-goalta/blob/master/materi-6/main.go
*/

package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type User struct {
	ID       int     `json: "id" form: "id" param: "id" query: "id"`
	Fullname string  `json: "fullname" form: "fullname" param: "fullname" query: "fullname"`
	DOB      string  `json: "dob" form: "dob" param: "dob" query: "dob"`
	Height   int     `json: "height" form: "height" param: "height" query: "height"`
	Weight   int     `json: "weight" form: "weight" param: "weight" query: "weight"`
	Age      int     `json: "age" form: "age" param: "age" query: "age"`
	BMI      float64 `json: "bmi" form: "bmi" param: "bmi" query: "bmi"`
	Status   string  `json: "status" form: "status" param: "status" query: "status"`
}

// Store data in Users map
var users map[int]*User
var seq int

func init() {
	// I put myself alone before anyone else
	users = map[int]*User{
		1: &User{
			ID:       1,
			Fullname: "Irvan Tristian",
			DOB:      "17/08/1945",
			Height:   178,
			Weight:   45,
		},
	}
	users[int(1)].calculateAge()
	users[int(1)].calculateBMI()
	users[int(1)].statusBMI()
	seq = 2
}

func main() {
	e := echo.New()

	// ROUTES
	// GET /users // list of users
	// GET /users/1 // get user with id 1
	// POST /users // create new user

	e.GET("/users", UsersController)
	e.GET("/users/:id", UserController)
	e.POST("/users", RegisterController)

	// I create special endpoint to get the context of the response, based on the quiz
	// GET /bmi
	// Request: {"fullname": "Rachman Kamil", "dob": "24/10/1992", "height": 175, "weight": 70}
	// Response: {"fullname": "Rachman Kamil", "age": "28", "bmi": 22.8, "status": normal}
	e.GET("/bmi", AnswerController) // API endpoint to test given the response
	e.GET("/bmi", AnswerController) // API endpoint to test given the response

	fmt.Println()
	fmt.Println("API Request:    ")
	fmt.Println()
	fmt.Println(`--->   curl http://localhost:8000/bmi -X GET -H "Content-Type: application/json" -d '{"fullname": "Rachman Kamil", "dob": "24/10/1992", "height": 175, "weight": 70}'`)
	fmt.Println()

	// Hidebanner & Start server
	e.HideBanner = true
	e.Start(":8000")
}

func UsersController(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func UserController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}

func RegisterController(c echo.Context) error {
	u := new(User)
	u.ID = seq
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user := &User{
		ID:       u.ID,
		Fullname: u.Fullname,
		DOB:      u.DOB,
		Height:   u.Height,
		Weight:   u.Weight,
		Age:      u.Age,
		BMI:      u.BMI,
		Status:   u.Status,
	}
	user.calculateAge()
	user.calculateBMI()
	user.statusBMI()
	users[int(u.ID)] = user
	seq++
	return c.JSON(http.StatusCreated, user)
}

func AnswerController(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	u.calculateAge()
	u.calculateBMI()
	u.statusBMI()
	var answer struct {
		Fullname string
		Age      int
		BMI      float64
		Status   string
	}
	answer.Fullname = u.Fullname
	answer.Age = u.Age
	answer.BMI = u.BMI
	answer.Status = u.Status
	return c.JSON(http.StatusCreated, answer)
}

// Calculate BMI
// BMI Formula: BMI = weight / (height * height)
func (u *User) calculateBMI() {
	weight := float64(u.Weight)
	height := float64(u.Height)
	bmi := weight / ((height / 100) * (height / 100))
	// Floor rounded like in the example of the quiz
	u.BMI = math.Floor(bmi*10) / 10
}

// Return BMI Status
// BMI < 18.5 --> Underweight
// 18.5 < BMI < 24.9 --> Normal
// 25 < BMI < 29.9 --> Overweight
// BMI > 30 --> Obesity
func (u *User) statusBMI() {
	if u.BMI < 18.5 {
		u.Status = "Underweight"
	} else if u.BMI < 24.9 {
		u.Status = "Normal"
	} else if u.BMI < 29.9 {
		u.Status = "Overweight"
	} else {
		u.Status = "Obesity"
	}
}

// Calculate Age
// Return diff between Time.now - Time.DateofBirth
func (u *User) calculateAge() {
	layOut := "02/01/2006" // dd/mm/yyyy
	dobTime, _ := time.Parse(layOut, u.DOB)
	ageYear := time.Now().Year() - dobTime.Year()
	// This is the approach to calculate age accurately
	dobDayMonth, _ := strconv.Atoi(strconv.Itoa(dobTime.Day()) + strconv.Itoa(monthToInt(dobTime.Month().String())))
	nowDayMonth, _ := strconv.Atoi(strconv.Itoa(time.Now().Day()) + strconv.Itoa(monthToInt(time.Now().Month().String())))
	if dobDayMonth > nowDayMonth {
		ageYear = ageYear - 1
	}
	u.Age = ageYear
}

// Helper function to calculate age/time accurately
// It return month in int
func monthToInt(month string) int {
	switch month {
	case "January":
		return 1
	case "February":
		return 2
	case "March":
		return 3
	case "April":
		return 4
	case "May":
		return 5
	case "June":
		return 6
	case "July":
		return 7
	case "August":
		return 8
	case "September":
		return 9
	case "October":
		return 10
	case "November":
		return 11
	case "December":
		return 12
	default:
		panic("Unrecognized month")
	}
}
