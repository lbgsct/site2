package main

import (
	"database/sql"
	"html/template"
		"log"
	"net/http"
	"strconv"
    "time"
	_ "github.com/lib/pq"
)


func DashboardPage(w http.ResponseWriter, r *http.Request) {
	var rows *sql.Rows
	var err error
	
	// Получаем параметры запроса
	queryParams := QueryParams{
		Departure: r.FormValue("departure"),
		Arrival:   r.FormValue("arrival"),
		Date:      r.FormValue("date"),
	}

	// Флаг, указывающий, нужно ли выполнить поиск или нет
	doSearch := queryParams.Departure != "" || queryParams.Arrival != "" || queryParams.Date != ""

	var query string
	var args []interface{}

	if doSearch {
		query = "SELECT * FROM UserFlights WHERE"
		args = []interface{}{}

		// Если есть Departure в запросе
		if queryParams.Departure != "" {
			query += " departure_city = $1"
			args = append(args, queryParams.Departure)
		}

		// Если есть Arrival в запросе
		if queryParams.Arrival != "" {
			if len(args) > 0 {
				query += " AND"
			}
			query += " arrival_city = $" + strconv.Itoa(len(args)+1)
			args = append(args, queryParams.Arrival)
		}

		// Если есть Date в запросе
		if queryParams.Date != "" {
			parsedDate, err := time.Parse("2006-01-02T15:04", queryParams.Date)
			if err != nil {
				http.Error(w, "Неправильный формат даты", http.StatusBadRequest)
				log.Println(err)
				return
			}
			formattedDate := parsedDate.Format("2006-01-02 15:04") // Форматируем дату до минут
			if len(args) > 0 {
				query += " AND"
			}
			query += " departure_date = $" + strconv.Itoa(len(args)+1)
			args = append(args, formattedDate)
		}

		rows, err = db.Query(query, args...)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()
	} else {
		rows, err = db.Query("SELECT * FROM UserFlights")
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()
	}

	var userFlights []UserFlights

	for rows.Next() {
		var flight UserFlights
		var minTicketPrice sql.NullFloat64 
		err := rows.Scan(
			&flight.FlightName,
			&flight.Airline,
			&flight.TotalSeats,
			&flight.AircraftModel,
			&flight.DepartureCity,
			&flight.DepartureAirport,
			&flight.DepartureTime,
			&flight.ArrivalCity,
			&flight.ArrivalAirport,
			&flight.ArrivalTime,
			&minTicketPrice,
		)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if minTicketPrice.Valid {
            flight.MinTicketPrice = minTicketPrice.Float64
        } else {
            flight.MinTicketPrice = 0 // Устанавливаем значение по умолчанию для NULL
        }

		userFlights = append(userFlights, flight)
	}

	tmpl, err := template.ParseFiles("/home/sergey/site2/templates/dash_board.html")
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, userFlights)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println(err)
	}
}
