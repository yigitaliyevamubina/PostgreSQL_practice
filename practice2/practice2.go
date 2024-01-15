package main

import (
	"database/sql"
	"encoding/json"

	"github.com/k0kubun/pp/v3"
	_ "github.com/lib/pq"
)

type Student struct {
	Id   int
	Name string
	Age  int
}

type Course struct {
	Id    int
	Name  string
	Price float64
}

type Courses_students struct {
	Id      int
	Course  Course
	Student Student
}

func main() {
	connection := "user=postgres password=mubina2007 dbname=postgres sslmode=disable"
	mydb, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	reqString := []byte(`
	{
		"course": {
			"name": "Golang",
			"price": 1763.999
		},
		"student": {
			"name": "Paul",
			"age": 19
		}
	}
	`)

	var courses_students Courses_students
	if err := json.Unmarshal(reqString, &courses_students); err != nil {
		panic(err)
	}

	var respCourse Course
	rowCourse := mydb.QueryRow("INSERT INTO najot_courses(name, price) VALUES($1, $2) returning id, name, price",
		courses_students.Course.Name,
		courses_students.Course.Price)

	if err := rowCourse.Scan(&respCourse.Id, &respCourse.Name, &respCourse.Price); err != nil {
		panic(err)
	}

	var respStudent Student
	rowStudent := mydb.QueryRow("INSERT INTO najot_students(name, age) VALUES($1, $2) returning id, name, age",
		courses_students.Student.Name,
		courses_students.Student.Age)

	if err := rowStudent.Scan(&respStudent.Id, &respStudent.Name, &respStudent.Age); err != nil {
		panic(err)
	}

	rowCourse_student := mydb.QueryRow("INSERT INTO najot_courses_students(course_id, student_id) VALUES($1, $2) returning id, course_id, student_id", respCourse.Id, respStudent.Id)
	if err := rowCourse_student.Scan(&courses_students.Id,
		&courses_students.Course.Id,
		&courses_students.Student.Id); err != nil {
		panic(err)
	}

	var respCourse_student Courses_students
	rowCourse_student2 := mydb.QueryRow("SELECT s.id, s.name, s.age, c.id, c.name, c.price FROM najot_students s JOIN najot_courses_students cs ON s.id = cs.student_id JOIN najot_courses c ON c.id = cs.course_id WHERE s.id = $1", respStudent.Id)
	if err := rowCourse_student2.Scan(&respCourse_student.Student.Id,
		&respCourse_student.Student.Name,
		&respCourse_student.Student.Age,
		&respCourse_student.Course.Id,
		&respCourse_student.Course.Name,
		&respCourse_student.Course.Price); err != nil {
		panic(err)
	}

	pp.Println(respCourse_student)

	var respCourse_student2 Courses_students

	rowStudent2 := mydb.QueryRow("SELECT id, name, age FROM najot_students WHERE id = $1", 3)
	if err := rowStudent2.Scan(&respCourse_student2.Student.Id,
		&respCourse_student2.Student.Name,
		&respCourse_student2.Student.Age); err != nil {
		panic(err)
	}

	rowCourse_student3 := mydb.QueryRow("SELECT course_id FROM najot_courses_students WHERE student_id = $1", respCourse_student2.Student.Id)
	if err := rowCourse_student3.Scan(&respCourse_student2.Course.Id); err != nil {
		panic(err)
	}

	rowCourse1 := mydb.QueryRow("SELECT name, price FROM najot_courses WHERE id = $1", respCourse_student2.Course.Id)
	if err := rowCourse1.Scan(&respCourse_student2.Course.Name, &respCourse_student2.Course.Price); err != nil {
		panic(err)
	}

	pp.Println(respCourse_student2)
}
