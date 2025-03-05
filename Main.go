package main
package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/gin-gonic/gin"
		"go.mongodb.org/mongo-driver/bson"
		"go.mongodb.org/mongo-driver/mongo"
		"go.mongodb.org/mongo-driver/mongo/options"	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Global variable to hold the MongoDB client
var Client *mongo.Client

// Function to connect to MongoDB
func ConnectDatabase() {
	// Set up MongoDB connection string (replace with your actual MongoDB URI)
	uri := "mongodb://localhost:27017" // For local MongoDB
	// If using MongoDB Atlas, replace with the URI provided by MongoDB Atlas
	// uri := "mongodb+srv://<username>:<password>@cluster0.mongodb.net/test?retryWrites=true&w=majority"

	clientOptions := options.Client().ApplyURI(uri)

	// Create the client and connect to MongoDB
	var err error
	Client, err = mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		os.Exit(1)
	}

	// Ping the database to check the connection
	err = Client.Ping(nil, readpref.Primary())
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
		os.Exit(1)
	}

	// Log success message if connected
	fmt.Println("Successfully connected to MongoDB!")
}



)

func main() {
	ConnectDatabase()
	r := gin.Default()

	// Routes
	r.POST("/students", addStudent)
	r.PUT("/students/:id", editStudent)
	r.GET("/students/:id", searchStudent)
	r.DELETE("/students/:id", deleteStudent)
	r.GET("/students", displayAllStudents)

	// Start server
	r.Run(":8080")
}

// Function to calculate total marks, average, and grade
func calculateStudentDetails(student *Student) {
	student.TotalMarks = 300
	student.TotalSecured = student.English + student.Science + student.Maths
	student.Average = float64(student.TotalSecured) / 3

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

// Add Student (POST /students)
func addStudent(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate student details
	calculateStudentDetails(&student)

	// Insert into MongoDB
	_, err := studentCollection.InsertOne(context.TODO(), student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student added successfully", "student": student})
}

// Edit Student (PUT /students/:id)
func editStudent(c *gin.Context) {
	studentID := c.Param("id")
	var updatedStudent Student

	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calculateStudentDetails(&updatedStudent)

	filter := bson.M{"id": studentID}
	update := bson.M{"$set": updatedStudent}

	_, err := studentCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
}

// Search Student (GET /students/:id)
func searchStudent(c *gin.Context) {
	studentID := c.Param("id")
	var student Student

	filter := bson.M{"id": studentID}
	err := studentCollection.FindOne(context.TODO(), filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, student)
}

// Delete Student (DELETE /students/:id)
func deleteStudent(c *gin.Context) {
	studentID := c.Param("id")
	filter := bson.M{"id": studentID}

	_, err := studentCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

// Display All Students (GET /students)
func displayAllStudents(c *gin.Context) {
	var students []Student
	cursor, err := studentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var student Student
		cursor.Decode(&student)
		students = append(students, student)
	}

	c.JSON(http.StatusOK, students)
}
