package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

type TodoList []Todo
