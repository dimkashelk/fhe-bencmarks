.PHONY: all build run clean help deps check

# Имя выходного бинарника
BINARY_NAME=lattigo-bench
OUTPUT_FILE=results.csv

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

all: deps build run ## Установить зависимости, собрать и запустить бенчмарки

help: ## Показать справку по доступным командам
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

check: ## Проверить требования (Go версия)
	@echo "$(GREEN)Проверка требований...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo "Ошибка: Go не установлен. Установите Go 1.20+"; exit 1; }
	@echo "Go версия: $$(go version)"
	@echo "✓ Все требования выполнены"

deps: check ## Установить зависимости
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	@if [ ! -f "go.mod" ]; then \
		echo "Инициализация go.mod..."; \
		go mod init bench; \
		go mod edit -replace github.com/tuneinsight/lattigo/v6=../; \
	fi
	@echo "Загрузка зависимостей..."
	@go mod tidy
	@echo "✓ Зависимости установлены"

build: deps ## Собрать бенчмарк
	@echo "$(GREEN)Сборка $(BINARY_NAME)...$(NC)"
	@go build -o $(BINARY_NAME) .
	@echo "✓ Сборка завершена: $(BINARY_NAME)"

run: build ## Запустить бенчмарки
	@echo "$(GREEN)Запуск бенчмарков...$(NC)"
	@echo "Это может занять 2-5 минут в зависимости от вашего процессора..."
	@./$(BINARY_NAME)
	@echo ""
	@echo "$(GREEN)✓ Бенчмарки завершены!$(NC)"
	@echo "Результаты сохранены в: $(OUTPUT_FILE)"
	@echo ""
	@echo "Первые 5 строк результатов:"
	@head -n 5 $(OUTPUT_FILE)

run-only: ## Запустить бенчмарки без пересборки
	@echo "$(GREEN)Запуск бенчмарков...$(NC)"
	@if [ ! -f "$(BINARY_NAME)" ]; then \
		echo "Ошибка: $(BINARY_NAME) не найден. Запустите 'make build' сначала."; \
		exit 1; \
	fi
	@./$(BINARY_NAME)
	@echo "$(GREEN)✓ Результаты сохранены в: $(OUTPUT_FILE)$(NC)"

view: ## Показать результаты
	@if [ -f "$(OUTPUT_FILE)" ]; then \
		cat $(OUTPUT_FILE); \
	else \
		echo "Файл $(OUTPUT_FILE) не найден. Запустите 'make run' сначала."; \
	fi

clean: ## Удалить собранные файлы и результаты
	@echo "$(GREEN)Очистка...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f $(OUTPUT_FILE)
	@echo "✓ Очистка завершена"

clean-all: clean ## Удалить всё включая go.mod и зависимости
	@echo "$(GREEN)Полная очистка...$(NC)"
	@rm -f go.mod go.sum
	@echo "✓ Полная очистка завершена"

test: ## Проверить что код компилируется
	@echo "$(GREEN)Проверка компиляции...$(NC)"
	@go build -o /dev/null .
	@echo "✓ Код успешно компилируется"

install: build ## Установить бинарник в GOPATH/bin
	@echo "$(GREEN)Установка $(BINARY_NAME) в GOPATH/bin...$(NC)"
	@go install
	@echo "✓ Установлено. Теперь можно запустить: $(BINARY_NAME)"

# Информация о проекте
info: ## Показать информацию о проекте
	@echo "$(GREEN)=== Lattigo Benchmark Tool ===$(NC)"
	@echo ""
	@echo "Этот инструмент измеряет производительность операций"
	@echo "гомоморфного шифрования для схем BFV, CKKS и BGV."
	@echo ""
	@echo "$(YELLOW)Измеряемые операции:$(NC)"
	@echo "  • Генерация ключей (релинеаризация, Galois)"
	@echo "  • Шифрование/Дешифрование"
	@echo "  • Сложение/Умножение"
	@echo "  • Релинеаризация"
	@echo "  • Вращение"
	@echo "  • Сериализация (обычная, ZLIB, Zstandard)"
	@echo ""
	@echo "$(YELLOW)Параметры:$(NC)"
	@echo "  • Размеры колец: 4096, 8192, 16384"
	@echo "  • Итераций на операцию: 100"
	@echo "  • Результаты в микросекундах (μs)"
	@echo ""
	@echo "$(YELLOW)Быстрый старт:$(NC)"
	@echo "  make all      # Всё в одной команде"
	@echo "  make help     # Список всех команд"
	@echo ""

