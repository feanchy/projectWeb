package handlers // package handlers — пакет с HTTP-обработчиками (handlers = обработчики запросов)

import ( // import — подключение нужных библиотек
	"database/sql"     // database/sql — работа с SQL-базой данных
	"fmt"              // fmt — вывод в консоль (отладочные сообщения)
	"html/template"    // html/template — рендер HTML-шаблонов
	"net/http"         // net/http — типы ResponseWriter, Request, Cookie и HTTP-функции

	"golang.org/x/crypto/bcrypt" // bcrypt — библиотека для безопасного хеширования паролей
)

type Handler struct { // type (тип) Handler (имя) struct (структура) — объект с зависимостями для всех обработчиков
	DB  *sql.DB              // DB — указатель на подключение к базе данных
	Tpl *template.Template   // Tpl — указатель на набор HTML-шаблонов
}

var Templates *template.Template // var Templates — глобальная переменная (сейчас не используется, дублирует h.Tpl)

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) { // func — функция; (h *Handler) — метод структуры Handler; w — ответ клиенту; r — входящий запрос

	cookie, err := r.Cookie("user") // r.Cookie — прочитать cookie с именем "user" из запроса браузера

	if err != nil { // if — если cookie нет (пользователь не залогинен)
		h.Tpl.ExecuteTemplate(w, "index.html", nil) // ExecuteTemplate — отрисовать шаблон index.html; nil — данных нет (гость)
		return // return — выйти из функции, дальше код не выполняется
	}

	h.Tpl.ExecuteTemplate(w, "index.html", cookie.Value) // cookie.Value — логин из cookie; передаём в шаблон для приветствия
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) { // RegisterHandler — обработчик страницы регистрации

	if r.Method == http.MethodPost { // r.Method — HTTP-метод запроса; MethodPost — "POST" (отправка формы)
		r.ParseForm() // ParseForm — разобрать тело POST-запроса и заполнить r.Form

		login := r.FormValue("login")       // FormValue — значение поля формы с name="login"
		password := r.FormValue("password") // FormValue — значение поля формы с name="password"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // GenerateFromPassword — хеш пароля; []byte — байты; DefaultCost — стандартная сложность
		if err != nil { // если хеширование не удалось
			http.Error(w, "hash error", 500) // http.Error — отправить клиенту текст ошибки; 500 — Internal Server Error
			return
		}

		_, err = h.DB.Exec( // Exec — выполнить SQL-команду INSERT
			"INSERT INTO users (login, password) VALUES (?, ?)", // ? — плейсхолдеры (защита от SQL-инъекций)
			login, string(hashedPassword), // login и хеш пароля подставляются как параметры, не как часть SQL
		)

		if err != nil { // ошибка вставки (например, логин уже занят — UNIQUE)
			fmt.Println("user already exists or db error:", err) // Println — вывести ошибку в консоль сервера
			http.Error(w, "User already exists", 400) // 400 — Bad Request (неверный запрос клиента)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect — перенаправить браузер на главную "/"; StatusSeeOther — код 303
		return
	}

	h.Tpl.ExecuteTemplate(w, "register.html", nil) // GET-запрос — показать HTML-форму регистрации
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) { // LoginHandler — обработчик входа в аккаунт

	if r.Method == http.MethodPost { // обрабатываем только отправку формы (POST)
		r.ParseForm() // разбираем поля login и password из тела запроса

		login := r.FormValue("login")       // логин, который ввёл пользователь
		password := r.FormValue("password") // пароль в открытом виде (из формы)

		var dbPassword string // var — объявить переменную; сюда запишем хеш пароля из базы

		err := h.DB.QueryRow("SELECT password FROM users WHERE login = ?", login).Scan(&dbPassword) // QueryRow — одна строка; Scan — записать результат в dbPassword; & — указатель

		if err == sql.ErrNoRows { // ErrNoRows — пользователь с таким login не найден в таблице users
			http.Error(w, "User not found", 401) // 401 — Unauthorized (не авторизован)
			return
		}

		if err != nil { // любая другая ошибка базы данных
			http.Error(w, "DB error", 500) // 500 — ошибка на стороне сервера
			return
		}

		err = bcrypt.CompareHashAndPassword( // CompareHashAndPassword — сравнить хеш из БД с введённым паролем
			[]byte(dbPassword),  // хеш из базы данных (байты)
			[]byte(password),    // пароль из формы (байты)
		)

		if err != nil { // пароли не совпали
			http.Error(w, "wrong password", 401) // неверный пароль
			return
		}

		http.SetCookie(w, &http.Cookie{ // SetCookie — отправить браузеру cookie (сессия «залогинен»)
			Name:  "user",  // Name — имя cookie
			Value: login,   // Value — значение cookie (логин пользователя)
			Path:  "/",     // Path — cookie действует на всех страницах сайта
		})

		fmt.Println("LOGIN SUCCESS:", login) // отладочный вывод успешного входа в консоль

		http.Redirect(w, r, "/", http.StatusSeeOther) // после входа перенаправляем на главную
		return
	}

	h.Tpl.ExecuteTemplate(w, "login.html", nil) // GET-запрос — показать форму входа
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) { // LogoutHandler — выход из аккаунта
	http.SetCookie(w, &http.Cookie{ // перезаписываем cookie, чтобы «разлогинить» пользователя
		Name:   "user",  // то же имя cookie
		Value:  "",      // пустое значение — пользователь больше не залогинен
		Path:   "/",     // путь действия cookie
		MaxAge: -1,      // MaxAge: -1 — удалить cookie у браузера немедленно
	})

	http.Redirect(w, r, "/", http.StatusSeeOther) // перенаправление на главную после выхода
}
