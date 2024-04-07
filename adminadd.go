package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func AdminAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		flights, err := GetAllFlights(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о рейсах", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		airports, err := GetAllAirports(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных об аэропортах", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		tickets, err := GetAllTickets(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о билетах", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		viptickets, err := GetAllVipTickets(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о VIP билетах", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		aircrafts, err := GetAllAircrafts(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о VIP билетах", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Отображение страницы с данными о рейсах и аэропортах
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/add.html")
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		data := struct {
			Flights    []Flights
			Airports   []Airports
			Tickets    []Tickets
			VipTickets []VipTickets
			Aircrafts  []Aircrafts
		}{
			Flights:    flights,
			Airports:   airports,
			Tickets:    tickets,
			VipTickets: viptickets,
			Aircrafts:  aircrafts,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	r.ParseForm()
	action := r.FormValue("action")
	switch action {
	case "addFlight":

		departureTimeString := r.FormValue("departureTime")
		departureTime, err := time.Parse("2006-01-02T15:04", departureTimeString)
		if err != nil {
			fmt.Println("Ошибка при парсинге времени прибытия:", err)
			http.Error(w, "Ошибка при парсинге времени прибытия", http.StatusInternalServerError)
			return
		}

		// Получение времени прибытия из формы и преобразование в тип time.Time
		arrivalTimeString := r.FormValue("arrivalTime")
		arrivalTime, err := time.Parse("2006-01-02T15:04", arrivalTimeString)
		if err != nil {
			fmt.Println("Ошибка при парсинге времени прибытия:", err)
			http.Error(w, "Ошибка при парсинге времени прибытия", http.StatusInternalServerError)
			return
		}

		// Создание объекта Flights из данных формы
		newFlight := Flights{
			FlightName:       r.FormValue("flightName"),
			Airline:          r.FormValue("airline"),
			AircraftModel:    r.FormValue("aircraftModel"),
			DepartureAirport: r.FormValue("departureAirport"),
			ArrivalAirport:   r.FormValue("arrivalAirport"),
			DepartureTime:    departureTime,
			ArrivalTime:      arrivalTime,
			Destination:      r.FormValue("destination"),
		}

		// Добавление рейса в базу данных
		added, err := AddFlight(newFlight)
		if err != nil {
			fmt.Println("Ошибка при добавлении рейса:", err)
			http.Error(w, "Ошибка при добавлении рейса", http.StatusInternalServerError)
			return
		}

		if added {
			http.Redirect(w, r, "/admin/add/", http.StatusSeeOther)
			return
		} else {
			fmt.Fprintln(w, "Не удалось добавить рейс.")
		}

	case "addAirport":
		airport := Airports{
			AirportName: r.FormValue("airport_name"),
			AirportCity: r.FormValue("airport_city"),
		}

		success, err := AddAirport(airport)
		if err != nil {
			http.Error(w, "Ошибка добавления аэропорта", http.StatusInternalServerError)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/add/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка добавления аэропортa в базу данных")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	case "addVipTickets":
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Ошибка парсинга цены билета", http.StatusBadRequest)
			log.Println(err)
			return
		}

		// Получаем значение из формы
		personalConciergeValue := r.FormValue("personal_concierge")

		// Преобразуем строковое значение в булево
		var personalConcierge bool
		if personalConciergeValue == "yes" {
			personalConcierge = true
		} else {
			personalConcierge = false
		}
		vipticket := VipTickets{
			FlightName:        r.FormValue("flight_name"),
			Price:             price,
			SeatNumber:        r.FormValue("seat_number"),
			MealChoice:        r.FormValue("meal_choice"),
			DrinkAlkohol:      r.FormValue("drink_alcohol_choice"),
			PersonalConcierge: personalConcierge,
		}

		success, err := AddVipTicket(vipticket)
		if err != nil {
			http.Error(w, "Ошибка добавления билета", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/add/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка добавления билета в базу данных")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	case "addTickets":
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Ошибка парсинга цены билета", http.StatusBadRequest)
			log.Println(err)
			return
		}

		ticket := Tickets{
			FlightName: r.FormValue("flight_name"),
			Price:      price,
			SeatNumber: r.FormValue("seat_number"),
		}

		success, err := AddTicket(ticket)
		if err != nil {
			http.Error(w, "Ошибка добавления билета", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/add/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка добавления билета в базу данных")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	case "addAircraft":
		totalSeatsStr := r.FormValue("total_seats")
		totalSeatsInt32, err := strconv.ParseInt(totalSeatsStr, 10, 32)
		if err != nil {
			http.Error(w, "Ошибка парсинга количества мест", http.StatusBadRequest)
			log.Println(err)
			return
		}

		aircraft := Aircrafts{
			AircraftModel:     r.FormValue("aircraft_model"),
			TotalSeats:        int32(totalSeatsInt32),
			AircraftCondition: r.FormValue("aircraft_condition"),
		}

		success, err := AddAircraft(aircraft)
		if err != nil {
			http.Error(w, "Ошибка добавления самолета", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/add/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка добавления самолета в базу данных")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	default:
		http.Error(w, "Переход в дефолт", http.StatusBadRequest)

	}
}

func AddFlight(flight Flights) (bool, error) {
	// Вызываем хранимую процедуру с переданными аргументами
	_, err := db.Exec("CALL InsertFlight($1, $2, $3, $4, $5, $6, $7, $8)",
		flight.FlightName, flight.Airline, flight.AircraftModel, flight.DepartureAirport, flight.ArrivalAirport, flight.DepartureTime, flight.ArrivalTime, flight.Destination)
	if err != nil {
		fmt.Println("Ошибка при вызове процедуры InsertFlight:", err)
		return false, err
	}

	return true, nil
}

func AddAirport(airport Airports) (bool, error) {
	_, err := db.Exec("CALL InsertAirport ($1, $2)", airport.AirportName, airport.AirportCity)
	if err != nil {
		return false, err
	}
	if err != nil {
		fmt.Println("Ошибка при вызове процедуры InsertFlight:", err)
		return false, err
	}

	return true, nil
}

func AddTicket(ticket Tickets) (bool, error) {
	_, err := db.Exec("CALL InsertTickets ($1, $2, $3)", ticket.FlightName, ticket.Price, ticket.SeatNumber)
	if err != nil {
		return false, err
	}
	if err != nil {
		fmt.Println("Ошибка при вызове процедуры InsertTickets:", err)
		return false, err
	}

	return true, nil
}

func AddVipTicket(vipticket VipTickets) (bool, error) {
	_, err := db.Exec("CALL InsertVipTickets ($1, $2, $3, $4, $5, $6)", vipticket.FlightName, vipticket.Price, vipticket.SeatNumber, vipticket.MealChoice, vipticket.DrinkAlkohol, vipticket.PersonalConcierge)
	if err != nil {
		return false, err
	}
	if err != nil {
		fmt.Println("Ошибка при вызове процедуры InsertTickets:", err)
		return false, err
	}

	return true, nil
}

func AddAircraft(aircraft Aircrafts) (bool, error) {
	_, err := db.Exec("CALL AddAircraft ($1, $2, $3)", aircraft.AircraftModel, aircraft.TotalSeats, aircraft.AircraftCondition)
	if err != nil {
		return false, err
	}
	if err != nil {
		fmt.Println("Ошибка при вызове процедуры AddAircraft:", err)
		return false, err
	}

	return true, nil
}
