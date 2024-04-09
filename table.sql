-- Создание родительской таблицы без разделения
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    user_lastname VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(50) DEFAULT 'basic_user'
);

-- Создание родительской таблицы для разделения
-- Создание таблицы Users_partitioned с разделением по году регистрации
CREATE TABLE Users_partitioned (
    user_id SERIAL,
    username VARCHAR(50) NOT NULL,
    user_lastname VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(50) DEFAULT 'basic_user'
) PARTITION BY RANGE (EXTRACT(YEAR FROM registration_date));


-- Создание дочерних таблиц (разделов) для каждого года
CREATE TABLE Users_2024 PARTITION OF Users_partitioned
FOR VALUES FROM (2024) TO (2025);

CREATE TABLE Users_2024 PARTITION OF Users_partitioned
FOR VALUES FROM (2025) TO (2026);


--Админ
CREATE TABLE Admins (
    admin_id SERIAL PRIMARY KEY, -- Уникальный идентификатор админа
    email VARCHAR(100) UNIQUE NOT NULL, -- Email админа
    password_hash VARCHAR(100) NOT NULL, -- Захэшированный пароль админа
    admin_role VARCHAR(50) DEFAULT 'admin' -- Роль админа (может быть различными, например, "admin", "superadmin")
);

-- Рейсы
create table Flights (
    flight_id SERIAL primary key,     --не показывать пользователю
    flight_name VARCHAR(50) unique not null,   --номер рейса
    airline VARCHAR(100) not null, --имя компании перевозчика
    aircraft_model VARCHAR(100) references Aircrafts(aircraft_model), --модель самолета
    airport_departure  VARCHAR(100) references Airports(airport_name),  --аэропорт отправления
    airport_arrival  VARCHAR(100) references Airports(airport_name),    --аэропорт прибытия
    departure_time TIMESTAMPTZ not null,   --время отправления
    arrival_time TIMESTAMPTZ not null,     --время прибытия   
    destination VARCHAR(100) not null   --дистанция
);


--Билеты
create table Tickets (
	ticket_id SERIAL primary key,  -- Идентификатор билета
    flight_name VARCHAR(50) references Flights(flight_name),  -- На какой рейс
    ticket_price NUMERIC(8, 2) not null,  -- Цена
    seat_number VARCHAR(10) not null -- Номер места
);

--Вип-билеты как наследуемая от билетов
create table VipTickets (
    meal_choice VARCHAR(50) default 'Meat medium rare', --выбор обеда
    drink_alcohol_choice VARCHAR(50) default 'Vodka', --выбор напитков
    personal_concierge BOOLEAN default true--персональный консьерж
) inherits (Tickets);


--Бронирование Билетов
create table Bookings (
    booking_id SERIAL primary key,   -- Не показывать пользователю
    user_id INT references Users(user_id),     -- От кого бронь
	id_ticket int references Tickets(ticket_id),  --место в обычном салоне
    booking_status VARCHAR(20) DEFAULT 'Pending'
    );   -- По умолчанию ожидает подтверждения

-- Аэропорты
create table Airports (
    airport_id SERIAL primary key,
    airport_name VARCHAR(100) unique not null,
    airport_city VARCHAR(100) default 'Moscow'  --город 
);

-- Информация о Самолете
create table Aircrafts (
    aircraft_id SERIAL primary key,   --номер самолета
    aircraft_model VARCHAR(100) unique not null,  --модель
    total_seats INT not null ,  --общее количество мест
    aircraft_condition VARCHAR(50) default 'Flies' --состояние самолета 
);
