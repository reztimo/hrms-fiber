// models/employee.go
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name"`
	Salary float64            `json:"salary"`
	Age    float64            `json:"age"`
}
