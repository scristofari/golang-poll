package api

import (
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	s          *mgo.Session
	err        error
	database   string = "poll"
	collection string = "poll"
)

func init() {
	s, err = mgo.Dial(os.Getenv("APP_HOST"))
	if err != nil {
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)

	// @travis mongodb not have text search enabled
	if os.Getenv("APP_ENV") != "travis" {
		c := s.DB(database).C(collection)
		index := mgo.Index{
			Key: []string{"$text:name"},
		}

		err := c.EnsureIndex(index)
		if err != nil {
			panic(err)
		}
	}
}

func ListPolls(db *mgo.Database, qf *QueryFilter) ([]Poll, error) {
	var polls []Poll

	if err := db.C(collection).Find(nil).Skip(qf.Offset).Limit(qf.Limit).Sort(qf.Sort).All(&polls); err != nil {
		return nil, err
	}

	return polls, nil
}

func GetPoll(db *mgo.Database, id string) (*Poll, error) {
	p := new(Poll)

	if err := db.C(collection).FindId(bson.ObjectIdHex(id)).One(&p); err != nil {
		return nil, err
	}

	return p, nil
}

func InsertPoll(db *mgo.Database, p *Poll) error {
	return db.C(collection).Insert(p)
}

func UpdatePoll(db *mgo.Database, p *Poll) error {
	return db.C(collection).UpdateId(p.Id, p)
}

func DeletePoll(db *mgo.Database, id string) error {
	return db.C(collection).RemoveId(bson.ObjectIdHex(id))
}
