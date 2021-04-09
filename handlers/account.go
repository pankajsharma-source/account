package handlers

import (
	"context"
	"fmt"
	"github.com/pankajsharma-source/user-profile/data"
	"log"
	"net/http"
)

// Products is a http.Handler
type User struct {
	l *log.Logger
}

type KeyUser struct{}

// NewProducts creates a products handler with the given logger
func NewUser(l *log.Logger) *User {
	return &User{l}
}

// getProducts returns the products from the data store
func (u *User) GetUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET User")

	keys, ok := r.URL.Query()["id"] // was "id" before
	if !ok || len(keys[0]) < 1 {
		u.l.Println("Url Param 'key' is missing")
		return
	}
	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	user, err := data.GetUser(key)

	if err != nil {
		http.Error(rw, "User Not Found", http.StatusBadRequest)
		return
	}
	err = user.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal user to Json", http.StatusInternalServerError)
	}
}

func (u *User) AddUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST User")
	user := r.Context().Value(KeyUser{}).(data.User)
	data.AddUser(&user)
}

func (u User) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := data.User{}

		err := user.FromJSON(r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Error reading user", http.StatusBadRequest)
			return
		}

		// validate the product
		err = user.Validate()
		if err != nil {
			u.l.Println("[ERROR] validating user", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating user: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)

	})
}
