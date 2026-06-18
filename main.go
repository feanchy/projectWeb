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

	sqlStmt := `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT UNIQUE,
	password TEXT);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	templates := template.Must(template.ParseGlob("templates/*.html"))

	h := &handlers.Handler{
		DB:  db,
		Tpl: templates,
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("/login", h.LoginHandler)
	mux.HandleFunc("/register", h.RegisterHandler)
	mux.HandleFunc("/logout", h.LogoutHandler)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
