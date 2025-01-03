# REST API для управления заметками

## Обзор
Этот проект представляет собой REST API для управления заметками, написанный на **Golang**. Используется слоистая архитектура для улучшения модульности и масштабируемости. В проекте применяются следующие технологии:

- **[FX](https://github.com/uber-go/fx)**: Фреймворк для внедрения зависимостей.
- **[Gorm](https://gorm.io/)**: ORM-библиотека для работы с базой данных.
- **[Zap](https://github.com/uber-go/zap)**: Библиотека для структурированного логирования.
- **[Viper](https://github.com/spf13/viper)**: Управление конфигурацией.
- **[Redis](https://redis.io/)**: Хранилище данных в памяти для кэширования.

## Функционал
- CRUD-операции (создание, чтение, обновление, удаление) для заметок.
- Масштабируемая, модульная архитектура с соблюдением принципов "чистого кода".
- Структурированное логирование с помощью Zap для упрощения отладки.
- Кэширование с Redis для повышения производительности.
- Настраиваемость через переменные окружения с использованием Viper.

## Структура проекта
Проект организован по принципу **слоистой архитектуры** для лучшей изоляции ответственности:

```plaintext
├── cmd               # Точка входа приложения
│   └── api           # Основной исполняемый файл
├── internal          # Внутренняя логика приложения
│   ├── config        # Управление конфигурацией с Viper
│   ├── models        # Определение моделей данных
│   ├── repository    # Логика доступа к данным
│   ├── service       # Бизнес-логика
│   ├── transport     # HTTP-обработчики для API
│   └── migrations    # Миграции базы данных
├── pkg               # Общие утилиты и вспомогательные модули
│   ├── cache         # Логика работы с Redis
│   └── logger        # Настройка логирования с Zap
├── .env              # Переменные окружения
├── docker-compose.yml# Настройка контейнеров Docker
├── go.mod            # Зависимости Go
├── Makefile          # Утилиты для сборки и запуска
└── main.go           # Инициализация приложения
```

### Описание слоёв
1. **Config**: Управление конфигурацией с использованием Viper.
2. **Logger**: Настройка структурированного логирования с Zap.
3. **Models**: Определение сущностей данных, используемых в приложении.
4. **Repository**: Реализация доступа к данным через Gorm.
5. **Service**: Бизнес-логика приложения.
6. **Transport**: HTTP-обработчики для маршрутов API.
7. **Cache**: Кэширование с использованием Redis.

## Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/removespread/todo-api.git
   cd todo-api
   ```

2. Установите зависимости:
   ```bash
   go mod tidy
   ```

3. Настройте переменные окружения:
   Необходимо отредактировать `.env` файл и выставить следующие параметры:
   Переменные выставлены как пример
   ```env
    HTTP_PORT=:8080
    POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
    REDIS_DSN=redis://localhost:6379/0
   ```

4. Запустите приложение:
   ```bash
   go run main.go
   ```

