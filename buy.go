package main

import (
"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)


func BuyPage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        log.Println("Метод запроса не поддерживается:", r.Method)
        http.Error(w, "Метод запроса не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    session, err := store.Get(r, "session-name")
    if err != nil {
        log.Println("Ошибка получения сессии:", err)
        http.Error(w, "Ошибка получения сессии", http.StatusInternalServerError)
        return
    }

    flightName := r.URL.Query().Get("flightName")
    if flightName == "" {
        log.Println("Параметр 'flightName' не найден")
        http.Error(w, "Параметр 'flightName' не найден", http.StatusBadRequest)
        return
    }

    seatNumber := r.URL.Query().Get("seat")
    if flightName == "" {
        log.Println("Параметр 'seat' не найден")
        http.Error(w, "Параметр 'seat' не найден", http.StatusBadRequest)
        return
    }

    IDticket, _ := getIDticketBD(seatNumber, flightName)

    flightSeatPrices, err := GetFlightSeatsAndPrices(flightName)
    if err != nil {
        log.Println("Ошибка получения данных о полетах и ценах:", err)
        http.Error(w, "Ошибка получения данных о полетах и ценах", http.StatusInternalServerError)
        return
    }

    flightVIPSeatPrices, err := GetFlightVIPSeatsAndPrices(flightName)
    if err != nil {
        log.Println("Ошибка получения данных о полетах и ценах:", err)
        http.Error(w, "Ошибка получения данных о полетах и ценах", http.StatusInternalServerError)
        return
    }

    data := struct {
        SeatPrices     []FlightSeatPrice
        VipSeatPrices  []FlightVIPSeatPrice
    }{
        SeatPrices:     flightSeatPrices,
        VipSeatPrices:  flightVIPSeatPrices,
    }

    fmt.Println("Flight Seat Prices:")
    for _, price := range data.SeatPrices {
        fmt.Printf("Flight Name: %s, Seat Number: %s, Ticket Price: %f\n", price.FlightName, price.SeatNumber, price.TicketPrice)
    }

    tmpl := template.New("buy_page.html").Funcs(template.FuncMap{
        "json": func(v interface{}) string {
            jsonData, err := json.Marshal(v)
            if err != nil {
                return err.Error() // Можно обработать ошибку в удобном для вас формате
            }
            return string(jsonData)
        },
    })

    tmpl, err = tmpl.ParseFiles("/home/sergey/site2/templates/buy_page.html")
    if err != nil {
        log.Println("Ошибка загрузки шаблона:", err)
        http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        log.Println("Ошибка отображения данных в шаблоне:", err)
        http.Error(w, "Ошибка отображения данных в шаблоне", http.StatusInternalServerError)
        return
    }

    email, ok := session.Values["email"].(string)
    if !ok {
        log.Println("Email не найден в сессии")
        http.Error(w, "Email не найден в сессии", http.StatusInternalServerError)
        return
    }

    IDuser, _ := getUserIDByEmail(email)

    success, err := AddBookingToDB(IDuser, IDticket)
	if err != nil {
		log.Println("Ошибка при добавлении бронирования:", err)
		http.Error(w, "Ошибка при добавлении бронирования", http.StatusInternalServerError)
		return
	}

	// Если успешно добавлено, можно произвести перенаправление на другую страницу или сообщить об успешной покупке
	if success {
		http.Redirect(w, r, "/dashboard/", http.StatusSeeOther)
		return
	}

}
func GetFlightSeatsAndPrices(flightName string) ([]FlightSeatPrice, error) {
    rows, err := db.Query("SELECT * FROM FlightSeatsAndPrices WHERE flight_name = $1", flightName)
    if err != nil {
        log.Println("Ошибка выполнения SQL-запроса:", err)
        return nil, err
    }
    defer rows.Close()

    var flightSeatPrices []FlightSeatPrice

    for rows.Next() {
        var flightSeatPrice FlightSeatPrice
        err := rows.Scan(&flightSeatPrice.FlightName, &flightSeatPrice.SeatNumber, &flightSeatPrice.TicketPrice)
        if err != nil {
            log.Println("Ошибка при сканировании данных результата запроса:", err)
            return nil, err
        }
        flightSeatPrices = append(flightSeatPrices, flightSeatPrice)
    }

    if err = rows.Err(); err != nil {
        log.Println("Ошибка при чтении результатов запроса:", err)
        return nil, err
    }

    return flightSeatPrices, nil
}


func GetFlightVIPSeatsAndPrices(flightName string) ([]FlightVIPSeatPrice, error) {
    rows, err := db.Query("SELECT * FROM FlightVIPSeatsAndPrices WHERE flight_name = $1", flightName)
    if err != nil {
        log.Println("Ошибка выполнения SQL-запроса:", err)
        return nil, err
    }
    defer rows.Close()

    var flightVIPSeatPrices []FlightVIPSeatPrice

    for rows.Next() {
        var flightVIPSeatPrice FlightVIPSeatPrice
        err := rows.Scan(&flightVIPSeatPrice.FlightName, &flightVIPSeatPrice.SeatNumber, &flightVIPSeatPrice.TicketPrice, &flightVIPSeatPrice.MealChoice, &flightVIPSeatPrice.DrinkAlkohol, &flightVIPSeatPrice.PersonalConcierge)
        if err != nil {
            log.Println("Ошибка при сканировании данных результата запроса:", err)
            return nil, err
        }
        flightVIPSeatPrices = append(flightVIPSeatPrices, flightVIPSeatPrice)
    }

    if err = rows.Err(); err != nil {
        log.Println("Ошибка при чтении результатов запроса:", err)
        return nil, err
    }

    return flightVIPSeatPrices, nil
}

func getUserIDByEmail(email string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return 0, err
	}
    log.Printf("Айдишник юзера: %d", userID)
	return userID, nil
}

func getIDticketBD(seatNumber, flightName string) (int, error) {
	var ticketID int
	err := db.QueryRow("SELECT ticket_id FROM tickets WHERE seat_number = $1 and flight_name = $2", seatNumber, flightName).Scan(&ticketID)
	if err != nil {
		return 0, err
	}
    log.Printf("Айдишник билета: %d", ticketID)
	return ticketID, nil
}

func AddBookingToDB(userID, ticketID int) (bool, error) {
    log.Printf("UserID: %d, TicketID : %d", userID, ticketID)

    _, err := db.Exec("CALL InsertBookin($1, $2)", userID, ticketID)
    if err != nil {
        log.Println("Ошибка при выполнении InsertBookin:", err)
        return false, err
    }
    return true, nil // Возвращаем true при успешном добавлении, иначе false
}
