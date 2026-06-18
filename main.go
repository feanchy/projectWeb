package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"projectWeb/handlers"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("sqlite", "./database/app.db")
	if err != nil {
		panic(err)
	}

	templates := template.Must(template.ParseGlob("templates/*.html"))

	h := &handlers.Handler{
		DB:  db,
		Tpl: templates,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("/login", h.LoginHandler)
	mux.HandleFunc("/register", h.RegisterHandler)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
