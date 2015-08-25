package api

import (
	"gopkg.in/mgo.v2"
)

var (
	s        *mgo.Session
	err      error
	host     string = "127.0.0.1"
	database string = "sparck-poll"
)

func init() {
	s, err = mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)

	c := s.DB(database).C("poll")

	// @travis mongodb not have text search enabled
	/*
		index := mgo.Index{
			Key: []string{"$text:name"},
		}

		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	*/
}
