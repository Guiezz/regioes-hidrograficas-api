package model

import "time"

// Tabela por Tipologia (Obra, Estudo, etc)
type TypologyStats struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	BasinID        uint   `json:"basin_id"` // <--- NOVO: Vincula à bacia (Salgado)
	ReferenceMonth string `json:"reference_month"`

	Category     string  `json:"category"`
	Count        int     `json:"count"`
	SumPotential float64 `json:"sum_potential"`
	Percentage   float64 `json:"percentage"`

	CreatedAt time.Time `json:"created_at"`
}

func (TypologyStats) TableName() string {
	return "indicadores_tipologia"
}
