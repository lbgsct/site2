package main

import (
	//"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"fmt"

	_ "github.com/lib/pq"
)


func AdminAdd (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { 
		// Обработка GET-запроса (отображение формы)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/add.html")
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
	r.ParseForm()
	action := r.FormValue("action")
	switch action {
	case "addFlight":
		r.ParseForm()
		//action := r.FormValue("action")
		departureTimeString := r.FormValue("departure_time")
		departureTime, err := time.Parse("2006-01-02T15:04", departureTimeString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		arrivalTimeString := r.FormValue("arrival_time")
		arrivalTime, err := time.Parse("2006-01-02T15:04", arrivalTimeString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ticketPriceStr := r.FormValue("ticket_price")
		ticketPrice, err := strconv.ParseFloat(ticketPriceStr, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flight := Flights{
			Airline:        r.FormValue("airline"),
			departure_time: departureTime,
			arrival_time:   arrivalTime,
			origin:         r.FormValue("origin"),
			destination:    r.FormValue("destination"),
			ticket_price:   ticketPrice,
		}

		success, err := AddFlight(flight)
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
	case "addAirport":
		airport := Airports{
			airport_code:          r.FormValue("airport_code"),
			airport_name:          r.FormValue("airport_name"),
			location:              r.FormValue("location"),
			other_characteristics: r.FormValue("other_characteristics"),
		}

		success, err := AddAirport(airport)
		if err != nil {
			http.Error(w, "Failed to add airport", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			fmt.Fprintf(w, "Airport added successfully!")
			return
		} else {
			log.Println("Failed to add airport to DB")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)

	}
}

func AddFlight(flight Flights) (bool, error) {
	result, err := db.Exec("select AddNewFlight($1, $2, $3, $4, $5, $6)",
		flight.Airline, flight.departure_time, flight.arrival_time, flight.origin, flight.destination, flight.ticket_price)
	if err != nil {
		fmt.Println("ХУЛИ НЕ РАБОТАЕТ")
		return false, err
	}

	// Проверяем количество строк, затронутых операцией вставки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("ПИЗДИЩА")
		return false, err
	}
	fmt.Println("Rows affected:", rowsAffected)
	// Возвращаем true, если была вставлена хотя бы одна строка (успешное добавление пользователя)
	return rowsAffected > 0, nil
}

func AddAirport(airport Airports) (bool, error) {
	res, err := db.Exec("INSERT INTO Airports (airport_code, airport_name, location, other_characteristics) VALUES ($1, $2, $3, $4)",
		airport.airport_code, airport.airport_name, airport.location, airport.other_characteristics)
	if err != nil {
		return false, err
	}
	// Проверяем количество строк, затронутых операцией вставки
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	fmt.Println("Rows affected:", rowsAffected)
	// Возвращаем true, если была вставлена хотя бы одна строка (успешное добавление пользователя)
	return rowsAffected > 0, nil
}