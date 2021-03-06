# Домашнее задание к занятию "6.1. Типы и структура СУБД"

## Модуль 6. Администрирование баз данных

### Студент: Иван Жиляев

## Задача 1

>Архитектор ПО решил проконсультироваться у вас, какой тип БД 
>лучше выбрать для хранения определенных данных.
>
>Он вам предоставил следующие типы сущностей, которые нужно будет хранить в БД:
>
>- Электронные чеки в json виде
>- Склады и автомобильные дороги для логистической компании
>- Генеалогические деревья
>- Кэш идентификаторов клиентов с ограниченным временем жизни для движка аутенфикации
>- Отношения клиент-покупка для интернет-магазина
>
>Выберите подходящие типы СУБД для каждой сущности и объясните свой выбор.

Рассмотрим предлагаемые сущности:
>- Электронные чеки в json виде

Для этого случая я бы предложил СУБД документо-ориентированного типа: он позволит помещать чеки в БД без какого-либо преобразования, все поля json перенесутся в базу и запись о документе будет дополнена его метаданными.

>- Склады и автомобильные дороги для логистической компании

В этом случае я бы предложил графовую СУБД: в узлах будут склады, а дороги - их отношениями. 

>- Генеалогические деревья

Из списка данного на лекции больше всего подходит сетевая модель СУБД: она также как и генеалогические деревья иерархична и позволяет каждому узлу присвоить отношения с несколькими объектами.

>- Кэш идентификаторов клиентов с ограниченным временем жизни для движка аутенфикации

Пожалуй, эта задача для баз ключ-значение, работающих в оперативной памяти: небольшой TTL и ожидаемо небольшие величины в каждой записи не позволят базе сильно разрастись, в базах этого типа есть встроенный механизм работы с TTL. Да и сама специфика данных, как мне кажется, вполне допускает пониженную доступность, характеризующую эти СУБД.

>- Отношения клиент-покупка для интернет-магазина

Эти данные стоит хранить в реляционных СУБД: здесь важна и максимальная согласованность, и возможность работать с транзакциями. К тому же и у клиента, и у его покупки может быть большое количество дополнительных свойств, которые будет удобно вынести в отдельные таблицы - так управлять базой будет гораздо проще.

## Задача 2

>Вы создали распределенное высоконагруженное приложение и хотите классифицировать его согласно 
>CAP-теореме. Какой классификации по CAP-теореме соответствует ваша система, если
>(каждый пункт - это отдельная реализация вашей системы и для каждого пункта надо привести классификацию):
>
>- Данные записываются на все узлы с задержкой до часа (асинхронная запись)
>- При сетевых сбоях, система может разделиться на 2 раздельных кластера
>- Система может не прислать корректный ответ или сбросить соединение
>
>А согласно PACELC-теореме, как бы вы классифицировали данные реализации?

### CAP-теорема 

>- Данные записываются на все узлы с задержкой до часа (асинхронная запись)

Этот пункт явно не соответствует понятию `Согласованности`, т.к. после завершения транзакции не все узлы будут иметь новые данные (до завершения очередной синхронизации) и, соответственно, будут какое-то время предоставлять неактуальные ответы на запросы.  
Система соответствует __AP-типу__.

>- При сетевых сбоях, система может разделиться на 2 раздельных кластера

Условие сообщает, что реализация обладает свойством `Устойчивость к разделению`. Подразумевается также её `Доступность`. Но вот `Согласованность` в данном случае невозможна - у нас два разделённых экземпляра системы.  
Система снова соответствует __AP-типу__.

>- Система может не прислать корректный ответ или сбросить соединение

Эта реализация не обладает свойством `Доступность`.  
Получается, что наша система соответствует __CP-типу__.

### PACELC-теорема

>- Данные записываются на все узлы с задержкой до часа (асинхронная запись)

В этой реализации нет `Согласованности`, так что если она будет разделена, то можно предполагать лишь её `Доступность`. В нормальном режиме функционирования у системы очевиден приоритет на небольшие задержки при обращениях.  
Система соответствует __типу AP / EL__.

>- При сетевых сбоях, система может разделиться на 2 раздельных кластера

Думаю, что разделение на кластеры подразумевает назначение в каждом из них своего master-а, так что в каждом будет обеспечена своя консистентность. Поведение системы в неразделённом состоянии не описано в условии, так что предположу, что при наличии возможности к разделению на кластеры было бы ценно поддерживать целостность и в нормальных условиях.  
Система соответствует __типу CP / EС__.

>- Система может не прислать корректный ответ или сбросить соединение

У реализации всё плохо с `Доступностью`. Как правило системы жертвуют ей в пользу прочих качеств.  
Система соответствует __типу CP / EС__.

## Задача 3

>Могут ли в одной системе сочетаться принципы BASE и ACID? Почему?

Мне кажется, что такое сочетание в одной системе возможно только если мы будем говорить про разные режимы работы этой системы. То есть, когда система не разделена, то она будет работать по принципам ACID, но во время аварий (разделения) будет переходить на принципы BASE. Например, если система будет стараться придерживаться принципа `Консистентности` в разделённом состоянии, то она просто перестанет работать до восстановления целостности. Зато, если при разделении будет подразумеваться `Конечная согласованность`, то всё будет в порядке.

## Задача 4

>Вам дали задачу написать системное решение, основой которого бы послужили:
>
>- фиксация некоторых значений с временем жизни
>- реакция на истечение таймаута
>
>Вы слышали о key-value хранилище, которое имеет механизм [Pub/Sub](https://habr.com/ru/post/278237/). 
>Что это за система? Какие минусы выбора данной системы?

Из продуктов в приведённой статье key-value хранилищем является только Redis.

Из плюсов - в нём есть встроенная поддержка таймеров для записей в хранилище. Однако реакцию на истечение таймера нужно будет реализовывать программно, сам Redis просто удаляет значение.   
Из минусов - Redis хранит все данные в оперативной памяти, так что потеря питания приведёт к потере данных.
