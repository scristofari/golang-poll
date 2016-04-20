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
	// Custom error for document not found : mgo.ErrNotFound
	errNotFound = errors.New("Document 'Poll' not found !")
)

// Will permit to filter polls by :
// - limit
// - offset
// - sort
type QueryFilter struct {
	Sort   string `schema:"sort"`
	Limit  int    `schema:"limit"`
	Offset int    `schema:"offset"`
}

// Set default values for the query filter.
// Gorilla schema will override filters if they are set in the query of the url.
func NewQueryFilter() *QueryFilter {
	qf := new(QueryFilter)
	qf.Limit = 10
	qf.Offset = 0
	qf.Sort = "updated_at"
	return qf
}

// List all the polls.
// Workflow :
// - Get the filters for the list
// - Get the db session thanks to the context
// - List polls
// - Render the result in json format
func listHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	qf := NewQueryFilter()
	if err := BindForm(r, qf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repo := new(Repository)
	polls, err := repo.ListPolls(db, qf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, polls, http.StatusOK)
}

// Get one poll.
// Workflow :
// - Get the id from the router and validate it
// - Get the db session thanks to the context
// - Get poll
// - Render the result in json format
func getHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	vars := mux.Vars(r)
	id := vars["poll"]
	if !bson.IsObjectIdHex(id) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	repo := new(Repository)
	p, err := repo.GetPoll(db, id)
	if err == mgo.ErrNotFound {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, p, http.StatusOK)
}

// Create one poll.
// Workflow :
// - Generate a new objectId
// - Deserialize the body from json to struct
// - Validate the struct
// - Get the db session thanks to the context
// - Insert the poll
// - Render the result in json format
func postHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	p := new(Poll)
	p.ID = bson.NewObjectId()

	if err := BindJSON(r, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := IsValid(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	repo := new(Repository)
	if err := repo.InsertPoll(db, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, p, http.StatusCreated)
}

// Update one poll
// Workflow :
// - Get the id from the router and validate it
// - Deserialize the body from json to struct
// - Validate the struct
// - Get the db session thanks to the context
// - Update the poll
// - Render the result in json format
func putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["poll"]
	if !bson.IsObjectIdHex(id) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)

	p := new(Poll)
	if err := BindJSON(r, p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := IsValid(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.UpdatedAt = time.Now()

	p.ID = bson.ObjectIdHex(id)
	repo := new(Repository)
	if err := repo.UpdatePoll(db, p); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, p, http.StatusOK)
}

// Delete one poll
// Workflow :
// - Get the id from the router and validate it
// - Get the db session thanks to the context
// - Delete the poll
// - Return a http status 204 No content
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	vars := mux.Vars(r)
	id := vars["poll"]
	if !bson.IsObjectIdHex(id) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	repo := new(Repository)
	if err := repo.DeletePoll(db, id); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Vote
// Workflow :
// - Get the ids (poll / answer) from the router and validate them
// - Get the db session thanks to the context
// - Vote for the poll
// - Return a http status 204 No content
func voteHandler(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(*mgo.Database)

	vars := mux.Vars(r)
	poll := vars["poll"]
	answer := vars["answer"]

	if !bson.IsObjectIdHex(poll) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	repo := new(Repository)
	if err := repo.VotePoll(db, poll, answer); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusNoContent)
}
