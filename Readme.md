# Укорачиватель ссылок

## Что это?

Тестовое задание от Ozon Fintech, предполагающее небольшой сервис по укорачиванию ссылок. Сервис должен представлять из
себя Docker-образ, использовать in-memory или psql хранилище данных (выбирается при запуске).

## Как взаимодействовать?

1. Добавить .env в корень проекта со следующими аргументами:
    * POSTGRES_PORT
    * POSTGRES_DB
    * POSTGRES_USER
    * POSTGRES_PASSWORD
    * REDIS_PORT

2. Отправлять запросы:

   **Создание короткой ссылки:**
   ```http request
   POST /shortener/
   ```
   Тело запроса
   ```json
   {"url": "https://www.example.com/"}
   ```
   Возможные ответы:
   ```json lines
   200 - OK
   {"full_link":"https://www.example.com","short_link":"xR1BZQhdhG"}
   
   500 - Internal Server Error
   ```

   **Получение длинной ссылки из короткой**
   ```http request
   GET /shortener/{short_link}
   ```
   Возможные ответы:
   ```json lines
   200 - OK
   {"full_link":"https://www.example.com","short_link":"xR1BZQhdhG"}
   
   404 - Not Found
   500 - Internal server Error
   ```

## Как запустить?

```shell
git clone https://github.com/ronmount/ozon_go.git && cd ozon_go

# Запуск тестов:
make test

# Запуск с in-memory хранилищем:
make memory

# Запуск с хранилищем Postgres:
make postgresql
```

## Примеры запросов:

```shell
# Создание короткой ссылки
curl --request POST \
  --url http://localhost:8080/shortener/ \
  --data url=https://ozon.ru/
  
# Получение длинной ссылки
curl --request GET --url http://localhost:8080/shortener/XSHC0bWFda
```
