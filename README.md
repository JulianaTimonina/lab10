# Лабораторная работа №10
## Вариант 9

### Задания

| Задание | Описание |
|---------|----------|
| М1 | 1. Создать простое API на Go (Gin) с 2-3 эндпоинтами |
| М2 | 3. Реализовать валидацию входных данных в Go |
| М3 | 5. Передавать сложные структуры данных (JSON) между сервисами |

---

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
├── h1/
│   ├── 
│   
└── h2/
    ├── 
```

## Запуск тестов

### М1
Создать простое API на Go (Gin) с 2-3 эндпоинтами

```bash
# Запуск
go mod tidy
go run main.go

# Запуск тестов
go test -v ./...
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