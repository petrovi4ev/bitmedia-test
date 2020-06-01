# Тестовое задание для Bitmedia.

- Установка `go get -u gitlab.com/petrovi4ev/bitmedia-test`
- Переходим в рабочую директорию `cd $GOPATH/src/gitlab.com/petrovi4ev/bitmedia-test`
- Копируем конфиг приложения `cp .env.example .env`
- В конфиге указываем параметры подключения к БД, порт на котором будет работать сервер.
- Выполняем `make build-run` для запуска приложения, если требуется заполнить БД предоставленными в задании тестовыми данными, то выполняем `make build-run-migrate`

## API:

#### Получение списка пользователей.
- URL: /users
- Метод: GET
- Дополнительные параметры:  
    `per-page` - количество элементов на странице,  
    `page` - номер нужной страницы
- Пример: `localhost:8080/users?per-page=4&page3`

#### Получение информации об одном пользователе.
- URL: /users/{id}
- Метод: GET
- Обязательные параметры:  
    `id` - id пользователя
- Пример: `localhost:8080/users/5ed3fcfe7c1cb71634268f46`

#### Создание пользователя.
- URL: /users
- Метод: POST
- Обязательные параметры:  
    JSON представление нового пользователя.
- Пример: 

```json
{
    "birth_date": "Tuesday, April 26, 7042 3:14 PM",
    "city": "Kyiv",
    "country": "Ukraine",
    "email": "johnroget@gmail.tech",
    "gender": "Male",
    "last_name": "Sergey"
}
```

#### Обновление пользователя.
- URL: /users/{id}
- Метод: PUT
- Обязательные параметры:  
    JSON представление полей, которые нужно изменить.
- Пример: 

```json
{
    "birth_date": "Tuesday, April 26, 7042 3:14 PM",
    "city": "London",
    "country": "UK",
    "email": "johnroget@gmail.tech",
    "gender": "Male",
    "last_name": "John"
}
```

#### Удаление пользователя.
- URL: /users/{id}
- Метод: DELETE
- Пример: `localhost:8080/users/5ed3fcfe7c1cb71634268f46`