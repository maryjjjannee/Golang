package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Course struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
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

func coursesHandler(w http.ResponseWriter, r *http.Request) { // handler function
	courseJSON, err := json.Marshal(CourseList)
	switch r.Method {
	case http.MethodGet:
		if err != nil {
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
		w.WriteHeader(http.StatusCreated)          // send a 201 Created response
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

func findID(ID int) (*Course, int) {
	for i, course := range CourseList { // for range loop to find the course by ID
		if course.ID == ID {
			return &course, i
		}
	}
	return nil, 0
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1]) // get the last segment of the URL path
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound) // if ID is not a valid integer, return 400 Bad Request
		return
	}
	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("No course with Id%d", ID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		courseJSON, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPut:
		var updatedCourse Course
		byteBody, err := io.ReadAll(r.Body) // read the request body
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(byteBody, &updatedCourse) // unmarshal the JSON into a Course struct
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedCourse.ID != ID {
			w.WriteHeader(http.StatusBadRequest) // if ID is provided, return 400 Bad Request
			return
		}
		course = &updatedCourse // update the course with the new values
		CourseList[listItemIndex] = *course // update the course in the list
		w.WriteHeader(http.StatusOK) // send a 200 OK response
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // if the method is not GET or PUT, return 405 Method Not Allowed
	}


}

func middlewareHandler(handler http.Handler)  http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before handler middleware started")
		handler.ServeHTTP(w, r)
		fmt.Println("After handler middleware finished")
	})
}

func main() {
	courseItemHandler := http.HandlerFunc(courseHandler)
	courseListHandler := http.HandlerFunc(coursesHandler)
	http.Handle("/course/", middlewareHandler(courseItemHandler)) //especially for the course item
	http.Handle("/course", middlewareHandler(courseListHandler)) 
	http.ListenAndServe(":8080", nil)
}
