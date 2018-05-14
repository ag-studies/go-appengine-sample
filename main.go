package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	t := taskqueue.NewPOSTTask("/worker", map[string][]string{"key": {"val"}})
	_, err := taskqueue.Add(ctx, t, "")
	if err != nil {
		log.Errorf(ctx, "Failed to add task", err.Error())
	}
	fmt.Fprintf(w, "Handler Success")
	log.Infof(ctx, "Handler succeeded")
}

func Worker(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Worker succeeded")
}

func init() {
	http.HandleFunc("/worker", Worker)
	http.HandleFunc("/", Handler)
}
