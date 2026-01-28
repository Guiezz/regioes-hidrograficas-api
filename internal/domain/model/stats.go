package model

import "time"

// Tabela consolidada (Indices Globais)
type ConsolidatedStats struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	BasinID        uint   `json:"basin_id"` // <--- NOVO: Vincula à bacia (Salgado)
	ReferenceMonth string `json:"reference_month"`

	IEA_Max      float64 `json:"iea_max"`
	IEA_Min      float64 `json:"iea_min"`
	IndiceAtual  float64 `json:"indice_atual"`
	IndiceGlobal float64 `json:"indice_global"`
	TotalActions int     `json:"total_actions"`

	CreatedAt time.Time `json:"created_at"`
}

func (ConsolidatedStats) TableName() string {
	return "indicadores_consolidados"
}
