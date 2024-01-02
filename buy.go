package main

import (
	//"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	_ "github.com/lib/pq"
)


func BuyPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/buy_page.html")
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

	// Обработка данных формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Получение данных из формы
	flightID := r.FormValue("flightID")
	passengerName := r.FormValue("passengerName")
	paymentMethod := r.FormValue("paymentMethod")
	username := r.FormValue("username") // Получаем ID пользователя из сессии или откуда-либо еще

	userIDInt, err := getUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Failed to fetch user ID", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	userID := strconv.Itoa(userIDInt)
	// Пример вызова функции добавления бронировки в базу данных
	success, err := AddBookingToDB(userID, flightID, passengerName, paymentMethod)
	if err != nil {
		http.Error(w, "Failed to add booking", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if success {
		// Успешно добавлено
		fmt.Fprintf(w, "Booking added successfully!")
		return
	}

	http.Error(w, "Failed to add booking", http.StatusInternalServerError)
}


func AddBookingToDB(userID, flightID, passengerName, paymentMethod string) (bool, error) {
	// Здесь реализуйте логику добавления бронировки в базу данных, используя переданные параметры

	// Пример вызова хранимой процедуры или SQL-запроса для добавления бронировки
	result, err := db.Exec("select AddBooking($1, $2, $3, $4)", userID, flightID, passengerName, paymentMethod)
	if err != nil {
		return false, err
	}

	// Проверяем количество строк, затронутых операцией вставки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// Возвращаем true, если была вставлена хотя бы одна строка (успешное добавление бронировки)
	return rowsAffected > 0, nil
}