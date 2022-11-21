package main

import (
	"context"
	"net/http"
)

// UserKey is what's used to store the user details in a request context
type contextKey int

const userKey contextKey = iota

// authCheckMiddleware for checking if a user name is passed in this request
func authCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("user")

		// skip authentication check on root route
		if r.URL.Path == "/" && username == "" {
			next.ServeHTTP(w, r)
			return
		}

		// if username is not passed in the request, return an error
		if username == "" {
			http.Error(w, "user not passed in query params", http.StatusUnauthorized)
			return
		}

		// mimicking auth token: check for user in question
		u := User{}
		err := db.QueryRow("SELECT id, name FROM users WHERE name = $1;", username).Scan(&u.ID, &u.Name)
		if err != nil {
			http.Error(w, "username is incorrect", http.StatusInternalServerError)
			return
		}

		if u.ID == 0 {
			http.Error(w, "username is incorrect", http.StatusUnauthorized)
			return
		}

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), userKey, &u)

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
