package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/register.html")
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	// Обработка POST-запроса (обработка формы)
	userName := r.FormValue("user_name")
	userLastname := r.FormValue("user_lastname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Добавление данных в базу данных
	err = AddUserToDB(userName, userLastname, email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Не удалось зарегистрировать пользователя", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// Если мы дошли до этой точки, пользователь успешно добавлен в базу данных
	http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
}

func AddUserToDB(userName, userLastname, email, password string) error {
	// Проверка существования пользователя с таким email
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Пользователь с такими данными уже существует")
	}

	_, err = db.Exec("CALL InsertUser($1, $2, $3, $4)", userName, userLastname, email, password)
	if err != nil {
		return err
	}
	return nil
}
