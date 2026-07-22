package model

import "time"

type BasinKPI struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	BasinID   uint      `json:"basin_id" gorm:"index"`
	Tab       string    `json:"tab"`        // infraestrutura | oferta | demanda | balanco
	ViewMode  string    `json:"view_mode"`  // atual | futuro
	Value     string    `json:"value"`
	Unit      string    `json:"unit"`
	Label     string    `json:"label"`
	Sublabel  string    `json:"sublabel"`
	Icon      string    `json:"icon"`
	Severity  string    `json:"severity"` // positive | warning | critical
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BasinKPI) TableName() string {
	return "kpis_bacia"
}
