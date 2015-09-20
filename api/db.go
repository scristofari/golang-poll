package api

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var (
	s        *mgo.Session
	err      error
	database string = "sparck-poll"
)

func init() {
	host := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	log.Print("Host : " + host)
	s, err = mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)

	// @travis mongodb not have text search enabled
	if os.Getenv("SPARCK_ENV") != "travis" {
		c := s.DB(database).C("poll")
		index := mgo.Index{
			Key: []string{"$text:name"},
		}

		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}
