package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

type ActionHandler struct {
	db *gorm.DB
}

func NewActionHandler(db *gorm.DB) *ActionHandler {
	return &ActionHandler{db: db}
}

// GetActions godoc
// @Summary      Listar Ações (Matriz/Planos)
// @Description  Retorna lista de ações com filtros ajustados para a estrutura atual do banco.
// @Tags         Actions
// @Produce      json
// @Param        basin_id   query   int     false  "ID da Bacia"
// @Param        eixo       query   string  false  "Filtro por nome do Eixo"
// @Param        tipologia  query   string  false  "Filtro por Tipologia"
// @Success      200        {object} map[string]interface{}
// @Router       /actions [get]
func (h *ActionHandler) GetActions(c *gin.Context) {
	basinID := c.Query("basin_id")

	// Filtros
	filterEixo := c.Query("eixo")
	filterPrograma := c.Query("programa")
	filterTypology := c.Query("tipologia")
	filterAno := c.Query("ano")
	search := c.Query("search")

	// Query Base
	query := h.db.Model(&model.Action{}).
		Select("acoes.*"). // Garante que a struct Action seja a principal
		Preload("Program.Axis").
		Preload("Measurements")

	// JOINS
	query = query.Joins("JOIN programas ON programas.id = acoes.program_id").
		Joins("JOIN eixos ON eixos.id = programas.axis_id").
		Joins("JOIN bacias ON bacias.id = eixos.basin_id")

	// 1. Filtro de Bacia
	if basinID != "" {
		query = query.Where("bacias.id = ?", basinID)
	}

	// 2. Filtros Dinâmicos

	// CORREÇÃO CRÍTICA:
	// O filtro de EIXO deve olhar para a tabela PROGRAMAS, pois é lá que "Demanda Hídrica" está salvo.
	if filterEixo != "" && filterEixo != "todos" {
		query = query.Where("programas.name ILIKE ?", "%"+filterEixo+"%")
	}

	// O filtro de PROGRAMA também olha para PROGRAMAS (ou seria redundante, mas mantemos para compatibilidade)
	if filterPrograma != "" && filterPrograma != "todos" {
		query = query.Where("programas.name ILIKE ?", "%"+filterPrograma+"%")
	}

	if filterTypology != "" && filterTypology != "todos" {
		query = query.Where("acoes.typology ILIKE ?", "%"+filterTypology+"%")
	}
	if search != "" {
		query = query.Where("acoes.description ILIKE ?", "%"+search+"%")
	}
	if filterAno != "" {
		query = query.Where("acoes.start_year <= ? AND acoes.end_year >= ?", filterAno, filterAno)
	}

	var actions []model.Action
	if err := query.Order("programas.name ASC, acoes.id ASC").Find(&actions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ações"})
		return
	}

	// --- CORREÇÃO DE MAPEAMENTO (POPULATE) ---
	// Aqui sobrescrevemos o AxisName para o Frontend receber o dado correto
	for i := range actions {
		// Cenário Atual:
		// Program.Name = "DEMANDA HÍDRICA" (O que queremos exibir como Eixo)
		// Axis.Name    = "Curu" (Nome da Bacia, não queremos exibir isso no card de Eixo)

		if actions[i].Program.Name != "" {
			// Joga o nome do Programa para o campo AxisName
			actions[i].AxisName = actions[i].Program.Name
		} else {
			// Fallback
			actions[i].AxisName = "Geral"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(actions),
		"data":  actions,
	})
}

// GetFilters godoc
// @Summary      Opções de Filtro
// @Description  Retorna listas corrigidas onde 'eixos' busca dados da tabela 'programas'.
// @Tags         Actions
// @Produce      json
// @Param        basin_id   query   int     false  "ID da Bacia"
// @Success      200        {object} map[string]interface{}
// @Router       /actions/filters [get]
func (h *ActionHandler) GetFilters(c *gin.Context) {
	basinID := c.Query("basin_id")
	if basinID == "" {
		basinID = "1"
	}

	// 1. Busca EIXOS (Corrigido: Busca na tabela 'programas')
	// Isso vai retornar ["DEMANDA HÍDRICA", "OFERTA HÍDRICA", ...]
	eixos := []string{}
	h.db.Table("programas").
		Joins("JOIN eixos ON eixos.id = programas.axis_id").
		Where("eixos.basin_id = ?", basinID).
		Distinct("programas.name").
		Order("programas.name ASC").
		Pluck("programas.name", &eixos)

	// 2. Busca PROGRAMAS
	// Como os níveis estão deslocados, o "Programa" é igual ao Eixo neste contexto.
	// Mantemos a mesma lista para o filtro não ficar vazio.
	programas := eixos

	// 3. Busca TIPOLOGIAS (Normal)
	tipologias := []string{}
	h.db.Table("acoes").
		Joins("JOIN programas ON programas.id = acoes.program_id").
		Joins("JOIN eixos ON eixos.id = programas.axis_id").
		Where("eixos.basin_id = ?", basinID).
		Distinct("typology").
		Order("typology ASC").
		Pluck("typology", &tipologias)

	c.JSON(http.StatusOK, gin.H{
		"eixos":      eixos,
		"programas":  programas,
		"tipologias": tipologias,
	})
}
