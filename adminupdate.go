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
		flights, err := GetAllFlights(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о рейсах", http.StatusInternalServerError)
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

		booking, err := GetAllBookings(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о бронировании", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/update.html")
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		
		data := struct {
			Flights   []Flights
			Tickets   []Tickets
			VipTickets []VipTickets
			BookingDetail []BookingDetail
		}{
			Flights:  flights,
			Tickets: tickets,
			VipTickets: viptickets,
			BookingDetail: booking,
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
	case "updateFlight":
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Ошибка при разборе данных формы:", err)
			http.Error(w, "Ошибка при разборе данных формы", http.StatusInternalServerError)
			return
		}
		
		departureTimeString := r.FormValue("departureTime")
		departureTime, err := time.Parse("2006-01-02T15:04", departureTimeString)
		if err != nil {
			fmt.Println("Ошибка при парсинге времени прибытия:", err)
			http.Error(w, "Ошибка при парсинге времени прибытия", http.StatusInternalServerError)
			return
		}
		departureTime = departureTime.Add(-3 * time.Hour)

		// Получение времени прибытия из формы и преобразование в тип time.Time
		arrivalTimeString := r.FormValue("arrivalTime")
		arrivalTime, err := time.Parse("2006-01-02T15:04", arrivalTimeString)
		if err != nil {
			fmt.Println("Ошибка при парсинге времени прибытия:", err)
			http.Error(w, "Ошибка при парсинге времени прибытия", http.StatusInternalServerError)
			return
		}
		arrivalTime = arrivalTime.Add(-3 * time.Hour)
		// Создание объекта Flights из данных формы
		updateFlight := Flights{
			FlightName:       r.FormValue("flightName"),
			DepartureTime:    departureTime,
			ArrivalTime:   	  arrivalTime,
		}
	
		// Добавление рейса в базу данных
		updated, err := UpdateFlight(updateFlight)
		if err != nil {
			fmt.Println("Ошибка при изменении рейса:", err)
			http.Error(w, "Ошибка при изменении рейса", http.StatusInternalServerError)
			return
		}
	
		if updated {
			http.Redirect(w, r, "/admin/update/", http.StatusSeeOther)
		} else {
			fmt.Fprintln(w, "Не удалось изменить рейс.")
		}

	case "updateTickets":
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Ошибка парсинга цены билета", http.StatusBadRequest)
			log.Println(err)
			return
		}

		updateticket := Tickets{
			FlightName: r.FormValue("flight_name"),
			Price: price,
			SeatNumber: r.FormValue("seat_number"),
			}

		success, err := UpdateTicketPrice(updateticket)
		if err != nil {
			http.Error(w, "Ошибка изменения билета", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/update/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка изменения билета в базе данных")
		}

	case "updateVipTickets":
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
			FlightName: r.FormValue("flight_name"),
			Price: price,
			SeatNumber: r.FormValue("seat_number"),
			MealChoice: r.FormValue("meal_choice"),
			DrinkAlkohol: r.FormValue("drink_alcohol_choice"),
			PersonalConcierge: personalConcierge,
		}

		success, err := UpdateVipTicket(vipticket)
		if err != nil {
			http.Error(w, "Ошибка изменения vip-билета", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/update/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка изменения vip-билета в базе данных")
		}
	case "updateBooking":
		bookingID := r.FormValue("booking_id")
		bookingStatus := r.FormValue("booking_status")

		// Преобразование bookingID в int
		bookingIDInt, err := strconv.Atoi(bookingID)
		if err != nil {
			http.Error(w, "Ошибка преобразования ID брони в число: "+err.Error(), http.StatusInternalServerError)
			return
		}

		booking := BookingDetail{
			BookingID:      bookingIDInt,
			BookingStatus:  bookingStatus,
		}

		success, err := UpdateBookingStatus(booking)
		if err != nil {
			http.Error(w, "Ошибка обновления статуса брони: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if success {
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка обновления статуса брони в базе данных")
		}
	
	
	default:
		http.Error(w, "Переход в дефолт", http.StatusBadRequest)

	}
}

func UpdateFlight(flight Flights) (bool, error) {
    // Вызываем хранимую процедуру с переданными аргументами
    _, err := db.Exec("CALL UpdateFlight($1, $2, $3)",
        flight.FlightName, flight.DepartureTime, flight.ArrivalTime,)
    if err != nil {
        fmt.Println("Ошибка при вызове процедуры UpdateFlight:", err)
        return false, err
    }

    return true, nil
}

func UpdateTicketPrice(ticket Tickets) (bool, error) {
	_, err := db.Exec("CALL UpdateTicketPrice ($1, $2, $3)", ticket.FlightName, ticket.SeatNumber, ticket.Price)
	if err != nil {
		return false, err
	}
	if err != nil {
        fmt.Println("Ошибка при вызове процедуры UpdateTicketPrice:", err)
        return false, err
    }

    return true, nil
}

func UpdateVipTicket(vipticket VipTickets) (bool, error) {
	_, err := db.Exec("CALL  UpdateVIPTicket($1, $2, $3, $4, $5, $6)", vipticket.FlightName,vipticket.SeatNumber, vipticket.Price, vipticket.MealChoice, vipticket.DrinkAlkohol, vipticket.PersonalConcierge)
	if err != nil {
		return false, err
	}
	if err != nil {
        fmt.Println("Ошибка при вызове процедуры InsertTickets:", err)
        return false, err
    }

    return true, nil
}

func UpdateBookingStatus(booking BookingDetail) (bool, error) {
	_, err := db.Exec("CALL update_booking_status($1, $2)", booking.BookingID, booking.BookingStatus)
	if err != nil {
		fmt.Println("Error calling update_booking_status procedure:", err)
		return false, err
	}

	return true, nil
}


