package main

type Student struct {
	Name        string  `json:"name" bson:"name"`
	ID          string  `json:"id" bson:"id"`
	English     int     `json:"english" bson:"english"`
	Science     int     `json:"science" bson:"science"`
	Maths       int     `json:"maths" bson:"maths"`
	TotalMarks  int     `json:"total_marks" bson:"total_marks"`
	TotalSecured int    `json:"total_secured" bson:"total_secured"`
	Average     float64 `json:"average" bson:"average"`
	Grade       string  `json:"grade" bson:"grade"`
}
