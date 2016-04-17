package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Bootstrap() {
	port := os.Getenv("APP_PORT")

	// Handle routes
	http.Handle("/", Handlers())
	// SERVE
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Handlers() *mux.Router {

	// ROUTING
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/api/v1/polls", MiddlewareHandler(listHandler)).Methods("GET")
	r.HandleFunc("/api/v1/polls/{poll}", MiddlewareHandler(getHandler)).Methods("GET")
	r.HandleFunc("/api/v1/polls", MiddlewareHandler(postHandler)).Methods("POST")
	r.HandleFunc("/api/v1/polls/{poll}", MiddlewareHandler(putHandler)).Methods("PUT")
	r.HandleFunc("/api/v1/polls/{poll}", MiddlewareHandler(deleteHandler)).Methods("DELETE")
	r.HandleFunc("/api/v1/polls/{poll}/answers/{answer}", MiddlewareHandler(deleteHandler)).Methods("POST")

	return r
}
