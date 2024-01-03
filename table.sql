-- Пользователи
create table Users (
    user_id SERIAL primary key,   --id по которому будет страничка пользователя
    username VARCHAR(50) unique not null,   --просто имя, мб добавлю фамилию 
    email VARCHAR(100) UNIQUE not null,   --вход по имейлу 
    password_hash VARCHAR(100) not null,  --захэшировать пароль
    registration_date TIMESTAMP default CURRENT_TIMESTAMP,  --дата регистрации
    role VARCHAR(50) default 'basic_user'    --роль для прав доступа
);

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
    airline VARCHAR(100) not null,    --имя компании перевозчика
    airport_departure INT references Airports(airport_name),  --аэропорт отправления
    airport_arrival INT references Airports(airport_name),    --аэропорт прибытия
    departure_time TIMESTAMPTZ not null,   --время отправления
    arrival_time TIMESTAMPTZ not null,     --время прибытия   
    destination VARCHAR(100) not null   --дистанция
);

--Билеты
create table Tickets (
	ticket_id SERIAL primary key,  -- Идентификатор билета
    flight_name INT references Flights(flight_name),  -- На какой рейс
    ticket_price NUMERIC(8, 2) not null,  -- Цена
    seat_number VARCHAR(10)  -- Номер места
);


--Вип-билеты как наследуемая от билетов
create table VipTickets (
    vip_ticket_id SERIAL primary key,
    flight_name INT references Flights(flight_name), -- На какой рейс (используем flight_id)
    vip_ticket_price NUMERIC(12, 2) not null, -- Цена
    vip_seat_number VARCHAR(10) not null, -- Номер места
    meal_choice VARCHAR(50) not null, --выбор обеда
    drink_choice VARCHAR(50) not null, --выбор напитков
    personal_concierge BOOLEAN not null --персональный консьерж
) inherits (Tickets);

--Бронирование Билетов
create table Bookings (
    booking_id SERIAL primary key,   -- Не показывать пользователю
    user_id INT references Users(user_id),     -- От кого бронь
    flight_name INT references Flights(flight_name),  -- На какой рейс
	seat_number VARCHAR(10) references Tickets(seat_number),  --место в обычном салоне
	vip_seat_number VARCHAR(10) references VipTickets(vip_seat_number),  --место в бизнес классе
    booking_status VARCHAR(20) DEFAULT 'Pending',   -- По умолчанию ожидает подтверждения
    CONSTRAINT check_ticket_type CHECK (
        (seat_number is not null and vip_seat_number is null) or
        (seat_number is null and vip_seat_number is not null)
    )
);


-- Создание таблицы Оплаты
create table Payments (
    payment_id SERIAL primary key,  
    booking_id INT references Bookings(booking_id),  
    amount NUMERIC(10, 2) not null,  --цена
    payment_method VARCHAR(50) not null,  --способ оплаты
    payment_date TIMESTAMPTZ default CURRENT_TIMESTAMP  --дата оплаты
);

-- Аэропорты
create table Airports (
    airport_id SERIAL primary key,
    airport_name VARCHAR(100) unique not null,
    airport_city VARCHAR(100) not null   --город 
);

-- Информация о Самолете
create table Aircrafts (
    aircraft_id SERIAL primary key,   --номер самолета
    aircraft_model VARCHAR(100) not null,  --модель
    _seats INT not null ,  --общее количество мест
    aircraft_condition VARCHAR(50) default  --состояние самолета 
);
