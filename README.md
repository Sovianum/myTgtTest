# myTgtTest
Написать HTTP сервис на языке Golang, который предоставляет API для отлеживания активности пользователей на сайте.

Сервис должен предоставлять следующий API:

1. Регистрация: для каждого пользователя должна быть доступна следующая информация: уникальный идентификатор, возраст, пол.
2. Добавление статистики для пользователя. Парамеры: дата и время, идентификатор пользователя, тип статистики.

Тип статистики может быть одним из: логин, лайк, комментарий, выход с сайта.

3. Статистика по дням:
Топ N пользователей, которые оставляют больше всего комментариев на указанный период времени, данные в ответе должны быть упорядочены по дням.

Пример API:
* Регистрация пользователя:
<br />POST /api/users 
<br />{"id":1, "age":20, "sex":"M"}

* Добавление статистики:
<br /> POST /api/users/stats 
<br /> {"user":1, "action":"like", "ts":"2017-06-30T14:12:34"}

* Отчеты:
<br /> GET /api/users/stats/top?date1=2017-06-20&date2=2017-06-30&action=comments&limit=10
<br /> {"items":{"date":"2017-06-20", "rows":[{"id":1, "age":20, "sex":"M", "count":100},...]}...}
