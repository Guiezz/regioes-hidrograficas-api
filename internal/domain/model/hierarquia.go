package model

import "time"

// Basin representa a Bacia Hidrográfica (ex: "Salgado", "Curu")
type Basin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Basin) TableName() string {
	return "bacias"
}

// Axis representa os Eixos (ex: "Oferta Hídrica", "Demanda", "Político-Institucional")
type Axis struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"not null"` // Ex: "Eixo Oferta Hídrica"
	BasinID  uint      `json:"basin_id"`             // Pertence a uma bacia
	Basin    Basin     `json:"-" gorm:"foreignKey:BasinID"`
	Programs []Program `json:"programs,omitempty"`
}

func (Axis) TableName() string {
	return "eixos"
}

// Program representa os Programas (ex: "Incremento da oferta hídrica superficial")
type Program struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	AxisID      uint   `json:"axis_id"`
	Axis        Axis   `json:"-" gorm:"foreignKey:AxisID"`
	Description string `json:"description"`
}

func (Program) TableName() string {
	return "programas"
}
