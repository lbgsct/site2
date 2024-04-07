package main

import (
	"database/sql"
	"fmt"
	"log"

	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newPassword"
	dbname   = "site"
)

var db *sql.DB // Переменная для хранения подключения к базе данных

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// Установка параметров куки-сессии
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15, // Например, устанавливает сессию на 15 минут
		HttpOnly: true,
	}
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login/", LoginPage)
	http.HandleFunc("/register/", Register)
	http.HandleFunc("/dashboard/", DashboardPage)
	http.HandleFunc("/dashboard/logout/", Logout)
	http.HandleFunc("/profile/", Profile)
	http.HandleFunc("/profile/edit/", Edit)
	http.HandleFunc("/admin/", AdminPage)
	http.HandleFunc("/buy/", BuyPage)
	http.HandleFunc("/admin/add/", AdminAdd)
	http.HandleFunc("/admin/update/", AdminUpdate)
	http.HandleFunc("/admin/delete/", AdminDelete)
	http.HandleFunc("/admin/logout/", Logout)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
