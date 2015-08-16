package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	port        string = "9001"
	errNotFound        = errors.New("Document 'Poll' not found")
)

func Bootstrap() {
	// Handle routes
	http.Handle("/", Handlers())

	// SERVE
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Handlers() *mux.Router {

	// ROUTING
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/api/v1/polls", MiddlewareHandler(listHandler)).Methods("GET")
	r.HandleFunc("/api/v1/polls/{id}", MiddlewareHandler(getHandler)).Methods("GET")
	r.HandleFunc("/api/v1/polls", MiddlewareHandler(postHandler)).Methods("POST")
	r.HandleFunc("/api/v1/polls/{id}", MiddlewareHandler(putHandler)).Methods("PUT")
	r.HandleFunc("/api/v1/polls/{id}", MiddlewareHandler(deleteHandler)).Methods("DELETE")

	return r
}

// ------- HANDLER CONTROLLER ------- ////

type ResultList struct {
	Total int    `json:"total"`
	Polls []Poll `json:"poll"`
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	q := NewQuery()
	q.ParseQuery(r)

	db := context.Get(r, "db").(*mgo.Database)

	var polls []Poll
	err := db.C("poll").Find(q.Query).Skip(q.Offset).Limit(q.Limit).Sort(q.Sort).All(&polls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, _ := db.C("poll").Find(q.Query).Count()

	result := ResultList{count, polls}
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(result)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	p := new(Poll)
	db := context.Get(r, "db").(*mgo.Database)
	if err := db.C("poll").FindId(bson.ObjectIdHex(id)).One(&p); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(p)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	p := new(Poll)
	p.Id = bson.NewObjectId()

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := p.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Created = time.Now()
	p.Updated = time.Now()

	db := context.Get(r, "db").(*mgo.Database)
	if err := db.C("poll").Insert(p); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	enc.Encode(p)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := context.Get(r, "db").(*mgo.Database)

	if !bson.IsObjectIdHex(id) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	p := new(Poll)
	if err := db.C("poll").FindId(bson.ObjectIdHex(id)).One(&p); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := p.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Updated = time.Now()

	bid := bson.ObjectIdHex(id)
	p.Id = bid
	if err := db.C("poll").UpdateId(bid, p); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(p)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		http.Error(w, errNotFound.Error(), http.StatusNotFound)
		return
	}

	db := context.Get(r, "db").(*mgo.Database)
	if err := db.C("poll").RemoveId(bson.ObjectIdHex(id)); err != nil {
		if err == mgo.ErrNotFound {
			http.Error(w, errNotFound.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
