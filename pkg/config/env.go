package config

import (
	"os"
	"strconv"
	"time"
)

// ========== Дополнительные helper методы (опционально) ==========

func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

/*
Функция парсит формат времени из переменных окружения например: 1h30m - 1 час 30 минут
  - “ns” - наносекунды
  - “us” (или “µs”) - микросекунды
  - “ms” - миллисекунды
  - “s” - секунды
  - “m” - минуты
  - “h” - часы
*/
func getDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return d
}
