# Домашнее задание к занятию "6.2. SQL"

## Модуль 6. Администрирование баз данных

### Студент: Иван Жиляев

## Задача 1

>Используя docker поднимите инстанс PostgreSQL (версию 12) c 2 volume, 
>в который будут складываться данные БД и бэкапы.
>
>Приведите получившуюся команду или docker-compose манифест.

Для создания инфраструктуры использую [docker-compose файл](task1/docker-compose.yml). Осталось запустить контейнер командой:

```
docker-compose up -d
```

## Задача 2

>В БД из задачи 1: 
>- создайте пользователя test-admin-user и БД test_db
>- в БД test_db создайте таблицу orders и clients (спeцификация таблиц ниже)
>- предоставьте привилегии на все операции пользователю test-admin-user на таблицы БД test_db
>- создайте пользователя test-simple-user  
>- предоставьте пользователю test-simple-user права на SELECT/INSERT/UPDATE/DELETE данных таблиц БД test_db
>
>Таблица orders:
>- id (serial primary key)
>- наименование (string)
>- цена (integer)
>
>Таблица clients:
>- id (serial primary key)
>- фамилия (string)
>- страна проживания (string, index)
>- заказ (foreign key orders)
>
>Приведите:
>- итоговый список БД после выполнения пунктов выше,
>- описание таблиц (describe)
>- SQL-запрос для выдачи списка пользователей с правами над таблицами test_db
>- список пользователей с правами над таблицами test_db

### Выполнение задачи

Подключимся к контейнеру для работы в консоли:

```
docker exec -it task1_sql-server_1 psql -d training -U root
```

Создать пользователя и базу можно командами:

```
CREATE ROLE "test-admin-user" WITH LOGIN;
CREATE DATABASE "test_db";
```

Выберем базу test_db для выполнения дальнейших команд через `\c test_db;` и создадим нужные таблицы:

```
DROP TABLE IF EXISTS orders;
CREATE TABLE orders (
    id serial PRIMARY KEY,
    наименование varchar(255),
    цена int
    );

DROP TABLE IF EXISTS clients;
CREATE TABLE clients (
    id serial PRIMARY KEY,
    фамилия varchar(255),
    "страна проживания" varchar(255),
    заказ int REFERENCES orders
);

DROP INDEX IF EXISTS "страна проживания index";
CREATE INDEX "страна проживания index" ON clients ("страна проживания");
```

Предоставим привилегии пользователю "test-admin-user":

```
GRANT ALL PRIVILEGES ON clients TO "test-admin-user";
GRANT ALL PRIVILEGES ON orders TO "test-admin-user";
```

Создадим пользователя "test-simple-user" и назначим ему права:

```
CREATE ROLE "test-simple-user" WITH LOGIN;
GRANT SELECT, INSERT, UPDATE, DELETE ON clients TO "test-simple-user";
GRANT SELECT, INSERT, UPDATE, DELETE ON orders TO "test-simple-user";
```

### Итог задачи

Приведу команды для получения требуемых данных и их вывод:

>- итоговый список БД после выполнения пунктов выше

```
test_db=# \l
                             List of databases
   Name    | Owner | Encoding |  Collate   |   Ctype    | Access privileges 
-----------+-------+----------+------------+------------+-------------------
 postgres  | root  | UTF8     | en_US.utf8 | en_US.utf8 | 
 template0 | root  | UTF8     | en_US.utf8 | en_US.utf8 | =c/root          +
           |       |          |            |            | root=CTc/root
 template1 | root  | UTF8     | en_US.utf8 | en_US.utf8 | =c/root          +
           |       |          |            |            | root=CTc/root
 test_db   | root  | UTF8     | en_US.utf8 | en_US.utf8 | 
(4 rows)
```

>- описание таблиц (describe)

```
test_db=# \d+ orders;
                                                           Table "public.orders"
    Column    |          Type          | Collation | Nullable |              Default               | Storage  | Stats target | Description 
--------------+------------------------+-----------+----------+------------------------------------+----------+--------------+-------------
 id           | integer                |           | not null | nextval('orders_id_seq'::regclass) | plain    |              | 
 наименование | character varying(255) |           |          |                                    | extended |              | 
 цена         | integer                |           |          |                                    | plain    |              | 
Indexes:
    "orders_pkey" PRIMARY KEY, btree (id)
Referenced by:
    TABLE "clients" CONSTRAINT "clients_заказ_fkey" FOREIGN KEY ("заказ") REFERENCES orders(id)
Access method: heap

test_db=# \d+ clients;
                                                             Table "public.clients"
      Column       |          Type          | Collation | Nullable |               Default               | Storage  | Stats target | Description 
-------------------+------------------------+-----------+----------+-------------------------------------+----------+--------------+-------------
 id                | integer                |           | not null | nextval('clients_id_seq'::regclass) | plain    |              | 
 фамилия           | character varying(255) |           |          |                                     | extended |              | 
 страна проживания | character varying(255) |           |          |                                     | extended |              | 
 заказ             | integer                |           |          |                                     | plain    |              | 
Indexes:
    "clients_pkey" PRIMARY KEY, btree (id)
    "страна проживания index" btree ("страна проживания")
Foreign-key constraints:
    "clients_заказ_fkey" FOREIGN KEY ("заказ") REFERENCES orders(id)
Access method: heap
```

>- SQL-запрос для выдачи списка пользователей с правами над таблицами test_db

```
test_db=# select relacl from pg_catalog.pg_class where relname='clients';
                                         relacl                                          
-----------------------------------------------------------------------------------------
 {root=arwdDxt/root,"\"test-admin-user\"=arwdDxt/root","\"test-simple-user\"=arwd/root"}
(1 row)

test_db=# select relacl from pg_catalog.pg_class where relname='orders';
                                         relacl                                          
-----------------------------------------------------------------------------------------
 {root=arwdDxt/root,"\"test-admin-user\"=arwdDxt/root","\"test-simple-user\"=arwd/root"}
(1 row)
```

>- список пользователей с правами над таблицами test_db

```
test_db=# \dp orders 
                                    Access privileges
 Schema |  Name  | Type  |       Access privileges        | Column privileges | Policies 
--------+--------+-------+--------------------------------+-------------------+----------
 public | orders | table | root=arwdDxt/root             +|                   | 
        |        |       | "test-admin-user"=arwdDxt/root+|                   | 
        |        |       | "test-simple-user"=arwd/root   |                   | 
(1 row)

test_db=# \dp clients
                                    Access privileges
 Schema |  Name   | Type  |       Access privileges        | Column privileges | Policies 
--------+---------+-------+--------------------------------+-------------------+----------
 public | clients | table | root=arwdDxt/root             +|                   | 
        |         |       | "test-admin-user"=arwdDxt/root+|                   | 
        |         |       | "test-simple-user"=arwd/root   |                   | 
(1 row)
```

## Задача 3

>Используя SQL синтаксис - наполните таблицы следующими тестовыми данными:
>
>Таблица orders
>
>|Наименование|цена|
>|------------|----|
>|Шоколад| 10 |
>|Принтер| 3000 |
>|Книга| 500 |
>|Монитор| 7000|
>|Гитара| 4000|
>
>Таблица clients
>
>|ФИО|Страна проживания|
>|------------|----|
>|Иванов Иван Иванович| USA |
>|Петров Петр Петрович| Canada |
>|Иоганн Себастьян Бах| Japan |
>|Ронни Джеймс Дио| Russia|
>|Ritchie Blackmore| Russia|
>
>Используя SQL синтаксис:
>- вычислите количество записей для каждой таблицы 
>- приведите в ответе:
>    - запросы 
>    - результаты их выполнения.

Наполним таблицы:

```
INSERT INTO orders (наименование, цена) VALUES
    ('Шоколад', '10'),
    ('Принтер', '3000'),
    ('Книга', '500'),
    ('Монитор', '7000'),
    ('Гитара', '4000');

INSERT INTO clients (фамилия, "страна проживания") VALUES
    ('Иванов Иван Иванович', 'USA'),
    ('Петров Петр Петрович', 'Canada'),
    ('Иоганн Себастьян Бах', 'Japan'),
    ('Ронни Джеймс Дио', 'Russia'),
    ('Ritchie Blackmore', 'Russia');
```

Количество записей в таблицах можно посчитать запросами:

```
test_db=# select COUNT(*) from clients;
 count 
-------
     5
(1 row)

test_db=# select COUNT(*) from orders;
 count 
-------
     5
(1 row)
```

## Задача 4

>Часть пользователей из таблицы clients решили оформить заказы из таблицы orders.
>
>Используя foreign keys свяжите записи из таблиц, согласно таблице:
>
>|ФИО|Заказ|
>|------------|----|
>|Иванов Иван Иванович| Книга |
>|Петров Петр Петрович| Монитор |
>|Иоганн Себастьян Бах| Гитара |
>
>Приведите SQL-запросы для выполнения данных операций.
>
>Приведите SQL-запрос для выдачи всех пользователей, которые совершили заказ, а также вывод данного запроса.
> 
>Подсказк - используйте директиву `UPDATE`.

Свяжем записи имеющихся таблиц запросами:

```
UPDATE clients SET заказ = 3 WHERE фамилия = 'Иванов Иван Иванович';
UPDATE clients SET заказ = 4 WHERE фамилия = 'Петров Петр Петрович';
UPDATE clients SET заказ = 5 WHERE фамилия = 'Иоганн Себастьян Бах';
```

Выведем пользователей, сделавших заказ:

```
test_db=# select * from clients where заказ is not NULL;
 id |       фамилия        | страна проживания | заказ 
----+----------------------+-------------------+-------
  1 | Иванов Иван Иванович | USA               |     3
  2 | Петров Петр Петрович | Canada            |     4
  3 | Иоганн Себастьян Бах | Japan             |     5
(3 rows)
```

## Задача 5

>Получите полную информацию по выполнению запроса выдачи всех пользователей из задачи 4 
>(используя директиву EXPLAIN).
>
>Приведите получившийся результат и объясните что значат полученные значения.

Результат выполнения:

```
test_db=# EXPLAIN select * from clients where заказ is not NULL;
                         QUERY PLAN                         
------------------------------------------------------------
 Seq Scan on clients  (cost=0.00..10.70 rows=70 width=1040)
   Filter: ("заказ" IS NOT NULL)
(2 rows)
```

Разберём значения:

- `cost=0.00..10.70` - оценка стоимости выполнения запроса; первое число - затраты на получение первой записи, а второе - затраты на обработку всего запроса; стоимость составляется на основе 5 факторов оценки, по её значению в большинстве случаев можно делать вывод о быстродействии запроса.

- `rows=70` - ожидаемое число строк к выводу; откуда взялось такое большое число на нашей маленькой таблице не ясно, но с ключом analyze фактические данные соответствуют размеру трблицы.

- `width=1040` - ожидаемое среднее значение длины каждой строки вывода запроса, значение в байтах.

- `(2 rows)` - количество записей, попавших под фильтр WHERE.

## Задача 6

>Создайте бэкап БД test_db и поместите его в volume, предназначенный для бэкапов (см. Задачу 1).
>
>Остановите контейнер с PostgreSQL (но не удаляйте volumes).
>
>Поднимите новый пустой контейнер с PostgreSQL.
>
>Восстановите БД test_db в новом контейнере.
>
>Приведите список операций, который вы применяли для бэкапа данных и восстановления.

Бэкап сделаем командой:

```
docker exec -it task1_sql-server_1 bash -c 'pg_dump -U root -W test_db > /data/backups/test_db.dump'
```

Запуск пустого PostgreSQL:

```
docker run --rm --name EmptySQL -e POSTGRES_PASSWORD=root -v $(pwd)/backups_test/:/dumps/ -d postgres:12
```

Для восстановления БД требуется сперва создать пустую БД куда необходимо восстановить данные:

```
docker exec -it EmptySQL bash -c 'psql -U postgres -c "CREATE DATABASE test_db"'
```

Теперь можно восстанавливать данные:

```
docker exec -it EmptySQL bash -c 'psql -U postgres test_db < /dumps/test_db.dump'
```
