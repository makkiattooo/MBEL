# Быстрый старт: Ваше первое приложение MBEL

Запустите готовое к производству многоязычное приложение с MBEL за **15 минут**. Это руководство охватывает настройку проекта, файлы переводов, компиляцию и развертывание.

---

## 1. Инициализируйте проект

```bash
# Создайте директорию проекта
mkdir -p hello-mbel && cd hello-mbel

# Инициализируйте модуль Go
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# Создайте структуру проекта
mkdir -p locales cmd dist

# Установите MBEL CLI
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

---

## 2. Напишите файлы переводов

### Создайте русские переводы (`locales/ru.mbel`)

```mbel
@namespace: hello
@lang: ru

app_name = "Привет MBEL"
app_version = "1.0.0"

greeting = "Добро пожаловать, {name}!"
goodbye = "До свидания, {name}! Время: {time}"

# Русский: [one] (1, 21, 31...), [few] (2-4, 22-24...), [many] (0, 5-20, 25-30...)
items_count(n) {
    [one]  => "У вас есть 1 элемент"
    [few]  => "У вас есть {n} элемента"
    [many] => "У вас есть {n} элементов"
    [other] => "У вас есть {n} элемента"
}

profile_updated(gender) {
    [male]   => "Он обновил свой профиль"
    [female] => "Она обновила свой профиль"
    [other]  => "Они обновили свой профиль"
}

ui.menu {
    home = "Главная"
    about = "О нас"
    contact = "Контакты"
    settings = "Параметры"
}

order_total = "Итого: {price} (включая налог)"
```

---

## 3. Проверьте и отформатируйте переводы

```bash
# Проверьте синтаксические ошибки
mbel lint locales/

# Автоматически отформатируйте
mbel fmt locales/

# Покажите статистику
mbel stats locales/
```

---

## 4. Скомпилируйте переводы

```bash
# Скомпилируйте в один файл JSON
mbel compile locales/ -o dist/translations.json

# Включите карту источника для отладки
mbel compile locales/ -o dist/translations.json -sourcemap
```

---

## 5. Создайте приложение Go

### Базовый пример (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. Инициализируйте Manager
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "ru",
		FallbackChain: []string{"ru", "en"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Простой поиск
	fmt.Println(m.Get("ru", "app_name", nil))

	// 3. Интерполяция переменных
	vars := mbel.Vars{"name": "Виктор"}
	greeting := m.Get("ru", "greeting", vars)
	fmt.Println(greeting)

	// 4. Множественное число
	for _, count := range []int{1, 2, 5} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("ru", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. Глобальное API
	mbel.Init(m)
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "Мария"}))
}
```

---

## 6. Добавьте тесты

```bash
go test -v ./...
```

---

## 7. Соберите и запустите

```bash
go build -o hello-mbel ./cmd
./hello-mbel
```

---

## 8. Развертывание в производстве

### Вариант A: Встройте JSON

```go
//go:embed dist/translations.json
var translationsJSON []byte
```

### Вариант B: Разверните с Docker

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN go build -o /tmp/hello-mbel ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/hello-mbel .
COPY dist/translations.json ./dist/
CMD ["./hello-mbel"]
```

---

## 9. Следующие шаги

1. **[Руководство](Manual.md)** — Полная документация
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — Глубокий технический анализ
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — Расширение MBEL
4. **[Лучшие практики безопасности](SECURITY.md)** — Предотвращение XSS

---

## Устранение неполадок

| Проблема | Решение |
|----------|---------|
| `no such file or directory: locales/` | `mkdir -p locales` |
| `Undefined: mbel` | Проверьте `go.mod` |
| `Syntax error at line 5` | `mbel lint locales/` |
| Перевод не найден | Проверьте ключ, код языка и цепь резервного варианта |
