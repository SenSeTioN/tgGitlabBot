package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	chihttp "github.com/sensetion/tgGitlabBot/internal/adapters/http"
	"github.com/sensetion/tgGitlabBot/internal/config"
	"github.com/sensetion/tgGitlabBot/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	//TODO: delete in release
	logger.PrettyStructurePrint("📋 Loaded configuration:", cfg)

	r := chihttp.Init(cfg)
	port := strconv.Itoa(cfg.Server.Port)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", port),
		Handler:           r,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", port)
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Shutdown)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Graceful shutdown не удался: %v. Принудительное закрытие...\n", err)

		// Принудительно закрываем соединения
		if closeErr := server.Close(); closeErr != nil {
			log.Printf("❌ Ошибка принудительного закрытия: %v\n", closeErr)
		}
	} else {
		log.Println("✅ Сервер остановлен корректно")
	}
}
