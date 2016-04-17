package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errNotFound = errors.New("Document 'Poll' not found !")
)

type QueryFilter struct {
	Sort   string `schema:"sort"`
	Limit  int    `schema:"limit"`
	Offset int    `schema:"offset"`
}

func (qf QueryFilter) SetDefault() {
	qf.Limit = 10
	qf.Offset = 0
	qf.Sort = "updated_at"
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	qf := new(QueryFilter)
	qf.SetDefault()

	if err := BindForm(r, qf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	polls, err := ListPolls(db, qf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJson(w, polls, http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := GetPoll(db, id)
	if err == mgo.ErrNotFound {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	WriteJson(w, p, http.StatusOK)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	p := new(Poll)
	p.Id = bson.NewObjectId()

	if err := BindJson(r, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := IsValid(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if err := InsertPoll(db, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJson(w, p, http.StatusCreated)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := context.Get(r, "db").(*mgo.Database)

	p := new(Poll)
	if err := BindJson(r, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := IsValid(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.UpdatedAt = time.Now()

	p.Id = bson.ObjectIdHex(id)
	if err := UpdatePoll(db, p); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	WriteJson(w, p, http.StatusOK)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	if err := DeletePoll(db, id); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
