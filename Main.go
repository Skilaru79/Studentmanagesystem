package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection variables
var studentCollection *mongo.Collection

func initMongoDB() {
	// MongoDB connection string
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping MongoDB
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the database and collection
	studentCollection = client.Database("studentdb").Collection("students")
	fmt.Println("✅ Connected to MongoDB!")
}

// Student struct (model)
type Student struct {
	ID      string  `json:"id" bson:"id"`
	Name    string  `json:"name" bson:"name"`
	Email   string  `json:"email" bson:"email"`
	English int     `json:"english" bson:"english"`
	Maths   int     `json:"maths" bson:"maths"`
	Science int     `json:"science" bson:"science"`
	Total   int     `json:"total" bson:"total"`
	Average float64 `json:"average" bson:"average"`
	Grade   string  `json:"grade" bson:"grade"`
}

// Helper function to calculate total, average, and grade
func calculateStudentDetails(s *Student) {
	s.Total = s.English + s.Maths + s.Science
	s.Average = float64(s.Total) / 3

	// Assign grade based on average
	if s.Average >= 90 {
		s.Grade = "A"
	} else if s.Average >= 75 {
		s.Grade = "B"
	} else if s.Average >= 60 {
		s.Grade = "C"
	} else {
		s.Grade = "D"
	}
}

// Create Student API
func createStudent(c *gin.Context) {
	var student Student
	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Calculate total, average, and grade
	calculateStudentDetails(&student)

	// Insert student into MongoDB
	_, err := studentCollection.InsertOne(context.TODO(), student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student added successfully!", "student": student})
}

// Get All Students API
func getStudents(c *gin.Context) {
	cursor, err := studentCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
		return
	}

	var students []Student
	if err = cursor.All(context.TODO(), &students); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	c.JSON(http.StatusOK, students)
}

// Get Single Student by ID
func getStudentByID(c *gin.Context) {
	studentID := c.Param("id")

	var student Student
	err := studentCollection.FindOne(context.TODO(), bson.M{"id": studentID}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// Update Student API
func updateStudent(c *gin.Context) {
	studentID := c.Param("id")
	var updatedStudent Student

	if err := c.BindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	calculateStudentDetails(&updatedStudent)

	// Update student in MongoDB
	filter := bson.M{"id": studentID}
	update := bson.M{"$set": updatedStudent}

	_, err := studentCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully!"})
}

// Delete Student API
func deleteStudent(c *gin.Context) {
	studentID := c.Param("id")

	_, err := studentCollection.DeleteOne(context.TODO(), bson.M{"id": studentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully!"})
}

func main() {
	// Initialize MongoDB connection
	initMongoDB()

	// Set up Gin router
	router := gin.Default()

	// Define API endpoints
	router.POST("/students", createStudent)
	router.GET("/students", getStudents)
	router.GET("/students/:id", getStudentByID)
	router.PUT("/students/:id", updateStudent)
	router.DELETE("/students/:id", deleteStudent)

	// Start the server
	fmt.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}
