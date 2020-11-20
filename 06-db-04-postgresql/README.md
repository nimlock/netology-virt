# Домашнее задание к занятию "6.4. PostgreSQL"

## Модуль 6. Администрирование баз данных

### Студент: Иван Жиляев

## Задача 1

>Используя docker поднимите инстанс PostgreSQL (версию 13). Данные БД сохраните в volume.
>
>Подключитесь к БД PostgreSQL используя `psql`.
>
>Воспользуйтесь командой `\?` для вывода подсказки по имеющимся в `psql` управляющим командам.
>
>**Найдите и приведите** управляющие команды для:
>- вывода списка БД
>- подключения к БД
>- вывода списка таблиц
>- вывода описания содержимого таблиц
>- выхода из psql

Для запуска контейнера используем манифест [docker-compose](docker-compose.yml).

```
docker-compose up -d
```

Подключимся в контейнер:

```
docker exec -it 06-db-04-postgresql_sql-server_1 psql
```

Приведу строки, описывающие необходимые управляющие команды для:

>- вывода списка БД

```
\l[+]   [PATTERN]      list databases
```

>- подключения к БД

```
  \c[onnect] {[DBNAME|- USER|- HOST|- PORT|-] | conninfo}
                         connect to new database (currently "root")
```

>- вывода списка таблиц

```
\d[S+]                 list tables, views, and sequences
```

>- вывода описания содержимого таблиц

```
\d[S+]  NAME           describe table, view, sequence, or index
```

>- выхода из psql

```
\q                     quit psql
```

## Задача 2

>Используя `psql` создайте БД `test_database`.
>
>Изучите [бэкап БД](https://github.com/netology-code/virt-homeworks/tree/master/06-db-04-postgresql/test_data).
>
>Восстановите бэкап БД в `test_database`.
>
>Перейдите в управляющую консоль `psql` внутри контейнера.
>
>Подключитесь к восстановленной БД и проведите операцию ANALYZE для сбора статистики по таблице.
>
>Используя таблицу [pg_stats](https://postgrespro.ru/docs/postgresql/12/view-pg-stats) столбец таблицы `orders` 
>с наибольшим средним значением размера элементов в байтах.
>
>**Приведите в ответе** команду, которую вы использовали для вычисления и полученный результат.


Для восстановления БД требуется сперва создать пустую БД куда необходимо восстановить данные:

```
docker exec -it 06-db-04-postgresql_sql-server_1 bash -c 'psql -U "$POSTGRES_USER" -c "CREATE DATABASE test_database"'
```

Теперь можно восстанавливать данные:

```
docker exec -i 06-db-04-postgresql_sql-server_1 bash -c 'psql -U "$POSTGRES_USER" test_database' < ./test_data/test_dump.sql
```

Зайдём в контейнер и выполним операцию ANALYZE над восстановленной базой с помощью команд:

```
docker exec -it 06-db-04-postgresql_sql-server_1 bash -c 'psql -U "$POSTGRES_USER"'
postgres=# \c test_database
You are now connected to database "test_database" as user "postgres".
test_database=# ANALYZE;
ANALYZE
```

Для получения максимального значения "среднего размера элемента" выполним запрос:

```
test_database=# select max(avg_width) from pg_stats where tablename = 'orders';
-[ RECORD 1 ]
max | 16
```

Чтобы вывести соответствующее имя столбца импользуем вложенный запрос:

```
test_database=# select attname as "Имя столбца", avg_width as "Наибольшее значение" from pg_stats where tablename = 'orders' and avg_width = (select max(avg_width) from pg_stats where tablename = 'orders');
 Имя столбца | Наибольшее значение 
-------------+---------------------
 title       |                  16
(1 row)
```


## Задача 3

>Архитектор и администратор БД выяснили, что ваша таблица orders разрослась до невиданных размеров и
>поиск по ней занимает долгое время. Вам, как успешному выпускнику курсов DevOps в нетологии предложили
>провести разбиение таблицы на 2 (шардировать на orders_1 - price>499 и orders_2 - price<=499).
>
>Предложите SQL-транзакцию для проведения данной операции.
>
>Можно ли было изначально исключить "ручное" разбиение при проектировании таблицы orders?

Нам нужно в рамках транзакции создать "рядом" с исходной таблицей такую же табицу по структуре, но с дополнением в виде шардирования. Затем требуется перенести данные из оригинальной таблицы в новую, при этом записи автоматически разделятся по шардам. В финале остаётся только переименовать таблицы.  
Напишем такую транзакцию:

```
BEGIN;

CREATE TABLE orders_with_shards (
    id integer NOT NULL,
    title character varying(80) NOT NULL,
    price integer DEFAULT 0
)
PARTITION BY RANGE (price);

CREATE TABLE orders_2 PARTITION OF orders_with_shards FOR VALUES FROM (MINVALUE) TO (500);

CREATE TABLE orders_1 PARTITION OF orders_with_shards FOR VALUES FROM (500) TO (MAXVALUE);

INSERT INTO orders_with_shards SELECT * FROM orders;

ALTER TABLE orders RENAME TO orders_old;

ALTER TABLE orders_with_shards RENAME TO orders;

COMMIT;
```

Изначально исключить "ручное" разбиение можно с помощью trigger-ов на INSERT в таблицу, чтобы при определённых условиях создавался новый шард. Однако довольно непросто заранее предугадать логику шардирования, плюс это может быть затратно по ресурсам так как триггер будет проверять каждую запись в таблицу.  
Вероятно компромиссом могло бы стать изначальное определение таблицы как разделяемой, но только с одной партицией с правилом размещения данных DEFAULT: в этом случае не потребовалось бы пересоздавать таблицу (она сразу будет секционированной), а можно было бы просто добавить шарды и переместить данные по секциям. Это тоже "ручной" вариант, но он немного проще.

## Задача 4

>Используя утилиту `pg_dump` создайте бекап БД `test_database`.
>
>Как бы вы доработали бэкап-файл, чтобы добавить уникальность значения столбца `title` для таблиц `test_database`?

[Бэкап](test_data/backup.dump) сделаем командой:

```
docker exec -i 06-db-04-postgresql_sql-server_1 bash -c 'pg_dump -U "$POSTGRES_USER" test_database' > ./test_data/backup.dump
```

Добавить уникальности значениям в столбце `title` можно, например, дополнив название годом издания и автором книги. Для этого достаточно изменить строки дампа, в которых хранятся пользовательские данные из таблиц. 
