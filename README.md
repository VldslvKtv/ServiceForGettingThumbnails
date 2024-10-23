ОБЩАЯ СТРУКТУРА:

client - клиентская часть
config - файлы конфигурации для сервера и ключи api (нужно создать самому файл api.yaml с параметрами api_key: <ваш ключ YouTube Data API v3>
                                                                                                     api_url: "https://www.googleapis.com/youtube/v3/videos")
gen - сгенерированные файлы из proto
migrations - миграции для бд
proto - прото-файлы
tests - тесты

internal:
├───app
│   │   grpc_app.go - управление сервером
│   │
│   └───thumb_app - компановка сервера и приложения
│           app.go
│
├───cmd - запуск сервера
│       main.go
│
├───config - загрузка конфигурации проекта
│       config.go
│
├───grpc - 
│       server.go -gRPC-хэндлеры
│
├───loader - основная бизнес логика (сервисный слой)
│       loader.go
│
├───migrator - миграции в бд
│       main.go
│
├───models - модели для бд
│       thumbnail.go
│
└───youtube - получение превью
        youtube.go

storage: - работа сервера на уровне БД
│   sqlite.go
│   thumbnails.db
│
└───storage_err
        storage_err.go

Taskfile.yaml - нужно установить утилиту task чтобы пользоваться (пример: task run - для запуска сервера и тамкже по аналогии с миграциями и генерацией протофайлов)

сервер запускается из cmd/main.go
клиент запускается из client/client.go
Запуск нужно производить на разных терминалах

Также, можно дергать "ручки" из Postman

Пример как пользоваться client.go:
go run client.go https://www.youtube.com/watch?v=RMmjRrLdT3w

go run client.go --async  https://www.youtube.com/watch?v=RMmjRrLdT3w https://www.youtube.com/watch?v=72Ydy05emnY


