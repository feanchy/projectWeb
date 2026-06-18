package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

type Handler struct {
	DB  *sql.DB
	Tpl *template.Template
}

var Templates *template.Template

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("user")

	if err != nil {
		h.Tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	h.Tpl.ExecuteTemplate(w, "index.html", cookie.Value)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()

		login := r.FormValue("login")
		password := r.FormValue("password")

		_, err := h.DB.Exec(
			"INSERT INTO users (login, password) VALUES (?, ?)",
			login, password,
		)

		if err != nil {
			fmt.Println("DB error", err)
			http.Error(w, "DB error", 500)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tpl.ExecuteTemplate(w, "register.html", nil)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()

		login := r.FormValue("login")
		password := r.FormValue("password")

		var dbPassword string

		err := h.DB.QueryRow("SELECT password FROM users WHERE login = ?", login).Scan(&dbPassword)

		if err != nil {
			fmt.Println("user not found or db error:", err)
			http.Error(w, "Invalid login", 401)
			return
		}

		if password != dbPassword {
			fmt.Println("wrong password")
			http.Error(w, "wrong password", 401)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "user",
			Value: login,
			Path:  "/",
		})

		fmt.Println("LOGIN SUCCESS:", login)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tpl.ExecuteTemplate(w, "login.html", nil)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
