package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sensetion/tgGitlabBot/internal/controller/http/response"
)

type HealthHandler struct {
	telegramHealthCheck func() error
	startTime           time.Time
}

func NewHealthHandler(telegramHealthCheck func() error) *HealthHandler {
	return &HealthHandler{
		telegramHealthCheck: telegramHealthCheck,
		startTime:           time.Now(),
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Ready - readiness check с проверкой всех зависимостей
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	allHealthy := true

	// 1. Проверка файловой системы (может ли сервер писать файлы)
	if err := h.checkFileSystem(); err != nil {
		checks["filesystem"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		checks["filesystem"] = "healthy"
	}

	// 2. Проверка горутин (утечка горутин = проблема)
	goroutines := runtime.NumGoroutine()
	if goroutines > 1000 { // Пороговое значение
		checks["goroutines"] = fmt.Sprintf("warning: %d goroutines", goroutines)
		allHealthy = false
	} else {
		checks["goroutines"] = fmt.Sprintf("healthy: %d goroutines", goroutines)
	}

	// 3. Проверка памяти
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryMB := m.Alloc / 1024 / 1024
	if memoryMB > 512 { // Больше 512MB - подозрительно
		checks["memory"] = fmt.Sprintf("warning: %dMB allocated", memoryMB)
	} else {
		checks["memory"] = fmt.Sprintf("healthy: %dMB allocated", memoryMB)
	}

	// 4. Проверка Telegram API (если передан checker)
	if h.telegramHealthCheck != nil {
		if err := h.telegramHealthCheck(); err != nil {
			checks["telegram"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			checks["telegram"] = "healthy"
		}
	}

	// 5. Uptime сервера
	uptime := time.Since(h.startTime)
	checks["uptime"] = uptime.String()

	// Формируем response
	status := "ready"
	statusCode := http.StatusOK

	if !allHealthy {
		status = "not ready"
		statusCode = http.StatusServiceUnavailable
	}

	response.JSON(w, statusCode, map[string]any{
		"status": status,
		"checks": checks,
	})
}

// checkFileSystem проверяет доступность файловой системы
func (h *HealthHandler) checkFileSystem() error {
	// Создаем временный файл для проверки записи
	tmpFile := filepath.Join(os.TempDir(), ".health-check")

	// Пытаемся записать файл
	if err := os.WriteFile(tmpFile, []byte("health-check"), 0644); err != nil {
		return fmt.Errorf("cannot write to filesystem: %w", err)
	}

	// Пытаемся прочитать файл
	if _, err := os.ReadFile(tmpFile); err != nil {
		return fmt.Errorf("cannot read from filesystem: %w", err)
	}

	// Удаляем временный файл
	if err := os.Remove(tmpFile); err != nil {
		return fmt.Errorf("cannot delete from filesystem: %w", err)
	}

	return nil
}
