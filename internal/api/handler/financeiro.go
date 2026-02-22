package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

type FinanceiroHandler struct {
	DB *gorm.DB
}

func NewFinanceiroHandler(db *gorm.DB) *FinanceiroHandler {
	return &FinanceiroHandler{DB: db}
}

// GetCustos godoc
// @Summary      Listar custos de planejamento
// @Description  Retorna o resumo geral e os custos por eixo de uma bacia específica via ID.
// @Tags         Financeiro
// @Accept       json
// @Produce      json
// @Param        basin_id   query    int  false  "ID da Bacia Hidrográfica"
// @Success      200  {object}   model.PlanoAcaoResponse
// @Failure      500  {object}  map[string]string "Erro interno ao buscar dados"
// @Router       /financeiro/custos [get]
func (h *FinanceiroHandler) GetCustos(c *gin.Context) {
	basinID := c.Query("basin_id")
	if basinID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O parâmetro basin_id é obrigatório"})
		return
	}

	var response model.PlanoAcaoResponse

	// 1. Buscar o Resumo Geral da bacia
	if err := h.DB.Where("basin_id = ?", basinID).First(&response.ResumoGeral).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar resumo financeiro"})
			return
		}
	}

	// 2. Buscar os Eixos com carregamento recursivo (Preload) de Períodos e Custos Variáveis
	if err := h.DB.Preload("Periodos.CustosVariaveis").
		Where("basin_id = ?", basinID).
		Find(&response.PlanoAcao).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar eixos de ação"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMatriz godoc
// @Summary      Listar matriz de ações e prioridades
// @Description  Retorna a matriz de ações filtrada por bacia.
// @Tags         Financeiro
// @Accept       json
// @Produce      json
// @Param        basin_id   query    int  false  "ID da Bacia Hidrográfica"
// @Success      200  {array}   model.ActionMatrix
// @Failure      500  {object}  map[string]string "Erro interno ao buscar dados"
// @Router       /financeiro/matriz [get]
func (h *FinanceiroHandler) GetMatriz(c *gin.Context) {
	var matriz []model.ActionMatrix
	query := h.DB

	if basinID := c.Query("basin_id"); basinID != "" {
		query = query.Where("basin_id = ?", basinID)
	}

	if err := query.Find(&matriz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar matriz de ações"})
		return
	}
	c.JSON(http.StatusOK, matriz)
}
