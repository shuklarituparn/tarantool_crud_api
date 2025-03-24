# Tarantool API (Key-Value Store на Tarantool)

![tarantoolAPI](/static/image.png)

Простое key-value хранилище с REST API, построенное на Tarantool и Go. Документация Swagger и метрики используя prometheus/grafana включена.

## Особенности

- CRUD операции для ключей и значений
- Мониторинг метрик через Prometheus
- Логирование операций
- Swagger-документация
- Поддержка JSON-сериализации

## Требования

- Go 1.23+
- Tarantool 3.0+
- Docker
- Docker compose

## Установка

1.  Склонируйте репозиторий:

    ```bash
    git clone https://github.com/shuklarituparn/tarantool_crud_api.git

    ```

2.  Теперь выполните следующую команду, чтоы убедиться, что вы находитесь в корневой директории проекта:

    ```bash
    cd tarantool_crud_api
    ```

3.  Создайте .env файл:
    ```bash
    TARANTOOL_PASSWORD=<ваш пароль>
    LOG_LEVEL=INFO

        ```

    `убедитесь, что пароль в .env файле и в init.lua файле (tarantool/init.lua) совпадают`

4.  Находясь в в корневой директории проекта, выполните следующую команду, чтобы запустить:

```bash
    docker compose up
```

5.  Cервис доступен по адресу `localhost:5005`

6.  Вы можете получить доступ к метрикам prometheus по адресу `localhost:5005/metrics`

`Графана доступна по адресу localhost:3000`

---

## Использование API

- Сервер будет доступен на http://localhost:5005.

---

#### Создать запись

```bash
curl -X POST http://localhost:5005/api/v1/kv \
  -H "Content-Type: application/json" \
  -H 'accept: application/json' \
  -d '{"key": "test", "value": {"data": "example"}}'
```

#### Получить запись

```bash
curl -X GET http://localhost:5005/api/v1/kv/test \
-H 'accept: application/json'
```

#### Обновить запись

```bash
curl -X PUT http://localhost:5005/api/v1/kv/test \
  -H "Content-Type: application/json" \
  -H 'accept: application/json' \
  -d '{"value": {"data": "updated"}}'
```

#### Удалить запись

```bash
curl -X DELETE http://localhost:5005/api/v1/kv/test \
-H 'accept: application/json'
```

---

### Документация API

Доступна через Swagger UI: http://localhost:5005/docs/

### Мониторинг

Метрики Prometheus доступны по адресу:
http://localhost:5005/metrics

### Логирование

Все операции логируются в консоль с указанием:

    Времени выполнения

    Типа операции

    Ключа

    Статуса выполнения
