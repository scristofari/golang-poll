package api

import (
	"gopkg.in/mgo.v2"
)

var (
	s        *mgo.Session
	err      error
	host     string = "localhost"
	database string = "poll"
)

func init() {
	s, err = mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)

	c := s.DB(database).C("poll")
	index := mgo.Index{
		Key: []string{"$text:name"},
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
