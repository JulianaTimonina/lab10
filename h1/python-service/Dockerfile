# Используем официальный образ Python
FROM python:3.11-slim

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы с зависимостями
COPY requirements.txt .

# Устанавливаем зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Копируем весь исходный код
COPY . .

# Открываем порт 8000
EXPOSE 8000

# Команда для запуска сервиса
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]