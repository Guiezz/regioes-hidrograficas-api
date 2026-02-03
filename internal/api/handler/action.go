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
// @Description  Retorna lista de ações com diversos filtros dinâmicos (Eixo, Programa, Ano, etc).
// @Tags         Actions
// @Produce      json
// @Param        basin_id   query   int     false  "ID da Bacia"
// @Param        ano        query   string  false  "Ano de vigência (ex: 2024)"
// @Param        eixo       query   string  false  "Filtro por nome do Eixo"
// @Param        programa   query   string  false  "Filtro por nome do Programa"
// @Param        tipologia  query   string  false  "Filtro por Tipologia"
// @Param        search     query   string  false  "Busca textual na descrição"
// @Success      200        {object} map[string]interface{}
// @Router       /actions [get]
func (h *ActionHandler) GetActions(c *gin.Context) {
	basinID := c.Query("basin_id")

	// Filtros opcionais vindos do frontend
	filterEixo := c.Query("eixo")
	filterPrograma := c.Query("programa")
	filterTypology := c.Query("tipologia")
	filterAno := c.Query("ano") // Filtra se a ação está vigente neste ano
	search := c.Query("search") // Busca textual

	// Prepara a Query base carregando os relacionamentos necessários
	// Preload carrega os dados das tabelas filhas (Measurements, Program, Axis)
	query := h.db.Model(&model.Action{}).
		Preload("Program.Axis"). // Carrega Programa e Eixo
		Preload("Measurements")  // Carrega o histórico de medições

	// 1. Filtro Obrigatório de Bacia (via join com Program e Axis se necessário,
	// mas como salvamos ReservatorioNome na Ação, podemos usar ele ou o basin_id das tabelas pai)
	// Como seu modelo Action não tem basin_id direto, filtramos pelo nome ou join.
	// O jeito mais seguro no seu modelo atual é via JOINs:
	query = query.Joins("JOIN programas ON programas.id = acoes.program_id").
		Joins("JOIN eixos ON eixos.id = programas.axis_id").
		Joins("JOIN bacias ON bacias.id = eixos.basin_id")

	if basinID != "" {
		query = query.Where("bacias.id = ?", basinID)
	}

	// 2. Filtros Dinâmicos
	if filterEixo != "" {
		query = query.Where("eixos.name ILIKE ?", "%"+filterEixo+"%")
	}
	if filterPrograma != "" {
		query = query.Where("programas.name ILIKE ?", "%"+filterPrograma+"%")
	}
	if filterTypology != "" {
		query = query.Where("acoes.typology ILIKE ?", "%"+filterTypology+"%")
	}
	if search != "" {
		query = query.Where("acoes.description ILIKE ?", "%"+search+"%")
	}
	if filterAno != "" {
		// Filtra ações ativas naquele ano (Start <= Ano <= End)
		query = query.Where("acoes.start_year <= ? AND acoes.end_year >= ?", filterAno, filterAno)
	}

	var actions []model.Action
	result := query.Find(&actions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ações"})
		return
	}

	// Retorna lista filtrada
	c.JSON(http.StatusOK, gin.H{
		"count":   len(actions),
		"filters": c.Request.URL.Query(),
		"data":    actions,
	})
}

// GetFilters godoc
// @Summary      Opções de Filtro
// @Description  Retorna listas únicas de Eixos, Programas e Tipologias para preencher combos do frontend.
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

	// CORREÇÃO: Inicializar como array vazio []string{} e não var nil
	eixos := []string{}
	h.db.Model(&model.Axis{}).Where("basin_id = ?", basinID).Pluck("name", &eixos)

	programas := []string{}
	h.db.Table("programas").
		Joins("JOIN eixos ON eixos.id = programas.axis_id").
		Where("eixos.basin_id = ?", basinID).
		Pluck("programas.name", &programas)

	tipologias := []string{}
	// Precisamos fazer um join complexo para filtrar tipologias SÓ desta bacia
	h.db.Table("acoes").
		Joins("JOIN programas ON programas.id = acoes.program_id").
		Joins("JOIN eixos ON eixos.id = programas.axis_id").
		Where("eixos.basin_id = ?", basinID).
		Distinct("typology").
		Pluck("typology", &tipologias)

	c.JSON(http.StatusOK, gin.H{
		"eixos":      eixos,
		"programas":  programas,
		"tipologias": tipologias,
	})
}
