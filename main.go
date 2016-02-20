package main

import (
	ae "appengine"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

// error response contains everything we need to use http.Error
type handlerError struct {
	Error   error
	Message string
	Code    int
}

type IdGetter func(ae.Context, int64) (interface{}, error)

func createIdGetterHandler(getter IdGetter) handler {
	return handler(func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
		id := mux.Vars(r)["id"]
		log.Printf("get all somethings for id: %v", id)

		intid, e := strconv.ParseInt(id, 10, 64)
		if e != nil {
			return nil, &handlerError{e, "id should be an integer", http.StatusBadRequest}
		}

		things, e := getter(ae.NewContext(r), intid)
		if e != nil {
			return nil, &handlerError{e, "failed to get all somethings", 1000}
		}
		return things, nil
	})
}

// a custom type that we can use for handling errors and formatting responses
type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

func init() {
	log.Print("init")

	r := mux.NewRouter()
	
	http.Handle("/", r)

}
// attach the standard ServeHTTP method to our handler so the http library can call it
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(w, r)

	// check for errors
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		log.Printf("ERROR: response from method is nil\n")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// send the response and log
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}