# mzhn сервис

Сервис объёдиняет в себе работу с базой знаний FAQ и сервис классификации

## Deploy

Перед деплоем нужно заполнить `.env` конфиг (пример можно взять из файла `example.env`)
```bash
cp example.env .env
```

Был описан docker compose файл со всем окружением необходимым для работы сервиса (postgres)
```
docker compose up --build
```

## Endpoints

### Predict

`POST /predict`
```json
{
	"question":"Как узнать версию приложения студия RUTUBE"
}
```
```json
{
	"answer": "Для этого в правом верхнем углу нажмите на фото своего профиля и под кнопкой \"Выйти\" будет указана версия приложения Студия RUTUBE.",
	"class_1": "ОТСУТСТВУЕТ",
	"class_2": ""
}
```

### FAQ CRUD

Загрузка ответа на вопрос(FAQ)

`POST /faq`
```json
{
  "question": "вопрос",
  "answer": "ответ",
  "classifier1": "УПРАВЛЕНИЕ АККАУНТОМ", // nullable
  "classifier2": "", // nullable
}
```
```json
{
  "id": "<uuid v4>"
}
```

`GET /faq` Query-параметры: пагинация `?limit` и `?offset`
```json
{
  "items": [
    {
      "id": "<uuid v4>"
      "question": "string",
      "answer": "string",
      "classifier1": "string?",
      "classifier2": "string?",
      "createdAt": "datetime",
      "updatedAt": "datetime?"
    }
  ],
  "total": "326" // Общее количество записей в таблице
}
```


`GET /faq/:id`
```json
{
  "id": "<uuid v4>"
  "question": "string",
  "answer": "string",
  "classifier1": "string?",
  "classifier2": "string?",
  "createdAt": "datetime",
  "updatedAt": "datetime?"
}
```
