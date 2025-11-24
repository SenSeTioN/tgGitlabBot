package domain

type Repository struct {
	ID             string   `json:"id" mapstructure:"id"`
	TelegramChatID string   `json:"telegram_channel_id" mapstructure:"telegram_channel_id"`
	Branches       []string `json:"branches" mapstructure:"branches"`
	Enabled        bool     `json:"enabled" mapstructure:"enabled"`
}

// HasBranch проверяет, нужно ли мониторить данную ветку
func (r *Repository) HasBranch(branch string) bool {
	// Если не указаны конкретные ветки, мониторим все
	if len(r.Branches) == 0 {
		return true
	}

	// Проверяем, есть ли ветка в списке
	for _, b := range r.Branches {
		if b == branch {
			return true
		}
	}

	return false
}

// IsEnabled проверяет, активен ли репозиторий
func (r *Repository) IsEnabled() bool {
	return r.Enabled
}
