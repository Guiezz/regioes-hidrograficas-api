package model

import "time"

// Section representa uma seção de texto do plano (Identificação, Infraestrutura, etc)
type Section struct {
	ID uint `json:"id" gorm:"primaryKey"`

	// Referência opcional à bacia (caso expanda para outras bacias no futuro)
	BasinID *uint `json:"basin_id"`

	Number  string `json:"number"`                   // Ex: "1", "1.1", "3.4.1"
	Title   string `json:"title"`                    // Ex: "Infraestrutura Hídrica"
	Content string `json:"content" gorm:"type:text"` // O texto longo
	Level   int    `json:"level"`                    // 1, 2, 3... (Nível hierárquico para o frontend indentar)

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Section) TableName() string {
	return "secoes_plano"
}
