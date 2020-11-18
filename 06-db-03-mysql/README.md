# Домашнее задание к занятию "6.3. MySQL"

## Модуль 6. Администрирование баз данных

### Студент: Иван Жиляев

## Задача 1

>Используя docker поднимите инстанс MySQL (версию 8). Данные БД сохраните в volume.
>
>Изучите [бэкап БД](https://github.com/netology-code/virt-homeworks/tree/master/06-db-03-mysql/test_data) и 
>восстановитесь из него.
>
>Перейдите в управляющую консоль `mysql` внутри контейнера.
>
>Используя команду `\h` получите список управляющих команд.
>
>Найдите команду для выдачи статуса БД и **приведите в ответе** из ее вывода версию сервера БД.
>
>Подключитесь к восстановленной БД и получите список таблиц из этой БД.
>
>**Приведите в ответе** количество записей с `price` > 300.
>
>В следующих заданиях мы будем продолжать работу с данным контейнером.

Для запуска контейнера используем манифест [docker-compose](docker-compose.yml).

```
docker-compose up -d
```

Восстановим БД командой:

```
docker exec -i 06-db-03-mysql_sql-server_1 sh -c 'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD" -b "$MYSQL_DATABASE"' < ./test_data/test_dump.sql
```

Подключимся к контейнеру интерактивно:

```
docker exec -it 06-db-03-mysql_sql-server_1 mysql -uroot -p
```

Статус MySQL-сервера можно получить командой `\s` или `status`. Вывод сообщает, что версия сервера __Server version: 8.0.22 MySQL Community Server - GPL__.

Выберем базу для работы и получим список таблиц из этой БД (вывод сообщает что есть только таблица `orders`):

```
use test_db;
show tables;
```

Определим количество записей с `price` > 300:

```
mysql> select count(*) from orders where price > 300;
+----------+
| count(*) |
+----------+
|        1 |
+----------+
1 row in set (0.21 sec)
```

## Задача 2

>Создайте пользователя test в БД c паролем test-pass, используя:
>- плагин авторизации mysql_native_password
>- срок истечения пароля - 180 дней 
>- количество попыток авторизации - 3 
>- максимальное количество запросов в час - 100
>- аттрибуты пользователя:
>    - Фамилия "Pretty"
>    - Имя "James"
>
>Предоставьте привелегии пользователю `test` на операции SELECT базы `test_db`.
>    
>Используя таблицу INFORMATION_SCHEMA.USER_ATTRIBUTES получите данные по пользователю `test` и 
>**приведите в ответе к задаче**.

Создадим пользователя:

```
CREATE USER IF NOT EXISTS
  'test'@'localhost' IDENTIFIED WITH mysql_native_password BY 'test-pass'
  REQUIRE NONE
  WITH MAX_QUERIES_PER_HOUR 100
  FAILED_LOGIN_ATTEMPTS 3
  PASSWORD EXPIRE INTERVAL 180 DAY
  ATTRIBUTE '{"fname": "James", "lname": "Pretty"}';
```

Предоставим привелегии пользователю `test` на операции SELECT на все таблицы базы `test_db`:

```
GRANT SELECT ON test_db.* TO 'test'@'localhost';
```

Получить данные по пользователю `test` можно командой:

```
mysql> SELECT * FROM INFORMATION_SCHEMA.USER_ATTRIBUTES WHERE USER = 'test';
+------+-----------+---------------------------------------+
| USER | HOST      | ATTRIBUTE                             |
+------+-----------+---------------------------------------+
| test | localhost | {"fname": "James", "lname": "Pretty"} |
+------+-----------+---------------------------------------+
1 row in set (0.00 sec)
```

## Задача 3

>Установите профилирование `SET profiling = 1`.
>Изучите вывод профилирования команд `SHOW PROFILES;`.
>
>Исследуйте, какой `engine` используется в таблице БД `test_db` и **приведите в ответе**.
>
>Измените `engine` и **приведите время выполнения и запрос на изменения из профайлера в ответе**:
>- на `MyISAM`
>- на `InnoDB`


Для определения используемого движка выполним команду:

```
mysql> SELECT ENGINE FROM information_schema.tables WHERE table_schema = 'test_db'\G;
*************************** 1. row ***************************
ENGINE: InnoDB
1 row in set (0.00 sec)
```

Сменим движок:

```
mysql> ALTER TABLE orders engine=MyISAM;
Query OK, 5 rows affected (0.14 sec)
Records: 5  Duplicates: 0  Warnings: 0

mysql> ALTER TABLE orders engine=InnoDB;
Query OK, 5 rows affected (0.30 sec)
Records: 5  Duplicates: 0  Warnings: 0
```

Посмотрим записи в профайлере:

```
mysql> SHOW PROFILES\G;

<вывод сокращён для удобства>

*************************** 14. row ***************************
Query_ID: 36
Duration: 0.13994525
   Query: ALTER TABLE orders engine=MyISAM
*************************** 15. row ***************************
Query_ID: 37
Duration: 0.29994150
   Query: ALTER TABLE orders engine=InnoDB
```

## Задача 4 

>Изучите файл `my.cnf` в директории /etc/mysql.
>
>Измените его согласно ТЗ (движок InnoDB):
>- Скорость IO важнее сохранности данных
>- Нужна компрессия таблиц для экономии места на диске
>- Размер буффера с незакомиченными транзакциями 1 Мб
>- Буффер кеширования 30% от ОЗУ
>- Размер файла логов операций 100 Мб
>
>Приведите в ответе измененный файл `my.cnf`.

Для соответствия требованиям ТЗ конфигурационный файл был изменён до следующего состояния:

```
root@af5d3de0189b:/etc/mysql# cat my.cnf
# Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; version 2 of the License.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301 USA

#
# The MySQL  Server configuration file.
#
# For explanations see
# http://dev.mysql.com/doc/mysql/en/server-system-variables.html

[mysqld]
pid-file        = /var/run/mysqld/mysqld.pid
socket          = /var/run/mysqld/mysqld.sock
datadir         = /var/lib/mysql
secure-file-priv= NULL

innodb_buffer_pool_size = 300M  # Need to be calculate manual
innodb_log_file_size = 100M
innodb_log_buffer_size = 1M
innodb_file_per_table = ON
innodb_flush_method = O_DSYNC
innodb_flush_log_at_trx_commit = 2

# Custom config should go here
!includedir /etc/mysql/conf.d/

```
