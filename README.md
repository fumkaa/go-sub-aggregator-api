# Subscription Aggregator API

## Стек технологий

- Framework: Chi
- Database: PostgreSQL
- Logging: Slog
- Validation: Go-playground/validator & Gorilla Schema
- Config: Cleanenv (YAML/ENV)
- Documentation: OpenAPI 3.0 / Swagger

# Архитектура проекта

```
sub-aggregator/
├── internal/
│   ├── app/                     # [Layer: Orchestrator] Сборка приложения
│   │   ├── http/                
│   │   │   └── app.go           // Инициализация HTTP-сервера и управление его жизненным циклом
│   │   └── app.go               // Главная точка сборки всех компонентов
│   ├── config/                  
│   │   └── config.go            // Загрузка и парсинг конфигурации (библиотека cleanenv)
│   ├── domain/                  # [Layer: Domain] Модели данных и ошибки
│   │   ├── models/              // Бизнес-объекты (Subscription)
│   │   │   └── subscription.go
│   │   └── errors.go            // Глобальные ошибки домена (например, ErrSubNotFound)
│   ├── lib/                     // Общие вспомогательные утилиты
│   │   └── logger/              
│   │       └── sl/
│   │           └── sl.go        // Настройка structured logging (slog)
│   ├── repository/              # [Layer: Repository] Работа с данными
│   │   └── postgres/            // Реализация хранилища для PostgreSQL (SQL запросы через pgx)
│   │       ├── postgres.go      // Инициализация пула соединений
│   │       └── subscription.go  // CRUD операции и агрегация стоимости подписок
│   ├── services/                # [Layer: Usecase] Бизнес-логика
│   │   └── subscription/        
│   │       └── service.go       // Логика работы с подписками
│   └── transport/               # [Layer: Delivery] Транспортный уровень, внешние интерфейсы
│       └── http/                
│           ├── handler.go       // Обработчики запросов 
│           ├── router.go        // Описание маршрутов и подключение middleware
│           ├── server.go        
│           └── middleware/      
│               └── validator.go // Валидация запросов
├── cmd/
│   ├── sub-aggregator/          # Точка входа для основного приложения
│   │   └── main.go              // Инициализация конфига, логгера, БД и запуск сервера
│   └── migrator/                # Утилита для выполнения миграций БД
│       └── main.go              // Логика подключения к БД и запуска файлов миграций
├── configs/
│   └── local.yaml               // Конфигурация для локальной разработки
├── migrations/                  # SQL-скрипты для управления схемой базы данных
│   ├── 000001_create_subscriptions_table.up.sql
│   ├── 000001_create_subscriptions_table.down.sql
│   ├── 000002_create_indexes.up.sql
│   ├── 000002_create_indexes.down.sql
│   ├── 000003_create_trigger_modtime.up.sql
│   └── 000003_create_trigger_modtime.down.sql
├── swagger/                     # Документация API
│   └── openapi.yaml             
├── .env
├── .gitignore                   
├── go.mod
├── go.sum                       
└── README.md                    
```