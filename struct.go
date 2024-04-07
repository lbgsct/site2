package main

import (
	"time"
)

type QueryParams struct {
	Departure string
	Arrival   string
	Date      string
}

type UserFlights struct {
	FlightName       string
	Airline          string
	TotalSeats		 int64
	AircraftModel	 string
	DepartureCity	 string
	DepartureAirport string
	DepartureTime    time.Time
	ArrivalCity 	 string
	ArrivalAirport   string
	ArrivalTime      time.Time
	MinTicketPrice   float64
}

type FlightSeatPrice struct {
    FlightName  string
	SeatNumber	string
    TicketPrice float64
}

type FlightVIPSeatPrice struct {
    FlightName  string
	SeatNumber	string
    TicketPrice float64
	MealChoice	string
	DrinkAlkohol string
	PersonalConcierge bool
}


type Flights struct {
    FlightName       string    `json:"flight_name"`
    Airline          string    `json:"airline"`
    AircraftModel    string    `json:"aircraft_model"`
    DepartureAirport string    `json:"departure_airport"`
    ArrivalAirport   string    `json:"arrival_airport"`
    DepartureTime    time.Time `json:"departure_time"`
    ArrivalTime      time.Time `json:"arrival_time"`
    Destination      string    `json:"destination"`
}

type Airports struct {
    AirportName string
    AirportCity string
}

type Tickets struct {
	FlightName      string
	Price 			float64
	SeatNumber		string

}

type VipTickets struct {
	FlightName      string
	Price 			float64
	SeatNumber		string
	MealChoice		string
	DrinkAlkohol 	string
	PersonalConcierge	bool
}

type Aircrafts struct {
	AircraftModel       string
	TotalSeats			int32
	AircraftCondition	string

}

type User struct {
	UserID	int
	UserName	string
	UserLastname	string
	Email	string
	PasswordHash	string
}

type BookingDetail struct {
    BookingID     int
    UserID        int
    TicketID      int32
    FlightName    string
    TicketPrice   float64
    SeatNumber    string
    BookingStatus string
}


type UserPartition struct {
	UserID          int
	Username        string
	UserLastname    string
	Email           string
	PasswordHash    string
	RegistrationDate string
	Role            string
}