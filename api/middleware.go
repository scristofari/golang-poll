package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

//  Log all request
//  Apache access_log format
func LogHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := &LogResponseWriter{
			ResponseWriter: w,
		}
		f(nw, r)

		log.Println(fmt.Sprintf("%s %s %s %d %d", r.Method, r.URL.String(), r.Proto, nw.Status(), nw.Size()))
	}
}

// Get a new mongo session by request : "crash friendly".
// Close the session after using it.
func DBHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := s.Copy()
		defer db.Close()
		context.Set(r, "db", db.DB(database))
		f(w, r)
	}
}

// Print the right headers when the method "OPTIONS" is called.
// Allow all origin (*)
func CorsHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		f(w, r)
	}
}

// Map all sub middlewares
func MiddlewareHandler(f http.HandlerFunc) http.HandlerFunc {
	return LogHandler(CorsHandler(DBHandler(f)))
}
