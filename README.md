# Balance-avito API
путь к секретам от бд можно посмотреть в models/config.go

###### **Get /metrics**

Метрики


###### **GET /api/v1/balance**

Выводит баланс пользователя по указанному id и currency <br>
Пример запроса : http://localhost:8080/api/v1/balance?id=1&currency=USD
конвертер валют использует данные из Европейского центробанка
: <br>
```json
{
    "amount":1840.232288037166,
    "currency":"USD"
}
```

###### **GET /api/v1/transactions**

Выводит транзакции пользователя по указанному id, количество пользователей надо указать в параметре count <br>
пример запроса : http://localhost:8080/api/v1/transactions?id=1&count=3&sort=amount&onSort=DESC  <br>
также можно указать по чему сортировать в параметре sort и как сортировать в параметре onSort: <br>
```json
[
  {
    "user_id":1,
    "comment":"Начисление средств",
    "amount":23000,
    "date":"2020-09-27T15:06:03.137484+03:00"
    },
  {
    "user_id":1,
    "comment":"Начисление средств",
    "amount":11000,
    "date":"2020-09-27T00:23:08.527134+03:00"
  },
  {
    "user_id":1,
    "comment":"Перевод средств",
    "amount":11000, 
    "date":"2020-09-27T00:23:08.527134+03:00"
  }
]
```


###### **POST /api/v1/balance**

Создает пользователя с балансом  <br>
на входе:
 ```json
  {
    "user_id":1, 
    "amount":10000
  }
```


###### **PUT /api/v1/accrual**

Начисляет пользователю указанную сумму, а также записывает транзакцию  <br>
на входе:
 ```json
  {
      "user_id":1, 
      "amount":10000
  }
```

###### **PUT /api/v1/debit**

Списывает у пользователя указанную сумму, а также записывает транзакцию  <br>
на входе:
 ```json
  {
      "user_id":1, 
      "amount":10000
  }
```

###### **PUT /api/v1/transfer/{id:[0-9]+}**

Списывает у пользователя в теле запроса и начисляет пользователю в урле, также записывает транзакции  <br>
на входе:
 ```json
  {
      "user_id":1, 
      "amount":10000
  }
```

###### **DELETE /api/v1/balance/{id:[0-9]+}**

Удаляет пользователя по id  <br>




  