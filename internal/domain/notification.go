package domain

// Notification представляет уведомление для отправки в Telegram
type Notification struct {
	ChatID     string
	Message    string
	ParseMode  string // "Markdown" или "HTML"
	RetryCount int
}
