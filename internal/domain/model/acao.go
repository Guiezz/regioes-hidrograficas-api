package model

import "time"

type Action struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	ProgramID uint    `json:"program_id"`
	Program   Program `json:"program" gorm:"foreignKey:ProgramID"`

	// Identificação
	ReservatorioNome string `json:"reservatorio_nome"`            // Coluna A
	Description      string `json:"description" gorm:"type:text"` // Coluna D
	Typology         string `json:"typology"`                     // Coluna E
	Source           string `json:"source"`                       // Coluna F

	// Financeiro
	TotalBudget float64 `json:"total_budget"`
	BudgetUnit  string  `json:"budget_unit"`

	// Cronograma
	StartYear int `json:"start_year"`
	EndYear   int `json:"end_year"`

	// Responsáveis (Opcional, mas bom ter)
	ResponsavelPrincipal string `json:"responsavel_principal"`
	OrgaosEnvolvidos     string `json:"orgaos_envolvidos"`

	// --- NOVOS CAMPOS (O erro estava aqui: eles faltavam) ---
	ExecutionPerc float64 `json:"execution_perc"` // Métrica (0.0 a 1.0)
	PDPWeight     int     `json:"pdp_weight"`     // Peso Dificuldade (1-5)
	IEA           float64 `json:"iea"`            // Índice Ajustado

	// Relacionamentos
	Measurements []Measurement `json:"measurements" gorm:"foreignKey:ActionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Measurement struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	ActionID uint `json:"action_id"`

	ReferenceMonth string `json:"reference_month"`

	// Indicadores
	ExecutionPerc float64 `json:"execution_perc"`
	PDPWeight     int     `json:"pdp_weight"`
	IEA           float64 `json:"iea"`
	IEARelativo   float64 `json:"iea_relativo"`

	MeasuredAt time.Time `json:"measured_at"`
}

func (Action) TableName() string      { return "acoes" }
func (Measurement) TableName() string { return "medicoes" }
