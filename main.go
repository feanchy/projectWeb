package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Home page")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login page")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
