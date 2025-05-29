package main

import (
	"encoding/json"
	"fmt"
)

type employee struct { // Key must start with a capital letter 
	ID int
	EmployeeName string
	Tel string
	Email string
}

func main() {
	data,_ := json.Marshal(&employee{101, "John Doe", "123-456-7890", "Johndoe@gmail.com"})
  	fmt.Println(string(data))
}