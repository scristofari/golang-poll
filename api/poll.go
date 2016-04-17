package api

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Poll struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Question  string        `json:"question" bson:"question"`
	Answers   []Answer      `json:"answers" bson:"answers"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
	Status    int           `json:"status" bson:"status"`
}

type Answer struct {
	Label string `json:"label" bson:"label"`
	Votes int    `json:"votes" bson:"votes"`
}
