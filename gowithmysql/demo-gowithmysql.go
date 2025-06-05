package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func creatingTable(db *sql.DB) {
	query := `CREATE TABLE users (
		id INT AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME,
		PRIMARY KEY (id)
	)`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func Insert(db *sql.DB){
	var username string
	var password string
	fmt.Scan(&username) // Scan user input for username
	fmt.Scan(&password) // Scan user input for password
	createdAt := time.Now() // Get the current time

	result, err := db.Exec("INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)", username, password, createdAt)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId() // Get the last inserted ID
	fmt.Println(id)
}

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

		query := "SELECT id, course_name, price, instructor FROM sys.online_course WHERE id = ?"
		if err := db.QueryRow(query, inputId).Scan(&id, &course_name, &price, &instructor); err != nil { // query the database
			log.Fatal(err)
		}
		fmt.Println(id, course_name, price, instructor)
	}
}

func delete(db *sql.DB) {
	var deleteid int
	fmt.Scan(&deleteid)
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, deleteid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Delete successfully")
}

func main() {
	db, err := sql.Open("mysql", "root:254428@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		fmt.Println("Fail to connect to the database")
	} else {
		fmt.Println("Connected to the database successfully")
	}
	// creatingTable(db) // Create the users table
	// Insert(db)
	delete(db)
}
