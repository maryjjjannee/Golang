package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Course struct {
	ID  int     `json:"id"`
	Name string  `json:"name"`
	Price      float64 `json:"price"`
	Instructor string  `json:"instructor"`
}

var CourseList []Course

func init() {
	CourseJSON := `[
	{
		"id": 1,
		"name": "Go Programming",
		"price": 29.99,
		"instructor": "John Doe"
	},
	{
		"id": 2,
		"name": "Python Programming",
		"price": 19.99,
		"instructor": "Jane Smith"
	},
	{
		"id": 3,
		"name": "JavaScript Programming",
		"price": 24.99,
		"instructor": "Alice Johnson"
	}
	]`
	err := json.Unmarshal([]byte(CourseJSON), &CourseList)
	if err != nil {
		log.Fatal(err)
	}
}

func courseHandler(w http.ResponseWriter, r *http.Request) { // handler function
	courseJSON, err := json.Marshal(CourseList) 
	switch r.Method {
	case http.MethodGet:
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPost:
		var newCourse Course
		Bodybyte, err := io.ReadAll(r.Body) // read the request body
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(Bodybyte, &newCourse) // unmarshal the JSON into a Course struct
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest) // if ID is provided, return 400 Bad Request
			return
		}
		newCourse.ID = getNextID()
		CourseList = append(CourseList, newCourse) // add the new course to the list
		w.WriteHeader(http.StatusCreated) // send a 201 Created response
		return 
	}
}

func getNextID() int {
	hightestID := 1
	for _, course := range CourseList { // find the highest ID in the list
		if course.ID > hightestID {
			hightestID = course.ID
		}
	}
	return hightestID + 1 
}

func main() {
	http.HandleFunc("/course", courseHandler) // register the handler function
	http.ListenAndServe(":8080", nil) // start the server
}
