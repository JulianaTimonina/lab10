# Лабораторная работа №10
## Вариант 9

### Задания

| Задание | Описание |
|---------|----------|
| М1 | 1. Создать простое API на Go (Gin) с 2-3 эндпоинтами |
| М2 | 3. Реализовать валидацию входных данных в Go |
| М3 | 5. Передавать сложные структуры данных (JSON) между сервисами |
| H1 | 3. Добавить аутентификацию (JWT) в Go-сервисе и проверять токены из Python + 5. Развернуть оба сервиса в Docker Compose с общей сетью. |

---
```
В h1 выполнены задания 3 и 5 из повышенной сложности,
так как выполнение задания 5 предполагает выполнение задания 3.
```

### Структура репозитория

```
lab9/
├── README.md
├── .gitignore
├── m1/
│   ├── go.sum
│   ├── go.mod
│   ├── main.go
│   ├── main_test.go
│   ├── handlers/
│   │   └── task.go
│   ├── middleware/
│   │   └── logger.go
│   └── models/
│       └── task.go
│
├── m2/
│   ├── go.mod
│   ├── go.sum
│   ├── validator.go
│   └── validator_test.go
│
├── m3/
│   ├── fastapi-service/
│   │   ├── app/
│   │   │   ├── main.py
│   │   │   ├── models.py
│   │   │   └── client.py
│   │   ├── tests/
│   │   │   ├── test_main.py
│   │   │   └── test_client.py
│   │   └── pyproject.toml
│   └── gin-service/
│       ├── main.go
│       ├── go.mod
│       ├── go.sum
│       └── main_test.go 
│   
└── h1/
    ├── docker-compose.yml	#задание 5
    ├── go-service/
    │   ├── auth/
    │   │   └── jwt.go
    │   ├── handlers/
    │   │   └── handlers.go
    │   ├── middleware/
    │   │   └── auth.go
    │   ├── main.go
    │   ├── main_test.go
    │   ├── Dockerfile	#задание 5
    │   ├── go.mod
    │   └── go.sum
    └── python-service/
        ├── auth/
        │   ├── __init__.py
        │   └── jwt_validator.py
        ├── tests/
        │   ├── test_jwt_validator.py
        │   └── test_main.py
        ├── Dockerfile	#задание 5
        ├── main.py
        ├── requirements.txt
        └── pyproject.toml 

```

## Запуск тестов

### М1
Создать простое API на Go (Gin) с 2-3 эндпоинтами

```bash
# Запуск
go mod tidy
go run main.go

# Запуск тестов
go test -v
```

#### М2
Реализовать валидацию входных данных в Go

```bash
# Запуск и тестирование
go mod tidy
go test -v
go test -bench=. -benchmem
```

#### М3
Передавать сложные структуры данных (JSON) между сервисами

```bash
# Запуск GO
cd gin-service
go mod tidy

go run main.go

# Запуск Python
cd fastapi-service
pip install -e .

uvicorn app.main:app --reload --port 8000

#Для проверки работоспособности
curl http://localhost:8080/health
curl http://localhost:8000/health


#Запуск тестов
cd gin-service
go test -v

cd fastapi-service
pytest tests/ -v
```

#### h1
Добавить аутентификацию (JWT) в Go-сервисе и проверять токены из Python
+
Развернуть оба сервиса в Docker Compose с общей сетью. 

```bash
# Запуск GO
cd go-service
go mod tidy

go run main.go

# Запуск Python
cd python-service
pip install -r requirements.txt

uvicorn main:app --reload

#Для проверки работоспособности
curl http://localhost:8080/public
curl http://localhost:8000/health


#Запуск тестов
cd go-service
go test -v

cd python-service
pytest tests/ -v

#Запуск в Docker Compose
docker-compose up -d
docker-compose ps
```