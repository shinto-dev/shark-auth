package main

import (
	"net/http"

	"shark-auth/apis"
)

func main() {
	http.Handle("/", apis.Router())
	http.ListenAndServe(":8080", nil)
}
