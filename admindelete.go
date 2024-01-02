package main

import (
	//"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	_ "github.com/lib/pq"
)

func AdminDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы редактирования)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/delete.html")
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
	case "editFlight":
		r.ParseForm()
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
		airportIDStr := r.FormValue("airport_id")
		airportID, err := strconv.Atoi(airportIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		flight := Fligh{
			ID:            airportID,
			Airline:       r.FormValue("airline"),
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			Origin:        r.FormValue("origin"),
			Destination:   r.FormValue("destination"),
			TicketPrice:   ticketPrice,
		}

		success, err := UpdateFlighs(flight)
		if err != nil {
			http.Error(w, "Failed to edit flight", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
			return
		} else {
			log.Println("Failed to edit flight in DB")
		}
	case "editAirport":

		airportIDStr := r.FormValue("airport_id")
		airportID, err := strconv.Atoi(airportIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		airport := Airport{
			ID:                   airportID,
			AirportCode:          r.FormValue("airport_code"),
			AirportName:          r.FormValue("airport_name"),
			Location:             r.FormValue("location"),
			OtherCharacteristics: r.FormValue("other_characteristics"),
		}

		success, err := UpdateAirport(airport)
		if err != nil {
			http.Error(w, "Failed to edit airport", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			fmt.Fprintf(w, "Airport edited successfully!")
			return
		} else {
			log.Println("Failed to edit airport in DB")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)

	}
}


func deleteFlight(flightID int) (bool, error) {
	res, err := db.Exec("DELETE FROM Airports WHERE flight_id = $1", flightID)
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

func deleteAirport(airportCode string) (bool, error) {
	res, err := db.Exec("DELETE FROM Airports WHERE airport_code = $1", airportCode)
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