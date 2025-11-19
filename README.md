# Telegram Gitlab Bot

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

_Установка для MAC_:

```bash
brew install go-task
```

[Установка для других платформ](https://taskfile.dev/docs/installation).

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
