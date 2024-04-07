package main

import (
	"html/template"
	"net/http"
    "database/sql"
    "log"

	//"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func GetUserProfileBD(email string) (*User, error) {

	query := "SELECT username, user_lastname, email FROM Users WHERE email = $1"
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.UserName, &user.UserLastname, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Profile(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Ошибка получения сессии: %v", err)
        return
    }

    email, ok := session.Values["email"].(string)
    if !ok {
        http.Error(w, "Email не найден в сессии", http.StatusInternalServerError)
        log.Println("Email не найден в сессии")
        return
    }

    userProfile, err := GetUserProfileBD(email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Ошибка получения профиля пользователя: %v", err)
        return
    }

    userID, err := getUserIDByEmail(email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Ошибка получения ID пользователя: %v", err)
        return
    }

    userBookings, err := GetUserBookingsByID(db, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Ошибка получения бронирований пользователя: %v", err)
        return
    }

    data := struct {
        UserName     string
        UserLastname string
        Email        string
        Bookings     []BookingDetail
    }{
        UserName:     userProfile.UserName,
        UserLastname: userProfile.UserLastname,
        Email:        userProfile.Email,
        Bookings:     userBookings,
    }

    tmpl, err := template.ParseFiles("/home/sergey/site2/templates/profile.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Ошибка парсинга шаблона: %v", err)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Ошибка выполнения шаблона: %v", err)
        return
    }
}

func GetUserBookingsByID(db *sql.DB, userID int) ([]BookingDetail, error) {
    log.Printf("Начало GetUserBookingsByID для пользователя с ID: %d", userID)

    var bookings []BookingDetail

    rows, err := db.Query("SELECT DISTINCT b.booking_id, b.user_id, b.id_ticket, t.flight_name, t.ticket_price,"+
        " t.seat_number, b.booking_status FROM Bookings b INNER JOIN Tickets t ON b.id_ticket = t.ticket_id WHERE b.user_id = $1", userID)
    if err != nil {
        log.Printf("Ошибка выполнения запроса к базе данных: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var booking BookingDetail

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
            log.Printf("Ошибка сканирования данных из результата запроса: %v", err)
            return nil, err
        }

        bookings = append(bookings, booking)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Ошибка при обработке результатов запроса: %v", err)
        return nil, err
    }

    log.Printf("Завершение GetUserBookingsByID для пользователя с ID: %d", userID)

    return bookings, nil
}

