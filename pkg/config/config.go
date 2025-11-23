package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sensetion/tgGitlabBot/internal/domain"
	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig   `mapstructure:"server"`
	GitLab       GitLabConfig   `mapstructure:"gitlab"`
	Telegram     TelegramConfig `mapstructure:"telegram"`
	LogLevel     string         `mapstructure:"log_level"`
	Repositories []domain.Repository
}

type ServerConfig struct {
	Port              int           `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	Shutdown          time.Duration `mapstructure:"shutdown_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	CompressSize      int           `mapstructure:"compress_size"`
}

type GitLabConfig struct {
	WebhookSecret string `mapstructure:"webhook_secret"`
}

type TelegramConfig struct {
	BotToken   string        `mapstructure:"bot_token"`
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"max_retries"`
}

func Load() (*Config, error) {
	// 1. Загружаем .env файлы (godotenv)
	if err := loadEnvFiles(); err != nil {
		log.Printf("Warning: %v", err)
	}

	v := viper.New()
	// Имя файла без расширения (будет искать config.yaml, config.yml, config.json и т.д.)
	v.SetConfigName("config")
	// Тип конфигурационного файла (yaml, json, toml и т.д.)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config") // Относительный путь от корня проекта
	v.AddConfigPath(".")        // Текущая директория (fallback)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// v.AddConfigPath("/etc/tgGitlabBot/")  // Системная директория (для production)
	// v.AddConfigPath("$HOME/.tgGitlabBot") // Home directory пользователя (опционально)
	// Префикс для env переменных (TGBOT_SERVER_PORT, TGBOT_LOG_LEVEL и т.д.)
	// v.SetEnvPrefix("TGBOT")

	// Явная привязка вложенных полей к env переменным
	// Формат: BindEnv("yaml.путь.к.полю", "ENV_VARIABLE_NAME")
	// v.BindEnv("server.port", "SERVER_PORT")

	// Устанавливаем дефолтные значения, если не найдены в файле или env
	// v.SetDefault("server.port", 8080)

	if err := v.ReadInConfig(); err != nil {
		// Если файл не найден - это критическая ошибка
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := loadRepositories(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load repositories: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	log.Printf("Configuration loaded successfully: %d repositories", len(cfg.Repositories))

	return &cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}

	if c.GitLab.WebhookSecret == "" {
		return fmt.Errorf("gitlab webhook secret is required")
	}

	if c.Telegram.BotToken == "" {
		return fmt.Errorf("telegram bot token is required")
	}

	if len(c.Repositories) == 0 {
		return fmt.Errorf("at least one repository must be configured")
	}

	return nil
}

// loadEnvFiles загружает .env файлы с приоритетами
func loadEnvFiles() error {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Загружает в порядке приоритета (последний = высший)
	files := []string{
		".env",                      // Базовые значения
		fmt.Sprintf(".env.%s", env), // Environment-specific
		".env.local",                // Локальные (высший приоритет)
	}

	// Фильтруем только существующие файлы
	existingFiles := []string{}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			existingFiles = append(existingFiles, file)
		}
	}

	if len(existingFiles) == 0 {
		log.Println("No .env files found")
		return nil
	}

	// ✅ Загружаем все файлы сразу (порядок важен!)
	if err := godotenv.Load(existingFiles...); err != nil {
		return fmt.Errorf("failed to load env files: %w", err)
	}

	log.Printf("✓ Loaded env files: %v", existingFiles)
	return nil
}

func loadRepositories(cfg *Config) error {
	// Пути поиска repositories.json (в порядке приоритета)
	paths := []string{
		"config/repositories.json",
		"repositories.json",
		"/etc/tgGitlabBot/repositories.json",
	}

	var data []byte
	var err error
	var usedPath string

	// Ищем файл в разных локациях
	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			usedPath = path
			break
		}
	}
	if err != nil {
		return fmt.Errorf("repositories.json not found in any location: %w", err)
	}

	log.Printf("Loading repositories from: %s", usedPath)

	var repoConfig struct {
		Repositories []domain.Repository `json:"repositories"`
	}
	if err := json.Unmarshal(data, &repoConfig); err != nil {
		return fmt.Errorf("failed to parse repositories.json: %w", err)
	}

	cfg.Repositories = repoConfig.Repositories
	if len(cfg.Repositories) == 0 {
		return fmt.Errorf("repositories.json is empty")
	}

	return nil
}

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
