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

// Repository for poll
type Repository struct{}

func (r *Repository) ListPolls(db *mgo.Database, qf *QueryFilter) ([]Poll, error) {
	// Empty VS Nil slice
	// Empty => the json marshaler will print []
	// Nil => the json marshaler will print null
	// Empty is the good thing to do
	polls := []Poll{}
	if err := db.C(collection).Find(nil).Skip(qf.Offset).Limit(qf.Limit).Sort(qf.Sort).All(&polls); err != nil {
		return nil, err
	}
	return polls, nil
}

func (r *Repository) GetPoll(db *mgo.Database, id string) (*Poll, error) {
	p := new(Poll)
	if err := db.C(collection).FindId(bson.ObjectIdHex(id)).One(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) InsertPoll(db *mgo.Database, p *Poll) error {
	return db.C(collection).Insert(p)
}

func (r *Repository) UpdatePoll(db *mgo.Database, p *Poll) error {
	return db.C(collection).UpdateId(p.Id, p)
}

func (r *Repository) DeletePoll(db *mgo.Database, id string) error {
	return db.C(collection).RemoveId(bson.ObjectIdHex(id))
}

func (r *Repository) VotePoll(db *mgo.Database, pollId string, answerId string) error {
	p := new(Poll)
	key := "answers." + answerId + ".votes"
	change := mgo.Change{Update: bson.M{"$inc": bson.M{key: 1}}, ReturnNew: true}

	_, err := db.C(collection).FindId(bson.ObjectIdHex(pollId)).Apply(change, &p)
	return err
}
