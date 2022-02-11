<h1>Укорачиватель ссылок</h1>

<h2>Что это?</h2>
<p>Тестовое задание от Ozon Fintech, предполагающее небольшой сервис по укорачиванию ссылок. 
Сервис должен представлять из себя Docker-образ, использовать in-memory или psql хранилище
данных (выбирается при запуске).</p>

<h2>Как взаимодействовать?</h2>
1. Добавить .env в корень проекта со следующими аргументами:
   * POSTGRES_PORT
   * POSTGRES_DB
   * POSTGRES_USER
   * POSTGRES_PASSWORD
   * REDIS_PORT
2. Отправлять запросы:
   <b>Создание короткой ссылки:</b>
   ```http request
    POST /shortener/
   ```
   Тело запроса
   ```json
   {"url": "https://www.example.com/"}
   ```
   Возможные ответы:
   ```json lines
   201 - Created
   {"full_link":"https://www.example.com","short_link":"xR1BZQhdhG"}
   
   409 - Conflict
   {"full_link":"https://www.example.com","short_link":"xR1BZQhdhG"}
   
   500 - Internal Server Error
   ```
   
   <b>Получение длинной ссылки из короткой</b>
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

<h2>Как запустить?</h2>
```shell
git clone https://github.com/ronmount/ozon_go.git && cd ozon_go

# Запуск с хранилищем Redis:
docker-compose --profile redis up
# Запуск с хранилищем Postgres:
docker-compose --profile postgresql up
```

<h2>Примеры запросов:</h2>
```shell
# Создание короткой ссылки
curl --request POST \
  --url http://localhost:8080/shortener/ \
  --data url=https://ozon.ru/
  
# Получение длинной ссылки
curl --request GET \ 
  --url http://localhost:8080/shortener/XSHC0bWFda
```
