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

func AdminDelete(w http.ResponseWriter, r *http.Request) {

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

		aircrafts, err := GetAllAircrafts(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о самолетах", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		booking, err := GetAllBookings(db)
		if err != nil {
			http.Error(w, "Ошибка при получении данных о бронировании", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/delete.html")
		if err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		data := struct {
			Flights    []Flights
			Airports   []Airports
			Aircrafts  []Aircrafts
			BookingDetail []BookingDetail
		}{
			Flights:    flights,
			Airports:   airports,
			Aircrafts:  aircrafts,
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
	case "deleteFlight":

		FlightName := r.FormValue("flight_name")

		success, err := deleteFlight(FlightName)
		if err != nil {
			http.Error(w, "Ошибка удаления рейса: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка удаления рейса в базе данных")
			http.Error(w, "Ошибка удаления рейса в базе данных", http.StatusInternalServerError)
        	return
		}
	
	case "deleteAirport":
		
		AirportName := r.FormValue("airport_name")

		success, err := deleteAirport(AirportName)
		if err != nil {
			http.Error(w, "Ошибка удаления аэропорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка удаления аэропорта в базе данных")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

    case "deleteAircraft":

		AircraftModel := r.FormValue("aircraft_model")

		success, err := deleteAircraft(AircraftModel)
		if err != nil {
			http.Error(w, "Ошибка удаления самолета: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка удаления самолета в базе данных")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	case "deleteBooking":

		BookingID := r.FormValue("id_booking")

		BookingIDInt, err := strconv.Atoi(BookingID)
		if err != nil {
			http.Error(w, "Ошибка преобразования ID брони в число: "+err.Error(), http.StatusInternalServerError)
			return
		}

		success, err := deleteBooking(BookingIDInt)
		if err != nil {
			http.Error(w, "Ошибка удаления брони: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if success {
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка удаления брони в базе данных")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)

	}
}


func deleteFlight(FlightName string) (bool, error) {
	_, err := db.Exec("CALL DeleteFlight($1)", FlightName)
	if err != nil {
		return false, err
	}
	if err != nil {
        fmt.Println("Ошибка при вызове процедуры DeleteFlight:", err)
        return false, err
    }

    return true, nil
}

func deleteAirport(AirportName string) (bool, error) {
	_, err := db.Exec("CALL DeleteAirport ($1)", AirportName)
	if err != nil {
		return false, err
	}
	if err != nil {
        fmt.Println("Ошибка при вызове процедуры DeleteAirport:", err)
        return false, err
    }

    return true, nil
}

func deleteAircraft(AircraftModel string) (bool, error) {
	_, err := db.Exec("CALL DeleteAircraft ($1)", AircraftModel)
	if err != nil {
		return false, err
	}
	if err != nil {
        fmt.Println("Ошибка при вызове процедуры DeleteAircraft:", err)
        return false, err
    }

    return true, nil
}


func deleteBooking(ID_booking int) (bool, error){
	_, err := db.Exec("CALL delete_booking($1)", ID_booking)
	if err != nil {
		return false, err
	}
	if err != nil {
        fmt.Println("Ошибка при вызове процедуры DeleteAircraft:", err)
        return false, err
    }

    return true, nil
}