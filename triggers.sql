--ТРИГГЕРЫ НА КАСКАДНОЕ УДАЛЕНИЕ


-- Пользователи
alter table Users
drop constraint if exists users_bookings_fk,
add constraint users_bookings_fk foreign key (user_id) references Bookings(user_id) on delete cascade;

-- Рейсы
alter table Flights
drop constraint if exists flights_tickets_fk,
add constraint flights_tickets_fk foreign  key  (flight_name) references  Tickets(flight_name) on  delete  cascade;

-- Бронирование Билетов
alter table Bookings
drop constraint if exists bookings_flights_fk,
drop constraint if exists bookings_tickets_fk,
drop constraint if exists bookings_viptickets_fk,
add constraint bookings_flights_fk foreign key (flight_name) references Flights(flight_name) on delete cascade,
add constraint bookings_tickets_fk foreign key (ticket_price, seat_number) references Tickets(ticket_price, seat_number) on delete cascade,
add constraint bookings_viptickets_fk foreign key (vip_seat_number) references VipTickets(vip_seat_number) on delete cascade;


-- Оплаты
alter table Payments
drop constraint if exists payments_bookings_fk,
add constraint payments_bookings_fk foreign key (booking_id) references Bookings(booking_id) on delete cascade;

-- Аэропорты (в случае, если у вас есть внешние ключи, которые ссылаются на это)
alter  table  Airports
drop  constraint  if  exists  airports_flights_departure_fk,
drop  constraint  if  exists  airports_flights_arrival_fk,
add  constraint  airports_flights_departure_fk foreign  key  (airport_name) references  Flights(airport_departure) on  delete cascade ,
add  constraint  airports_flights_arrival_fk foreign  key  (airport_name) references  Flights(airport_arrival) on delete cascade;

--	ТРИГГЕРЫ  НА ПРОВЕРКУ ДАННЫХ

--правильный имейл

create OR REPLACE FUNCTION validate_user_email()
returns TRIGGER AS $$
begin
    if new.email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}$' THEN
        return new;
    else
        raise exception 'Invalid email format for Users table';
    end if;
end;
$$ language plpgsql;

create trigger users_email_trigger
before insert or update on Users
for each row execute function validate_user_email();

--правильное время отправки и прибытия

create  or  replace  function  validate_flight_times()
returns  trigger  as  $$
begin 
    if  new.departure_time < new.arrival_time then 
        return  new ;
    else 
        raise  exception  'Invalid flight times for Flights table';
    end  if ;
end ;
$$ language  plpgsql;

create trigger flights_time_trigger
before insert or update on Flights
for each row execute function validate_flight_times();

--цена билета

create or replace function validate_ticket_price()
returns trigger as $$
begin
    if new.ticket_price > 0 then
        return new;
    else
        raise exception 'invalid ticket price for tickets table';
    end if;
end;
$$ language plpgsql;

create trigger tickets_price_trigger
before insert or update on tickets
for each row execute function validate_ticket_price();


--цена вип билета

create or replace function validate_vip_ticket_price()
returns trigger as $$
begin
    if new.vip_ticket_price > 0 then
        return new;
    else
        raise exception 'invalid vip ticket price for viptickets table';
    end if;
end;
$$ language plpgsql;

create trigger viptickets_price_trigger
before insert or update on viptickets
for each row execute function validate_vip_ticket_price();

--проверка места

create or replace function validate_booking_seats()
returns trigger as $$
begin
    if (new.seat_number is not null and new.vip_seat_number is null) or (new.seat_number is null and new.vip_seat_number is not null) then
        return new;
    else
        raise exception 'invalid seat assignment for bookings table';
    end if;
end;
$$ language plpgsql;

create trigger bookings_seats_trigger
before insert or update on bookings
for each row execute function validate_booking_seats();

--сумма заказа
create or replace function validate_payment_amount()
returns trigger as $$
begin
    if new.amount > 0 then
        return new;
    else
        raise exception 'invalid payment amount for payments table';
    end if;
end;
$$ language plpgsql;

create trigger payments_amount_trigger
before insert or update on payments
for each row execute function validate_payment_amount();



