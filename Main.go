package main

import (
	"fmt"
	"os"
)

type Student struct {
	Name        string
	ID          string
	English     int
	Science     int
	Maths       int
	TotalMarks  int
	TotalSecured int
	Average     float64
	Grade       string
}

var students []Student

// Function to calculate total marks, average and grade
func calculateStudentDetails(student *Student) {
	student.TotalMarks = 300
	student.TotalSecured = student.English + student.Science + student.Maths
	student.Average = float64(student.TotalSecured) / 3
	// Calculate grade
	if student.Average >= 90 {
		student.Grade = "A"
	} else if student.Average >= 80 {
		student.Grade = "B"
	} else if student.Average >= 70 {
		student.Grade = "C"
	} else if student.Average >= 60 {
		student.Grade = "D"
	} else {
		student.Grade = "F"
	}
}

// Function to display all students
func displayAllStudents() {
	if len(students) == 0 {
		fmt.Println("No students available.")
		return
	}
	for _, student := range students {
		fmt.Printf("ID: %s, Name: %s, English: %d, Science: %d, Maths: %d, Total Secured: %d, Average: %.2f, Grade: %s\n",
			student.ID, student.Name, student.English, student.Science, student.Maths, student.TotalSecured, student.Average, student.Grade)
	}
}

// Function to add student
func addStudent() {
	var student Student
	fmt.Println("Enter student details:")
	fmt.Print("Name: ")
	fmt.Scan(&student.Name)
	fmt.Print("ID: ")
	fmt.Scan(&student.ID)
	fmt.Print("English Marks: ")
	fmt.Scan(&student.English)
	fmt.Print("Science Marks: ")
	fmt.Scan(&student.Science)
	fmt.Print("Maths Marks: ")
	fmt.Scan(&student.Maths)

	// Calculate total marks, average, and grade
	calculateStudentDetails(&student)

	// Add to the list of students
	students = append(students, student)
	fmt.Println("Student added successfully!")
}

// Function to edit student details
func editStudent() {
	var studentID string
	fmt.Print("Enter student ID to edit: ")
	fmt.Scan(&studentID)

	// Search for the student
	for i, student := range students {
		if student.ID == studentID {
			fmt.Println("Editing student:", student.Name)
			fmt.Print("New Name: ")
			fmt.Scan(&students[i].Name)
			fmt.Print("New English Marks: ")
			fmt.Scan(&students[i].English)
			fmt.Print("New Science Marks: ")
			fmt.Scan(&students[i].Science)
			fmt.Print("New Maths Marks: ")
			fmt.Scan(&students[i].Maths)

			// Recalculate the details after editing
			calculateStudentDetails(&students[i])
			fmt.Println("Student details updated successfully!")
			return
		}
	}
	fmt.Println("Student not found!")
}

// Function to search student by ID
func searchStudent() {
	var studentID string
	fmt.Print("Enter student ID to search: ")
	fmt.Scan(&studentID)

	// Search for the student
	for _, student := range students {
		if student.ID == studentID {
			fmt.Printf("Student Found: ID: %s, Name: %s, English: %d, Science: %d, Maths: %d, Total Secured: %d, Average: %.2f, Grade: %s\n",
				student.ID, student.Name, student.English, student.Science, student.Maths, student.TotalSecured, student.Average, student.Grade)
			return
		}
	}
	fmt.Println("Student not found!")
}

// Function to delete student
func deleteStudent() {
	var studentID string
	fmt.Print("Enter student ID to delete: ")
	fmt.Scan(&studentID)

	// Search and remove the student
	for i, student := range students {
		if student.ID == studentID {
			students = append(students[:i], students[i+1:]...)
			fmt.Println("Student deleted successfully!")
			return
		}
	}
	fmt.Println("Student not found!")
}

func main() {
	for {
		fmt.Println("\nStudent Management System")
		fmt.Println("1. Add Student")
		fmt.Println("2. Edit Student Details")
		fmt.Println("3. Search Student")
		fmt.Println("4. Delete Student")
		fmt.Println("5. Display All Students")
		fmt.Println("6. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			addStudent()
		case 2:
			editStudent()
		case 3:
			searchStudent()
		case 4:
			deleteStudent()
		case 5:
			displayAllStudents()
		case 6:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
