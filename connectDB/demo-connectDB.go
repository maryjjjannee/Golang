package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func query(db *sql.DB) {
	var (
		id          int
		course_name string
		price       float64
		instructor  string
	)
	
	for {
		var inputId int
		fmt.Scan(&inputId) // Scan user input for course ID

		query := "SELECT id, course_name, price, instructor FROM coursedb.online_course WHERE id = ?"
		if err := db.QueryRow(query, inputId).Scan(&id, &course_name, &price, &instructor); err != nil { // query the database
			log.Fatal(err)
		}
		fmt.Println(id, course_name, price, instructor)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:254428@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		fmt.Println("Fail to connect to the database")
	} else {
		fmt.Println("Connected to the database successfully")
	}
	query(db)
}
