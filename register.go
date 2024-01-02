package main

import (
	"html/template"
	"log"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
)


func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/register.html")
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

	// Обработка POST-запроса (обработка формы)
	userName := r.FormValue("user_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Добавление данных в базу данных
	success, err := AddUserToDB(userName, email, password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if success {
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
		return
	} else {
		log.Println("Failed to add user to DB")
	}
}


func AddUserToDB(userName, email, password string) (bool, error) {
	result, err := db.Exec("select add_user($1, $2, $3)", userName, email, password)
	if err != nil {
		return false, err
	}

	// Проверяем количество строк, затронутых операцией вставки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	fmt.Println("Rows affected:", rowsAffected)
	// Возвращаем true, если была вставлена хотя бы одна строка (успешное добавление пользователя)
	return rowsAffected > 0, nil
}