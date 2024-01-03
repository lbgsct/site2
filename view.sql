--ИНФА ДЛЯ ПОЛЬЗОВАТЕЛЯ

CREATE VIEW UserFlights AS
SELECT 
    f.flight_name,
    f.airline,
    a.total_seats AS total_seats,
    f.aircraft_model,
    dep.airport_city AS departure_city,
    arr.airport_city AS arrival_city,
    f.departure_time,
    f.arrival_time,
    f.destination
FROM Flights f
JOIN Aircrafts a ON f.aircraft_model = a.aircraft_model
JOIN Airports dep ON f.airport_departure = dep.airport_id
JOIN Airports arr ON f.airport_arrival = arr.airport_id;

   