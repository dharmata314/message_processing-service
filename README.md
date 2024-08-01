# Backend сервис для обработки сообщений
## Развертывание
### Общие настройки
Общие настройки приложения содержатся в [конфиге](https://github.com/dharmata314/message_processing-service/tree/main/config). В зависимости от способа развертывания какие-либо параметры могут меняться. 
В конфиге содержатся основные данные, необходимые для работы приложения.

### Доступ к сервису по IP
Сервис доступен по адресу
```
87.242.106.218:8080 
```
Доступные следующие запросы к серверу


Регистрация пользователя:
```
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://87.242.106.218:8080 :8080/users/new
```
Авторизация пользователя:
```
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://87.242.106.218:8080 :8080/login
```

Изменение данных пользователя:
```
curl -X PATCH \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"email": "newEmail@email.com", "password": "NewPassword", "id": 1}' \
http://87.242.106.218:8080 :8080/users/{id}
```
Отправка сообщения:
```
curl -X POST \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"content": "test"}' \
http://87.242.106.218:8080 :8080/message
```
Получение статистика по сообщениям
```
curl -X GET \
-H "Authorization: Bearer <token>" \
http://87.242.106.218:8080 :8080/message/statistics
```
Удаление сообщения из базы данных:
```
curl -X DELETE \
-H "Authorization: Bearer <token>" \
http://87.242.106.218:8080 :8080/message/{id}
```


### Запуск приложения локально
Запуск приложения осуществляется через Docker

Для развертывания в Docker Compose создан файл [docker-compose.yml](https://github.com/dharmata314/message_processing-service/blob/main/docker-compose.yml)
Для запуска приложения необходимо запустить следующую команду из директории проекта
```
docker-compose up --build app
```
## Общее
Приложение представляет из себя сервис обработки сообщений в Kafka. 

Сервис устроен следующим образом:

Пользователь региструется в сервисе, авторизуется, затем отправляет сообщения в PostrgreSQL, после чего Kafka обрабатывает сообщения.

Доступны следующие операции:
 - Регистрация пользователя
 - Авторизация пользователя
 - Изменение данных пользователя
 - Добавление сообщений в базу данных
 - Получение статистики по обработанным сообщениям
 - Удаление сообщений из базы данных

Для доступа к большинству функционала (кроме регистрации и авторизации) необходим доступ по токену.
Токен выдается пользователю после авторизации.
В дальнейшем токен должен передаваться вместе с заголовком запроса:
```
Authorization: Bearer <token>
```

## Примеры запросов


Регистрация пользователя:
```
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://localhost:8080/users/new
```
Авторизация пользователя:
```
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://localhost:8080/login
```

Изменение данных пользователя:
```
curl -X PATCH \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"email": "newEmail@email.com", "password": "NewPassword", "id": 1}' \
http://localhost:8080/users/{id}
```
Отправка сообщения:
```
curl -X POST \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"content": "test"}' \
http://localhost:8080/message
```
Получение статистика по сообщениям
```
curl -X GET \
-H "Authorization: Bearer <token>" \
http://localhost:8080/message/statistics
```
Удаление сообщения из базы данных:
```
curl -X DELETE \
-H "Authorization: Bearer <token>" \
http://localhost:8080/message/{id}
```


