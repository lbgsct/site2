package main

import (
	"time"
)

type Flights struct {
	flight_id		int
	Airline        string
	departure_time time.Time
	arrival_time   time.Time
	origin         string
	destination    string
	ticket_price   float64
}

type Airports struct {
	airport_id				int32
	airport_code          string
	airport_name          string
	location              string
	other_characteristics string
}

type QueryParams struct {
	Departure string
	Arrival   string
	Date      string
}


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


type Fligh struct {
	ID             int
	Airline        string
	DepartureTime  time.Time
	ArrivalTime    time.Time
	Origin         string
	Destination    string
	TicketPrice    float64
}

type Airport struct {
	ID                    int
	AirportCode           string
	AirportName           string
	Location              string
	OtherCharacteristics  string
}
