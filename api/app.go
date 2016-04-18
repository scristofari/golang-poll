package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Bootstrap the api
// Serving / Routing
func Bootstrap() {
	port := os.Getenv("APP_PORT")

	// Handle routes
	http.Handle("/", Handlers())

	// serve
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Use gorilla mux
// CRUD + Vote Handlers
func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/api/v1/polls", MiddlewareHandler(listHandler)).Methods("GET")
	r.HandleFunc("/api/v1/polls/{poll}", MiddlewareHandler(getHandler)).Methods("GET")
	r.HandleFunc("/api/v1/polls", MiddlewareHandler(postHandler)).Methods("POST")
	r.HandleFunc("/api/v1/polls/{poll}", MiddlewareHandler(putHandler)).Methods("PUT")
	r.HandleFunc("/api/v1/polls/{poll}", MiddlewareHandler(deleteHandler)).Methods("DELETE")
	r.HandleFunc("/api/v1/polls/{poll}/answers/{answer}", MiddlewareHandler(voteHandler)).Methods("POST")

	return r
}
