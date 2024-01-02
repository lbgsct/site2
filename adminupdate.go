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

func AdminUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы редактирования)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/update.html")
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

func UpdateFlighs(flight Fligh) (bool, error) {
	query := `
        UPDATE flights
        SET
            airline = ?,
            departure_time = ?,
            arrival_time = ?,
            origin = ?,
            destination = ?,
            ticket_price = ?
        WHERE
            id = ?
    `

	// Предполагается, что у вас есть подключение к базе данных с именем db
	// Замените db на ваше актуальное соединение с базой данных

	result, err := db.Exec(query, flight.Airline, flight.DepartureTime, flight.ArrivalTime, flight.Origin, flight.Destination, flight.TicketPrice, flight.ID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func UpdateAirport(airport Airport) (bool, error) {
	query := `
        UPDATE airports
        SET
            airport_code = ?,
            airport_name = ?,
            location = ?,
            other_characteristics = ?
        WHERE
            id = ?
    `

	// Предполагается, что у вас есть подключение к базе данных с именем db
	// Замените db на ваше актуальное соединение с базой данных

	result, err := db.Exec(query, airport.AirportCode, airport.AirportName, airport.Location, airport.OtherCharacteristics, airport.ID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
