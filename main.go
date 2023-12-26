package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newPassword"
	dbname   = "site"
)

var db *sql.DB // Переменная для хранения подключения к базе данных

type AvailableFlight struct {
	FlightID         int
	Airline          string
	DepartureTime    time.Time
	ArrivalTime      time.Time
	TicketPrice      float64
	ArrivalAirport   string
	DepartureAirport string
}

type User struct {
	Email    string
	Password string
	UserName string
}

type Flights struct {
	Airline        string
	departure_time time.Time
	arrival_time   time.Time
	origin         string
	destination    string
	ticket_price   float64
}

type Airports struct {
	airport_code          string
	airport_name          string
	location              string
	other_characteristics string
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/login/", loginPage)
	http.HandleFunc("/register/", register)
	http.HandleFunc("/dashboard/", dashboardPage)
	http.HandleFunc("/admin/", adminPage)
	http.HandleFunc("/buy/", buyPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("/home/sergey/site2/templates/home_page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// Отображаем страницу
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func addUserToDB(userName, email, password string) (bool, error) {
	result, err := db.Exec("select add_user($1, $2, $3)", userName, email, password)
	if err != nil {
		return false, err
	}

	// Проверяем количество строк, затронутых операцией вставки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	fmt.Println("Rows affected:", rowsAffected)
	// Возвращаем true, если была вставлена хотя бы одна строка (успешное добавление пользователя)
	return rowsAffected > 0, nil
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/register.html")
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

	// Обработка POST-запроса (обработка формы)
	userName := r.FormValue("user_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Добавление данных в базу данных
	success, err := addUserToDB(userName, email, password)
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
}

func isValidUser(username, password string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1 AND password_hash= $2", username, password).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/login_page.html")
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

	username := r.FormValue("user_name")
	password := r.FormValue("password")

	// Проверка данных пользователя
	valid, err := isValidUser(username, password)
	if err != nil {
		http.Error(w, "Failed to check credentials", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if valid {
		if username == "admin@mail.ru" && password == "admin" {
			http.Redirect(w, r, "http://localhost:8080/admin/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "http://localhost:8080/dashboard/", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Println("Failed authentication attempt")
	}
}

func dashboardPage(w http.ResponseWriter, r *http.Request) {
	var rows *sql.Rows
	var err error

	// Флаг, указывающий, нужно ли выполнить поиск или нет
	doSearch := r.Method == http.MethodGet && r.URL.Query().Get("departure") != "" && r.URL.Query().Get("arrival") != ""

	// Выполнить поиск, если флаг установлен
	if doSearch {
		departure := r.FormValue("departure")
		arrival := r.FormValue("arrival")

		rows, err = db.Query("SELECT * FROM AvailableFlights WHERE departure_airport = $1 AND arrival_airport = $2",
			departure, arrival)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()
	} else {
		rows, err = db.Query("SELECT * FROM AvailableFlights")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()
	}

	var availableFlights []AvailableFlight

	for rows.Next() {
		var flight AvailableFlight
		if doSearch {
			err := rows.Scan(
				&flight.FlightID,
				&flight.Airline,
				&flight.DepartureTime,
				&flight.ArrivalTime,
				&flight.DepartureAirport,
				&flight.ArrivalAirport,
				&flight.TicketPrice,
			)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Println(err)
				return
			}
		} else {
			err := rows.Scan(
				&flight.FlightID,
				&flight.Airline,
				&flight.DepartureTime,
				&flight.ArrivalTime,
				&flight.DepartureAirport,
				&flight.ArrivalAirport,
				&flight.TicketPrice,
				// Дополнительные поля, которые отображаются только при отсутствии поиска
				// &flight.SomeOtherField,
			)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
		availableFlights = append(availableFlights, flight)
	}

	tmpl, err := template.ParseFiles("/home/sergey/site2/templates/dash_board.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, availableFlights)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}

func buyPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/buy_page.html")
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

	// Обработка данных формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Получение данных из формы
	flightID := r.FormValue("flightID")
	passengerName := r.FormValue("passengerName")
	paymentMethod := r.FormValue("paymentMethod")
	username := r.FormValue("username") // Получаем ID пользователя из сессии или откуда-либо еще

	userIDInt, err := getUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Failed to fetch user ID", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	userID := strconv.Itoa(userIDInt)
	// Пример вызова функции добавления бронировки в базу данных
	success, err := addBookingToDB(userID, flightID, passengerName, paymentMethod)
	if err != nil {
		http.Error(w, "Failed to add booking", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if success {
		// Успешно добавлено
		fmt.Fprintf(w, "Booking added successfully!")
		return
	}

	http.Error(w, "Failed to add booking", http.StatusInternalServerError)
}

func getUserIDByUsername(username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE user_name = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func addBookingToDB(userID, flightID, passengerName, paymentMethod string) (bool, error) {
	// Здесь реализуйте логику добавления бронировки в базу данных, используя переданные параметры

	// Пример вызова хранимой процедуры или SQL-запроса для добавления бронировки
	result, err := db.Exec("select AddBooking($1, $2, $3, $4)", userID, flightID, passengerName, paymentMethod)
	if err != nil {
		return false, err
	}

	// Проверяем количество строк, затронутых операцией вставки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// Возвращаем true, если была вставлена хотя бы одна строка (успешное добавление бронировки)
	return rowsAffected > 0, nil
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Обработка GET-запроса (отображение формы)
		tmpl, err := template.ParseFiles("/home/sergey/site2/templates/admin_board.html")
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

		success, err := addFlight(flight)
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

		success, err := addAirport(airport)
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
	case "deleteAirport":
		airportCode := r.FormValue("airport_code")
		success, err := deleteAirport(airportCode)
		if err != nil {
			http.Error(w, "Failed to delete airport", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			fmt.Fprintf(w, "Airport deleted successfully!")
			return
		} else {
			log.Println("Failed to delete airport from DB")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)
	case "deleteFlight":
		flightIDStr := r.FormValue("flight_id")
		flightID, err := strconv.Atoi(flightIDStr)
		success, err := deleteFlight(flightID)
		if err != nil {
			http.Error(w, "Failed to delete flight", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if success {
			fmt.Fprintf(w, "Flight deleted successfully!")
			return
		} else {
			log.Println("Failed to delete flight from DB")
		}
		http.Redirect(w, r, "http://localhost:8080/", http.StatusSeeOther)

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}

func addFlight(flight Flights) (bool, error) {
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

func addAirport(airport Airports) (bool, error) {
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
