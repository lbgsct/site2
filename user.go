package main

import (
	"html/template"
	"log"
	"net/http"
	_ "github.com/lib/pq"
)



func IsValidUser(username, password string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1 AND password_hash= $2", username, password).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/login_page.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	username := r.FormValue("user_name")
	password := r.FormValue("password")

	// Проверка данных пользователя
	valid, err := IsValidUser(username, password)
	if err != nil {
		http.Error(w, "Failed to check credentials", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if valid {
		if username == "admin@mail.ru" && password == "admin" {
			http.Redirect(w, r, "http://localhost:8080/admin/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "http://localhost:8080/dashboard/", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Println("Failed authentication attempt")
	}
}