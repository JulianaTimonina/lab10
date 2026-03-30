# Лабораторная работа №10
## Вариант 9

### Задания

| Задание | Описание |
|---------|----------|
| М1 | 1. Создать простое API на Go (Gin) с 2-3 эндпоинтами |
| М2 | 3. Реализовать валидацию входных данных в Go |

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
├── m2/
│   ├── go.mod
│   ├── go.sum
│   ├── validator.go
│   └── validator_test.go
├── m3/
│   ├── 
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
