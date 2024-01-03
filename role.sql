-- Создание ролей

-- Администратор
create role admin with LOGIN password 'admin_password' SUPERUSER;

-- Пользователь
create role basic_user with LOGIN password 'user_password';

-- Назначение прав доступа для ролей

-- Администратору назначаются полные права на все таблицы
grant all privileges on all tables in schema public to admin;

-- Пользователю basic_user назначаются права на чтение и модификацию таблиц Users и Bookings
grant select, insert, update on Users to basic_user;
grant select, insert on Bookings to basic_user;

-- Пользователь basic_user также может читать таблицы Flights, Tickets, Payments и Airports
grant select on Flights to basic_user;
grant select on Tickets to basic_user;
grant select on Payments to basic_user;
grant select on Airports to basic_user;
grant select on Aircrafts to basic_user;

-- Отдельно для VipTickets, так как она наследуется от Tickets
grant select, insert, update on VipTickets to basic_user;




--Секционирование журнала действий пользователей
create table UserActionsLogs (
    log_id SERIAL,
    user_id INT,
    action_time TIMESTAMP NOT NULL,
    action_details TEXT,
    primary key (log_id, action_time)
) partition by range (action_time);

--Создание секций
create table UserActionsLogs_2023_01 partition of UserActionsLogs
for values from ('2023-01-01') to ('2023-02-01');

create table UserActionsLogs_2023_02 partition of UserActionsLogs
for values from ('2023-02-01') to ('2023-03-01');

--Индексация секций
create index idx_user_id on UserActionsLogs (user_id);