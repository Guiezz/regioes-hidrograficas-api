package model

type PlanoAcaoResponse struct {
	ResumoGeral ResumoGeral `json:"resumoGeral"`
	PlanoAcao   []EixoAcao  `json:"planoAcao"`
}

type ResumoGeral struct {
	ID                 uint    `gorm:"primaryKey"`
	BasinID            uint    `json:"-" gorm:"uniqueIndex"` // Um resumo por bacia
	ValorTotalPrevisto float64 `json:"valorTotalPrevisto"`
	Percentual         float64 `json:"percentual"`
}

type EixoAcao struct {
	ID                  uint          `gorm:"primaryKey"`
	BasinID             uint          `json:"-" gorm:"index"`
	Eixo                string        `json:"eixo" gorm:"column:nome_eixo"` // Nome da coluna no banco
	ValorTotalProjetado float64       `json:"valorTotalProjetado"`
	Percentual          float64       `json:"percentual"`
	Periodos            []PeriodoAcao `json:"periodos" gorm:"foreignKey:EixoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type PeriodoAcao struct {
	ID              uint            `gorm:"primaryKey"`
	EixoID          uint            `json:"-" gorm:"index"`
	Intervalo       string          `json:"intervalo"`
	CustoFixo       float64         `json:"custoFixo"`
	CustosVariaveis []CustoVariavel `json:"custosVariaveis" gorm:"foreignKey:PeriodoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CustoVariavel struct {
	ID            uint    `gorm:"primaryKey"`
	PeriodoID     uint    `json:"-" gorm:"index"`
	Descricao     string  `json:"descricao"`
	ValorUnitario float64 `json:"valorUnitario"`
}
