package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi :)")
	}))
	http.ListenAndServe(":9000", nil)
}
