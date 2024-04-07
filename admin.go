package main

import (
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT user_id, username, user_lastname," +
	"email, registration_date, role FROM Users_partitioned_view")
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	// Обработка полученных данных
	var users []UserPartition
	for rows.Next() {
		var user UserPartition
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.UserLastname,
			&user.Email,
			&user.RegistrationDate,
			&user.Role,
		)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		users = append(users, user)
	}

	// Отображение данных на странице
	tmpl, err := template.ParseFiles("/home/sergey/site2/templates/admin_board.html")
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Передача данных в шаблон для отображения
	err = tmpl.Execute(w, users)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
	}
	/*tmpl, err := template.ParseFiles("/home/sergey/site2/templates/admin_board.html")
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// Отображаем страницу
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
	}*/
}


// Пример добавления админа
/*email := "admin@example.com"
password := "securepassword"

// Хеширование пароля с использованием bcrypt
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
	log.Fatal(err)
}

// Вставка данных нового админа в таблицу
_, err = db.Exec("INSERT INTO Admins (email, password_hash) VALUES ($1, $2)", email, string(hashedPassword))
if err != nil {
	log.Fatal(err)
}
fmt.Println("Admin added successfully")*/