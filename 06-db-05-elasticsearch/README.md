# Домашнее задание к занятию "6.5. Elasticsearch"

## Модуль 6. Администрирование баз данных

### Студент: Иван Жиляев

## Задача 1

>В этом задании вы потренируетесь в:
>- установке elasticsearch
>- первоначальном конфигурировании elastcisearch
>- запуске elasticsearch в docker
>
>Используя докер образ [centos:7](https://hub.docker.com/_/centos) как базовый и 
>[документацию по установке и запуску Elastcisearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/targz.html):
>
>- составьте Dockerfile-манифест для elasticsearch
>- соберите docker-образ и сделайте `push` в ваш docker.io репозиторий
>- запустите контейнер из получившегося образа и выполните запрос пути `/` c хост-машины
>
>Требования к `elasticsearch.yml`:
>- данные `path` должны сохраняться в `/var/lib`
>- имя ноды должно быть `netology_test`
>
>В ответе приведите:
>- текст Dockerfile манифеста
>- ссылку на образ в репозитории dockerhub
>- ответ `elasticsearch` на запрос пути `/` в json виде
>
>Подсказки:
>- возможно вам понадобится установка пакета perl-Digest-SHA для корректной работы пакета shasum
>- при сетевых проблемах внимательно изучите кластерные и сетевые настройки в elasticsearch.yml
>- при некоторых проблемах вам поможет docker директива ulimit
>- elasticsearch в логах обычно описывает проблему и пути ее решения
>
>Далее мы будем работать с данным экземпляром elasticsearch.

### Решение

Для создания образа подготовлен [Dockerfile](Dockerfile). Команды для создания образа и отправки его в репозиторий:

```
docker build -f Dockerfile -t my_elasticsearch .
docker tag my_elasticsearch nimlock/netology-homework-6.5
docker push nimlock/netology-homework-6.5
```

Чтобы выполнить требования к заданию в Dockerfile была определена переменная окружения `ES_PATH_DATA`, которая в дальнейшем указывается в конфигурационном файле `elasticsearch.yml`. Имя ноды будет браться из переменной `HOSTNAME`.

Для возможности определять конфигурацию сервиса без подключения в контейнер поднимем временный контейнер, скопируем из него конфиг-файлы и удалим временный контейнер. В дальнейшнем эти файлы будут монтироваться поверх файлов в рабочем контейнере.

_Прим.: Более хорошего/удобного решения для возможности управлять конфигами приложения в контейнере не знаю._

```
docker run -d --name my_elasticsearch_temp nimlock/netology-homework-6.5
docker cp my_elasticsearch:/app/elasticsearch-7.10.0/config/elasticsearch.yml ./config
docker cp my_elasticsearch:/app/elasticsearch-7.10.0/config/jvm.options ./config
docker cp my_elasticsearch:/app/elasticsearch-7.10.0/config/log4j2.properties ./config
docker rm -f my_elasticsearch_temp
```

Используя информацию из [оф.документации](https://www.elastic.co/guide/en/elasticsearch/reference/current/bootstrap-checks.html) для работы сервиса в standalone-режиме с возможностью доступа к API по "non-localhost"-адресам нужно сделать две вещи:

- в [elasticsearch.yml](config/elasticsearch.yml) задать опции `discovery.type: single-node` и `network.host: 0.0.0.0`
- отключить часть проверок при запуске с помощью определения переменной `es.enforce.bootstrap.checks=true` в [docker-compose.yml](docker-compose.yml)

_Прим: Эта настройка не используется в задании, но может быть полезна в будущем._

Запустим сервис задав его параметры с помощью манифеста [docker-compose](docker-compose.yml).

```
docker-compose up -d
```

### Ответы

- ссылка на [Dockerfile](Dockerfile)
- ссылка на [образ в DockerHub](https://hub.docker.com/repository/docker/nimlock/netology-homework-6.5)
- ответ `elasticsearch` на запрос пути `/` в json виде:

    ```
    ivan@kubang:~/study/netology-virt/06-db-05-elasticsearch$ curl -X GET localhost:9200/
    {
    "name" : "netology_test",
    "cluster_name" : "elasticsearch",
    "cluster_uuid" : "GziYkU7TQaS0ytXVMxAlMA",
    "version" : {
        "number" : "7.10.0",
        "build_flavor" : "default",
        "build_type" : "tar",
        "build_hash" : "51e9d6f22758d0374a0f3f5c6e8f3a7997850f96",
        "build_date" : "2020-11-09T21:30:33.964949Z",
        "build_snapshot" : false,
        "lucene_version" : "8.7.0",
        "minimum_wire_compatibility_version" : "6.8.0",
        "minimum_index_compatibility_version" : "6.0.0-beta1"
    },
    "tagline" : "You Know, for Search"
    }
    ```

## Задача 2

>В этом задании вы научитесь:
>- создавать и удалять индексы
>- изучать состояние кластера
>- обосновывать причину деградации доступности данных
>
>Ознакомтесь с [документацией](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html) 
>и добавьте в `elasticsearch` 3 индекса, в соответствии со таблицей:
>
>| Имя | Количество реплик | Количество шард |
>|-----|-------------------|-----------------|
>| ind-1| 0 | 1 |
>| ind-2 | 1 | 2 |
>| ind-3 | 2 | 4 |
>
>Получите список индексов и их статусов, используя API и **приведите в ответе** на задание.
>
>Получите состояние кластера `elasticsearch`, используя API.
>
>Как вы думаете, почему часть индексов и кластер находится в состоянии yellow?
>
>Удалите все индексы.
>
>**Важно**
>
>При проектировании кластера elasticsearch нужно корректно рассчитывать количество реплик и шард,
>иначе возможна потеря данных индексов, вплоть до полной, при деградации системы.

### Решение

Запросы на создание индексов:

```
curl -X PUT "localhost:9200/ind-1?pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  }
}
'
curl -X PUT "localhost:9200/ind-2?pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "number_of_shards": 2,
    "number_of_replicas": 1
  }
}
'
curl -X PUT "localhost:9200/ind-3?pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "number_of_shards": 4,
    "number_of_replicas": 2
  }
}
'
```

Список индексов и их статусы можно получить следующим запросом:

```
ivan@kubang:~/study/netology-virt/06-db-05-elasticsearch$ curl -X GET "localhost:9200/_cat/indices?pretty"
green  open ind-1 UNBbTycwR7q6zLr_u0xrXg 1 0 0 0 208b 208b
yellow open ind-3 gILROa3GQ6G7Jw0Vk68b4w 4 2 0 0 832b 832b
yellow open ind-2 v_wj-gWJQROwaT0rKuy4KA 2 1 0 0 416b 416b
```

Состояние кластера проверим запросом:

```
ivan@kubang:~/study/netology-virt/06-db-05-elasticsearch$ curl -X GET "localhost:9200/_cluster/health?pretty"
{
  "cluster_name" : "elasticsearch",
  "status" : "yellow",
  "timed_out" : false,
  "number_of_nodes" : 1,
  "number_of_data_nodes" : 1,
  "active_primary_shards" : 7,
  "active_shards" : 7,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 10,
  "delayed_unassigned_shards" : 0,
  "number_of_pending_tasks" : 0,
  "number_of_in_flight_fetch" : 0,
  "task_max_waiting_in_queue_millis" : 0,
  "active_shards_percent_as_number" : 41.17647058823529
}
```

Индексы `ind-2` и `ind-3` находятся в статусе `yellow` из-за того, что созданные реплики шард не могут "развернуться" так как у нас в кластере только одна нода (а на ней уже развёрнуты primary-шарды). Для решения проблемы нужно добавить ещё две ноды т.к. для `ind-3` задано две реплики и именно этот индекс диктует минимальное количество нод для работы без предупреждений.  
Весь кластер находится в статусе `yellow` так как часть индексов, расположенных в нём, находится в этом же статусе.

Удалить индексы можно запросами:

```
curl -X DELETE "localhost:9200/ind-1?pretty"
curl -X DELETE "localhost:9200/ind-2?pretty"
curl -X DELETE "localhost:9200/ind-3?pretty"
```

## Задача 3

>В данном задании вы научитесь:
>- создавать бэкапы данных
>- восстанавливать индексы из бэкапов
>
>Создайте директорию `{путь до корневой директории с elasticsearch в образе}/snapshots`.
>
>Используя API [зарегистрируйте](https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-register-repository.html#snapshots-register-repository) 
>данную директорию как `snapshot repository` c именем `netology_backup`.
>
>**Приведите в ответе** запрос API и результат вызова API для создания репозитория.
>
>Создайте индекс `test` с 0 реплик и 1 шардом и **приведите в ответе** список индексов.
>
>[Создайте `snapshot`](https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-take-snapshot.html) 
>состояния кластера `elasticsearch`.
>
>**Приведите в ответе** список файлов в директории со `snapshot`ами.
>
>Удалите индекс `test` и создайте индекс `test-2`. **Приведите в ответе** список индексов.
>
>[Восстановите](https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshots-restore-snapshot.html) состояние
>кластера `elasticsearch` из `snapshot`, созданного ранее. 
>
>**Приведите в ответе** запрос к API восстановления и итоговый список индексов.
>
>Подсказки:
>- возможно вам понадобится доработать `elasticsearch.yml` в части директивы `path.repo` и перезапустить `elasticsearch`

### Решение

Создадим каталог для снапшотов:

```
docker exec 06-db-05-elasticsearch_elastic-server_1 mkdir snapshots
```

Для удобства создадим "обезличенный" симлинк `elasticsearch-current` на рабочий каталог текущего релиза. Это позволит в будущем при изменении версии ПО не менять запросы к api. Если идея окажется полезной, то это изменение можно будет зафиксировать в Dockerfile.

```
docker exec 06-db-05-elasticsearch_elastic-server_1 bash -c 'cd .. && ln -snf elasticsearch-${ELASTIC_VER} elasticsearch-current'
```

Добавим путь к снапшотам в файл конфигурации `elasticsearch.yml`:

```
path:
  repo:
    - /app/elasticsearch-current/snapshots
```

Для применения изменений перезагрузим контейнер:

```
docker restart 06-db-05-elasticsearch_elastic-server_1
```

Теперь можно подключиться к контейнеру и зарегистрировать репозиторий:

```
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X PUT "localhost:9200/_snapshot/netology_backup?pretty" -H 'Content-Type: application/json' -d'
> {
>   "type": "fs",
>   "settings": {
>     "location": "/app/elasticsearch-current/snapshots"
>   }
> }
> '
{
  "acknowledged" : true
}
```

Создадим индекс `test`:

```
curl -X PUT "localhost:9200/test?pretty" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  }
}
'
```

Список индексов получим запросом:

```
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X GET "localhost:9200/_cat/indices?pretty"
green open test bCQDEg1mSH6a17rNb8uucw 1 0 0 0 208b 208b
```

При создании снапшота без указания параметров операции происходит бэкап всех данных в кластере, что нам и требуется. Так что запрос на создание снапшота будет таким:

```
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X PUT "localhost:9200/_snapshot/netology_backup/snapshot_1?wait_for_completion=true&pretty"
{
  "snapshot" : {
    "snapshot" : "snapshot_1",
    "uuid" : "5C50Sp9kSeCll7KWmaCyzQ",
    "version_id" : 7100099,
    "version" : "7.10.0",
    "indices" : [
      "test"
    ],
    "data_streams" : [ ],
    "include_global_state" : true,
    "state" : "SUCCESS",
    "start_time" : "2020-11-24T15:25:17.704Z",
    "start_time_in_millis" : 1606231517704,
    "end_time" : "2020-11-24T15:25:17.908Z",
    "end_time_in_millis" : 1606231517908,
    "duration_in_millis" : 204,
    "failures" : [ ],
    "shards" : {
      "total" : 1,
      "failed" : 0,
      "successful" : 1
    }
  }
}
```

Посмотрим на содержимое папки со снапшотами:

```
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ ls -la snapshots/
total 60
drwxr-xr-x 3 elasticuser elasticuser  4096 Nov 24 20:25 .
drwxr-xr-x 1 elasticuser root         4096 Nov 24 19:39 ..
-rw-r--r-- 1 elasticuser elasticuser   434 Nov 24 20:25 index-0
-rw-r--r-- 1 elasticuser elasticuser     8 Nov 24 20:25 index.latest
drwxr-xr-x 3 elasticuser elasticuser  4096 Nov 24 20:25 indices
-rw-r--r-- 1 elasticuser elasticuser 30685 Nov 24 20:25 meta-5C50Sp9kSeCll7KWmaCyzQ.dat
-rw-r--r-- 1 elasticuser elasticuser   266 Nov 24 20:25 snap-5C50Sp9kSeCll7KWmaCyzQ.dat
```

Удалим индекс `test` и создадим индекс `test-2` с последующим выводом списка индексов:

```
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X DELETE "localhost:9200/test?pretty"
{
  "acknowledged" : true
}
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X PUT "localhost:9200/test-2?pretty" -H 'Content-Type: application/json' -d'
> {
>   "settings": {
>     "number_of_shards": 1,
>     "number_of_replicas": 0
>   }
> }
> '
{
  "acknowledged" : true,
  "shards_acknowledged" : true,
  "index" : "test-2"
}
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X GET "localhost:9200/_cat/indices?pretty"
green open test-2 2xUhI231QD60CC5NgCtgaQ 1 0 0 0 208b 208b
```

Восстановим состояние кластера из снапшота:

```
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X POST "localhost:9200/_snapshot/netology_backup/snapshot_1/_restore?pretty"
{
  "accepted" : true
}
```

Итоговый список индексов содержит как созданный `test-2`, так и восстановленный `test` (но с другим хэшем):

```
[elasticuser@40d64c412cf2 elasticsearch-7.10.0]$ curl -X GET "localhost:9200/_cat/indices?pretty"
green open test-2 2xUhI231QD60CC5NgCtgaQ 1 0 0 0 208b 208b
green open test   8hqNPscQS7CTN6mn-9jRvw 1 0 0 0 208b 208b
```
