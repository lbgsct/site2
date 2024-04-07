package main

import (
    "database/sql"
    //"time"
)


func GetAllFlights(db *sql.DB) ([]Flights, error) {
    // Создаем срез для хранения данных о рейсах
    var flights []Flights

    // Выполняем запрос к базе данных для получения данных о рейсах
    rows, err := db.Query("SELECT flight_name, airline, aircraft_model, airport_departure, airport_arrival," + 
	" departure_time,  arrival_time, destination FROM flights")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Читаем данные из результата запроса и добавляем их в срез flights
    for rows.Next() {
        var flight Flights

        // Сканируем данные из результата запроса в переменные структуры Flights
        err := rows.Scan(
            &flight.FlightName,
            &flight.Airline,
            &flight.AircraftModel,
            &flight.DepartureAirport,
            &flight.ArrivalAirport,
            &flight.DepartureTime,
            &flight.ArrivalTime,
            &flight.Destination,
        )
        if err != nil {
            return nil, err
        }

        // Добавляем данные о рейсе в срез flights
        flights = append(flights, flight)
    }

    // Проверяем наличие ошибок при обходе результатов запроса
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Возвращаем данные о рейсах в виде среза []Flights
    return flights, nil
}


func GetAllAirports(db *sql.DB) ([]Airports, error) {
    var airports []Airports

    // Выполняем запрос к базе данных для получения данных об аэропортах
    rows, err := db.Query("SELECT airport_name, airport_city FROM Airports")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Читаем данные из результата запроса и добавляем их в срез airports
    for rows.Next() {
        var airport Airports

        // Сканируем данные из результата запроса в переменные структуры Airports
        err := rows.Scan(
            &airport.AirportName,
            &airport.AirportCity,
        )
        if err != nil {
            return nil, err
        }

        // Добавляем данные об аэропорте в срез airports
        airports = append(airports, airport)
    }

    // Проверяем наличие ошибок при обходе результатов запроса
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Возвращаем данные об аэропортах в виде среза []Airports
    return airports, nil
}

func GetAllTickets(db *sql.DB) ([]Tickets, error) {
    var tickets []Tickets

    // Выполняем запрос к базе данных для получения данных о билетах
    rows, err := db.Query("SELECT flight_name, ticket_price, seat_number FROM tickets")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Читаем данные из результата запроса и добавляем их в срез tickets
    for rows.Next() {
        var ticket Tickets

        // Сканируем данные из результата запроса в переменные структуры Tickets
        err := rows.Scan(
            &ticket.FlightName,
            &ticket.Price,
            &ticket.SeatNumber,
        )
        if err != nil {
            return nil, err
        }

        // Добавляем данные о билете в срез tickets
        tickets = append(tickets, ticket)
    }

    // Проверяем наличие ошибок при обходе результатов запроса
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Возвращаем данные о билетах в виде среза []Tickets
    return tickets, nil
}

func GetAllVipTickets(db *sql.DB) ([]VipTickets, error) {
    var viptickets []VipTickets

    // Выполняем запрос к базе данных для получения данных о билетах
    rows, err := db.Query("SELECT flight_name, ticket_price, seat_number, meal_choice,"+ 
	"drink_alcohol_choice, personal_concierge FROM VipTickets")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Читаем данные из результата запроса и добавляем их в срез tickets
    for rows.Next() {
        var vipticket VipTickets

        // Сканируем данные из результата запроса в переменные структуры Tickets
        err := rows.Scan(
            &vipticket.FlightName,
            &vipticket.Price,
            &vipticket.SeatNumber,
			&vipticket.MealChoice,
			&vipticket.DrinkAlkohol,
			&vipticket.PersonalConcierge,
        )
        if err != nil {
            return nil, err
        }

        // Добавляем данные о билете в срез tickets
        viptickets = append(viptickets, vipticket)
    }

    // Проверяем наличие ошибок при обходе результатов запроса
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Возвращаем данные о билетах в виде среза []Tickets
    return viptickets, nil
}

func GetAllAircrafts(db *sql.DB) ([]Aircrafts, error) {
    var aircrafts []Aircrafts

    // Выполняем запрос к базе данных для получения данных о билетах
    rows, err := db.Query("SELECT aircraft_model, total_seats, aircraft_condition FROM aircrafts")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Читаем данные из результата запроса и добавляем их в срез tickets
    for rows.Next() {
        var aircraft Aircrafts

        // Сканируем данные из результата запроса в переменные структуры Tickets
        err := rows.Scan(
            &aircraft.AircraftModel,
            &aircraft.TotalSeats,
            &aircraft.AircraftCondition,
        )
        if err != nil {
            return nil, err
        }

        // Добавляем данные о билете в срез tickets
        aircrafts = append(aircrafts, aircraft)
    }

    // Проверяем наличие ошибок при обходе результатов запроса
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Возвращаем данные о билетах в виде среза []Tickets
    return aircrafts, nil
}

func GetAllBookings(db *sql.DB) ([]BookingDetail, error) {
    var bookings []BookingDetail

    // Выполняем запрос к базе данных для получения данных о бронированиях
    rows, err := db.Query("SELECT * FROM Bookingview")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Читаем данные из результата запроса и добавляем их в срез bookings
    for rows.Next() {
        var booking BookingDetail

        // Сканируем данные из результата запроса в переменные структуры BookingDetail
        err := rows.Scan(
            &booking.BookingID,
            &booking.UserID,
            &booking.TicketID,
            &booking.FlightName,
            &booking.TicketPrice,
            &booking.SeatNumber,
            &booking.BookingStatus,
        )
        if err != nil {
            return nil, err
        }

        // Добавляем данные о бронировании в срез bookings
        bookings = append(bookings, booking)
    }

    // Проверяем наличие ошибок при обходе результатов запроса
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Возвращаем данные о бронированиях в виде среза []BookingDetail
    return bookings, nil
}