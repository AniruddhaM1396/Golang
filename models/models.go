// models/models.go
package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Username string             `json:"username,omitempty" bson:"username,omitempty"`
    Email    string             `json:"email,omitempty" bson:"email,omitempty"`
}

type Group struct {
    ID      primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
    Name    string               `json:"name,omitempty" bson:"name,omitempty"`
    Members []primitive.ObjectID `json:"members,omitempty" bson:"members,omitempty"`
}

// Add Expense, Comment, Settlement structs as needed