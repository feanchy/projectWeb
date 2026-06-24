package main // package (пакет) — группа файлов; main — точка входа, из неё запускается программа

import ( // import (импорт) — подключение внешних библиотек и пакетов
	"database/sql"     // database/sql — стандартная библиотека для работы с SQL-базами данных
	"fmt"              // fmt — форматированный вывод в консоль (Println, Printf и т.д.)
	"html/template"    // html/template — движок HTML-шаблонов с автоматическим экранированием
	"net/http"         // net/http — стандартная библиотека для HTTP-сервера и клиента
	"projectWeb/handlers" // projectWeb/handlers — наш локальный пакет с обработчиками запросов

	_ "modernc.org/sqlite" // _ (blank import) — импорт только ради side-effect: регистрация драйвера SQLite
)

var db *sql.DB // var (переменная) db (имя) *sql.DB (указатель на объект подключения к базе данных)

func main() { // func (функция) main (главная) — выполняется при запуске программы

	var err error // var err error — объявляем переменную err типа error для хранения ошибок
	db, err = sql.Open("sqlite", "./database/app.db") // sql.Open — открыть соединение; "sqlite" — драйвер; путь к файлу БД
	if err != nil { // if (если) err != nil (ошибка не пустая) — проверка, что Open не вернул ошибку
		panic(err) // panic — аварийная остановка программы с выводом ошибки
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS users ( // sqlStmt — строка SQL; := присваивание; CREATE TABLE — создать таблицу
	id INTEGER PRIMARY KEY AUTOINCREMENT, // id — столбец; INTEGER — целое число; PRIMARY KEY — первичный ключ; AUTOINCREMENT — авто-номер
	login TEXT UNIQUE, // login — столбец логина; TEXT — текст; UNIQUE — значение должно быть уникальным
	password TEXT);` // password — столбец пароля; TEXT — текстовый тип

	_, err = db.Exec(sqlStmt) // db.Exec — выполнить SQL без возврата строк; _ — игнорируем результат (число затронутых строк)
	if err != nil { // если Exec вернул ошибку (например, нет прав на запись)
		panic(err) // останавливаем программу
	}

	templates := template.Must(template.ParseGlob("templates/*.html")) // ParseGlob — загрузить все .html из папки; Must — panic при ошибке парсинга

	h := &handlers.Handler{ // h — указатель на структуру Handler; & — взять адрес в памяти
		DB:  db,        // DB — поле структуры; передаём подключение к базе
		Tpl: templates, // Tpl — поле для шаблонов; передаём загруженные HTML-шаблоны
	}

	mux := http.NewServeMux() // NewServeMux — создать маршрутизатор HTTP (сопоставляет URL и обработчики)
	fs := http.FileServer(http.Dir("./static")) // FileServer — сервер статических файлов; http.Dir — корневая папка ./static
	mux.Handle("/static/", http.StripPrefix("/static/", fs)) // Handle — зарегистрировать путь; StripPrefix — убрать /static/ из URL при поиске файла

	mux.HandleFunc("/", h.HomeHandler)           // HandleFunc — URL "/" → функция HomeHandler (главная страница)
	mux.HandleFunc("/login", h.LoginHandler)     // "/login" → вход пользователя
	mux.HandleFunc("/register", h.RegisterHandler) // "/register" → регистрация
	mux.HandleFunc("/logout", h.LogoutHandler)   // "/logout" → выход

	fmt.Println("Server started on http://localhost:8080") // Println — вывести сообщение в консоль о запуске сервера
	http.ListenAndServe(":8080", mux) // ListenAndServe — слушать порт 8080 и обрабатывать запросы через mux
}
