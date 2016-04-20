package api

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	StatusPublished = iota
	StatusUnpublished
)

type Poll struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name" validate:"required" conform:"name"`
	Question  string        `json:"question" bson:"question" validate:"required" conform:"name"`
	Answers   []Answer      `json:"answers" bson:"answers"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
	Status    int           `json:"status" bson:"status" validate:"lte=1"`
}

type Answer struct {
	Label string `json:"label" bson:"label" validate:"required" conform:"name"`
	Votes int    `json:"votes" bson:"votes"`
}
