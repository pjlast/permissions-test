package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// here we read from the request context and fetch out `"user"` key set in
	// the MyMiddleware example above.
	rawUser := r.Context().Value(userKey)

	if rawUser != nil {
		user := rawUser.(*User)
		fmt.Fprintf(w, "Hello %s!\n", user.Name)
		return
	}

	fmt.Fprintf(w, "Hello World!\n")
}

func getBatchChangeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
