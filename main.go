package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login/", LoginPage)
	http.HandleFunc("/register/", Register)
	http.HandleFunc("/login/user/", DashboardPage)
	http.HandleFunc("/login/admin/", AdminPage)
	http.HandleFunc("/login/user/buy/", BuyPage)
	http.HandleFunc("/login/admin/add/", AdminAdd)
	http.HandleFunc("/login/admin/update/", AdminUpdate)
	http.HandleFunc("/login/admin/delete/", AdminDelete)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func getUserIDByUsername(username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE user_name = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}





/*func adminSee(w http.ResponseWriter, r *http.Request) {
	// Получение данных о рейсах и аэропортах из базы данных (предположим, используется функция getFlightsFromDB и getAirportsFromDB)

	flights, err := getFlightsFromDB(db) // Получение данных о рейсах из базы данных
	if err != nil {
		http.Error(w, "Failed to get flights", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	airports, err := getAirportsFromDB(db) // Получение данных об аэропортах из базы данных
	if err != nil {
		http.Error(w, "Failed to get airports", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Отображение данных на веб-странице с использованием шаблонов HTML
	tmpl, err := template.ParseFiles("/home/sergey/site2/templates/see.html") // Путь к вашему шаблону HTML
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := struct {
		Flights  []Flights
		Airports []Airports
	}{
		Flights:  flights,
		Airports: airports,
	}

	err = tmpl.Execute(w, data) // Передача данных в шаблон для отображения
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}*/

/*func getFlightsFromDB(db *sql.DB) ([]Flights, error) {
	rows, err := db.Query("SELECT * FROM flights") // Выполнение запроса на получение всех данных из таблицы рейсов
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []Flights
	for rows.Next() {
		var flight Flights
		err := rows.Scan(&flight.ID, &flight.Airline, &flight.DepartureTime, &flight.ArrivalTime, &flight.Origin, &flight.Destination, &flight.TicketPrice)
		if err != nil {
			return nil, err
		}
		flights = append(flights, flight)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return flights, nil
}*/

/*func getAirportsFromDB(db *sql.DB) ([]Airports, error) {
	rows, err := db.Query("SELECT * FROM airports") // Выполнение запроса на получение всех данных из таблицы аэропортов
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports []Airports
	for rows.Next() {
		var airport Airports
		err := rows.Scan(&airport.ID, &airport.AirportCode, &airport.AirportName, &airport.Location, &airport.OtherCharacteristics)
		if err != nil {
			return nil, err
		}
		airports = append(airports, airport)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return airports, nil
}*/
