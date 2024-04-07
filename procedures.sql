--1ДОБАВЛЕНИЕ НОВОГО РЕЙСА добавить если нет рейса и сразу добавить аэропорт
create procedure InsertFlight (
    p_flight_name VARCHAR(50),
    p_airline VARCHAR(100),
    p_aircraft_model VARCHAR(100),
    p_airport_departure VARCHAR(100),
    p_airport_arrival VARCHAR(100),
    p_departure_time TIMESTAMPTZ,
    p_arrival_time TIMESTAMPTZ,
    p_destination VARCHAR(100)
)
language plpgsql as $$
declare 
	departure_exists int;  /*проверка наличия аэпорта отправления*/
	arrival_exists int;  /*проверка наличия аэропортв прибытия*/
begin

    /* проверка наличия аэропорта вылета в базе данных */
    select count(*) into departure_exists
    from airports
    where airport_name = p_airport_departure;

    /* Если аэропорт вылета отсутствует, добавляем его */
    if departure_exists = 0 then
        insert into airports (airport_name)
        values (p_airport_departure);
    end if;

    /* проверка наличия аэропорта прибытия в базе данных */
    select count(*) into arrival_exists
    from airports
    where airport_name = p_airport_arrival;

    /* Если аэропорт прибытия отсутствует, добавляем его */
    if departure_exists = 0 then
        insert into airports (airport_name)
        values (p_airport_arrival);
    end if;

    /* Добавление нового рейса */
    insert into Flights (flight_name, airline, aircraft_model, airport_departure, airport_arrival, departure_time, arrival_time, destination)
    values (p_flight_name, p_airline,p_aircraft_model, p_airport_departure, p_airport_arrival, p_departure_time, p_arrival_time, p_destination);

end;
$$;
/*вызываем через CALL AddFlight('Flight123', 'Airline XYZ', 'DepartureAirport', 'ArrivalAirport', 
'2024-01-04 08:00:00', '2024-01-04 12:00:00', 'Destination');  */



--2ДОБАВЛЕНИЕ НОВОГО ПОЛЬЗОВАТЕЛЯ

create or replace procedure InsertUser(
    p_username varchar(50),
    p_user_lastname varchar(50),
    p_email varchar(100),
    p_password_hash varchar(100)
)

language plpgsql as $$
begin
    insert into Users(username, user_lastname, email, password_hash)
    values (p_username, p_email, p_user_lastname, p_password_hash);
end;
$$;


--3ДОБАВЛЕНИЕ БИЛЕТА

create or replace procedure InsertTickets(
    p_flight_name INT ,  -- На какой рейс
    p_ticket_price NUMERIC(8, 2),
    p_seat_number VARCHAR(10)
) 

language plpgsql as $$
begin
    insert into Tickets(flight_name, ticket_price, seat_namber)
    values (p_flight_name, p_ticket_price, p_seat_namber);
end;
$$;


--4ДОБАВЛЕНИЕ ВИП-БИЛЕТА

create or replace procedure InsertTickets(
    p_flight_name INT ,  -- На какой рейс
    p_ticket_price NUMERIC(8, 2),
    p_seat_number VARCHAR(10),
	p_meal_choice VARCHAR(50),
	p_drink_alcohol_choice VARCHAR(50),
	p_personal_concierge BOOLEAN
) 

language plpgsql as $$
begin
    insert into Tickets(flight_name, ticket_price, seat_namber, meal_choice, drink_alcohol_choice, personal_concierge)
    values (p_flight_name, p_ticket_price, p_seat_namber, p_meal_choice, p_drink_alcohol_choice, p_personal_concierge);
end;
$$;


--5ДОБАВЛЕНИЕ БРОНИРОВАНИЯ

create or replace procedure InsertBooking(
    p_user_id int,
    p_flight_name VARCHAR(50),
    p_seat_number varchar(10),
    p_vip_seat_number varchar(10)
)
language plpgsql as $$
begin
    if (p_seat_number is not null and p_vip_seat_number is not null) or
       (p_seat_number is null and p_vip_seat_number is null) then
        return; -- операция не выполнена
    end if;

    insert into Bookings (user_id, flight_name, seat_number, vip_seat_number)
    values (p_user_id, p_flight_name, p_seat_number, p_vip_seat_number)
end;
$$;


--6ДОБАВЛЕНИЕ ОПЛАТЫ

create or replace procedure InsertPayment( 
    p_amount numeric(8,2),
    p_payment_method varchar(50)
)
language plpgsql as $$
begin
    insert into Payments (booking_id, amount, payment_method)
    values (p_booking_id, p_amount, p_payment_method);
end;
$$;




--7ДОБАВЛЕНИЕ АЭРОПОРТА
create or replace procedure InsertAirport(
    p_airport_name VARCHAR(100),
    p_airport_city VARCHAR(100)
)
language plpgsql as $$
begin
    insert into Airports (airport_name, airport_city)
    values (p_airport_name, coalesce(p_airport_city, 'Moscow'));
end;
$$;

--8ДОБАВЛЕНИЕ САМОЛЕТА

create or replace procedure AddAircraft(
    p_aircraft_model VARCHAR(100),
    p__total_seats INT,
    p_aircraft_condition VARCHAR(50)
)
language plpgsql as $$
begin 
    insert into Aircrafts (aircraft_model, _seats, aircraft_condition)
    values (p_aircraft_model, p_seats, coalesce(p_aircraft_condition, 'Flies'));
end;
$$;


--1ОТМЕНА РЕЙСА

create or replace procedure DeleteFlight(
    p_flight_name varchar(50)
)
language plpgsql as $$
begin
    -- Удаление рейса
    delete from Flights where flight_name = p_flight_name;
end;
$$;

--2УДАЛЕНИЕ АЭРОПОРТА

create or replace procedure DeleteAirport(
    p_airport_name VARCHAR(100)
) 
language plpgsql as $$
begin
    -- Удаление из таблицы Airports
    delete from Airports where airport_name = p_airport_name;
end;
$$;

--3УДАЛЕНИЕ БРОНИ НА БИЛЕТ

create or replace procedure DeleteBooking(
    p_booking_id INT
)
language plpgsql as $$
begin
    delete from Bookings where booking_id = p_booking_id;
end;
$$;

--1ОБНОВЛЕНИЕ КОЛИЧЕСТВА МЕСТ В САМОЛЕТЕ

create or replace procedure UpdateAvailableSeats(
    p_flight_name VARCHAR(50)
)
language plpgsql as $$
declare
    v_total_seats INT;
    v_booked_seats INT;
begin
    -- Получаем общее количество мест в самолете
    select total_seats into v_total_seats
    from Aircrafts a
    join Flights f on a.aircraft_id = f.aircraft_id
    where f.flight_name = p_flight_name;

    -- Получаем количество забронированных мест
    select COUNT(*) into v_booked_seats
    from Bookings
    where flight_name = p_flight_name;

    -- Обновляем количество доступных мест в самолете
    update Aircrafts
    set total_seats = v_total_seats - v_booked_seats
    from Flights
    where Aircrafts.aircraft_id = Flights.aircraft_id
    and Flights.flight_name = p_flight_name;
end;
$$;
