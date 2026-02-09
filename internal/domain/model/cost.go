package model

import "time"

// Cost represent os dados do arquivo "Custos.xlsx"
type Cost struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	BasinID    uint      `json:"basin_id"`
	Eixo       string    `json:"eixo"`
	ValorTotal string    `json:"valor_total"` // Mantido como string para preservar o formato "R$ ..."
	Percentual float64   `json:"percentual"`
	P2021_2025 string    `json:"p2021_2025"`
	P2025_2030 string    `json:"p2025_2030"`
	P2030_2035 string    `json:"p2030_2035"`
	P2035_2040 string    `json:"p2035_2040"`
	P2040_2045 string    `json:"p2040_2045"`
	P2045_2050 string    `json:"p2045_2050"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Cost) TableName() string {
	return "costs"
}
