package models

import "time"

type Emp struct {
	Employee
	Employer
}

// Employee it is User information from Organization server
type Employee struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Patronymic string    `json:"patronymic"`
	BirthDate  time.Time `json:"birth_date"`
	Email      string    `json:"email"`
}

// Employee it is Employer information from Organization server
type Employer struct {
	Name       string `json:"name"`
	City       string `json:"city"`
	Office     string `json:"office"`
	Department string `json:"department"`
}
