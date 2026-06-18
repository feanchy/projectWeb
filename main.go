package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "modernc.org/sqlite"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var db *sql.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()

		login := r.FormValue("login")
		password := r.FormValue("password")

		fmt.Println("REGISTER:", login, password)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "register.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()

		login := r.FormValue("login")
		password := r.FormValue("password")

		fmt.Println("LOGIN:", login, password)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "login.html", nil)

}

func main() {
	var err error
	db, err = sql.Open("sqlite", "./database/app.db")
	if err != nil {
		panic(err)
	}

	sqlStmt := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT,
	password TEXT
);
`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
