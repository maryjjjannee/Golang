package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Course struct {
	CourseID   int     `json:"course_id"`
	CourseName string  `json:"course_name"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
}

var Db *sql.DB
var courseList []Course

const coursePath = "courses"
const basePath = "/api"






func getCourse(courseid int) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows := Db.QueryRowContext(ctx, `SELECT 
	course_id, 
	course_name, 
	price, 
	image_url 
	FROM course_online 
	WHERE course_id = ?`, courseid)

	course := &Course{}
	err := rows.Scan(
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}

func removeCourse(courseID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := Db.ExecContext(ctx, `DELETE FROM course_online WHERE course_id = ?`, courseID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func getCourseList() ([]Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // set time 
	defer cancel()
	results, err := Db.QueryContext(ctx, `SELECT 
	course_id, 
	course_name, 
	price, 
	image_url 
	FROM course_online`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	courses := make([]Course, 0)
	for results.Next() {
		var course Course
		results.Scan(&course.CourseID, 
			&course.CourseName, 
			&course.Price, 
			&course.ImageURL)

		courses = append(courses, course)
	}
	return courses, nil
}

func insertProduct(course Course) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.ExecContext(ctx, `INSERT INTO course_online 
	(course_id,
	course_name, 
	price, 
	image_url
	) VALUES (?, ?, ?, ?)`, 
	 course.CourseID, 
	 course.CourseName, 
	 course.Price, 
	 course.ImageURL)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func handleCourses(w http.ResponseWriter, r *http.Request) { //s
	switch r.Method {
		case http.MethodGet:
			courseList, err := getCourseList()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		j, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
			
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
		case http.MethodPost:
			var course Course
			err := json.NewDecoder(r.Body).Decode(&course)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			CourseID, err := insertProduct(course)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(fmt.Sprintf(`{"course_id":%d}`, CourseID)))
		case http.MethodOptions:
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		
	
}}

func handleCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
		case http.MethodGet:
			course, err := getCourse(courseID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if course == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			j, err := json.Marshal(course)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusBadRequest)
			}
			_, err = w.Write(j)
			if err != nil {
				log.Fatal(err)
				
			}
		case http.MethodDelete:
			err := removeCourse(courseID)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		

}}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		handler.ServeHTTP(w, r)
	})
}
func SetupRoutes(apiBasePath string) {

	courseHandler := http.HandlerFunc(handleCourse)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))
	coursesHandler := http.HandlerFunc(handleCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))
}

func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:254428@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Db)
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)

}


func main() {
	SetupDB()	
	SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
