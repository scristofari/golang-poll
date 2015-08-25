package api

import (
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type HttpQuery struct {
	Query  map[string]interface{}
	Offset int
	Limit  int
	Sort   string
}

const (
	offset int    = 0
	limit  int    = 15
	sort   string = "-updated"
)

func NewQuery() *HttpQuery {
	query := new(HttpQuery)
	query.Limit = limit
	query.Offset = offset
	query.Sort = sort

	return query
}

func (qc *HttpQuery) ParseQuery(r *http.Request) {
	q := r.URL.Query()

	query := make(map[string]interface{})

	for key := range q {
		value := q.Get(key)
		if key == "q" && value != "" {
			query["$text"] = bson.M{"$search": value}
			continue
		}
		if key == "sinceid" && value != "" {
			query["_id"] = bson.M{"$gt": bson.ObjectIdHex(value)}
			continue
		}
		if key == "since" && value != "" {
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				continue
			}
			tm := time.Unix(i, 0)
			query["updated"] = bson.M{"$gte": tm}
			continue
		}
		if value != "" {
			vint, err := strconv.Atoi(value)
			if err == nil {
				query[key] = vint
			} else {
				query[key] = value
			}

		}
	}
	delete(query, "limit")
	delete(query, "offset")
	delete(query, "sort")

	qc.Query = query
	o := q.Get("offset")
	if o != "" {
		qc.Offset, _ = strconv.Atoi(o)
	}

	l := q.Get("limit")
	if l != "" {
		v, _ := strconv.Atoi(l)
		if v > 100 {
			qc.Limit = 100
		} else {
			qc.Limit = v
		}
	} else {
		qc.Limit = 20
	}

	s := q.Get("sort")
	if s != "" {
		qc.Sort = s
	}
}
