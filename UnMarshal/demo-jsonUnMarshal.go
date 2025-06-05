package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee struct { // Key must start with a capital letter 
	ID int
	EmployeeName string
	Tel string
	Email string
}

func main() {
	e := &employee{} // create a pointer to the employee struct
	err := json.Unmarshal([]byte(`{"ID":101,"EmployeeName":"John Doe","Tel":"123-456-7890","Email":"JohnDoe@gmail.com"}`), e)
	if err != nil { 
		log.Fatal(err) // Fatal will log the error and stop the program
	}
	fmt.Println(e)
}