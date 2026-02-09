package model

import "time"

type ActionMatrix struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	BasinID          uint      `json:"basin_id"`
	TipoMatriz       string    `json:"tipo_matriz"`
	SolicitacoesCBH  string    `json:"solicitacoes_cbh" gorm:"type:text"`
	AcoesEspecificas string    `json:"acoes_especificas" gorm:"type:text"`
	Programa         string    `json:"programa"`
	Instituicoes     string    `json:"instituicoes_envolvidas" gorm:"type:text"`
	Prioridade       string    `json:"prioridade"`
	CreatedAt        time.Time `json:"created_at"`
}

func (ActionMatrix) TableName() string { return "matriz_acoes_prioridade" }
