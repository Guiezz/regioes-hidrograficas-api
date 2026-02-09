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
// @Description  Retorna uma lista de custos, permitindo filtrar por bacia específica via ID.
// @Tags         Financeiro
// @Accept       json
// @Produce      json
// @Param        basin_id   query     int  false  "ID da Bacia Hidrográfica (ex: 1 para Curu)"
// @Success      200  {array}   model.Cost
// @Failure      500  {object}  map[string]string "Erro interno ao buscar dados"
// @Router       /financeiro/custos [get]
func (h *FinanceiroHandler) GetCustos(c *gin.Context) {
	var custos []model.Cost
	query := h.DB

	// Filtro por BasinID
	if basinID := c.Query("basin_id"); basinID != "" {
		query = query.Where("basin_id = ?", basinID)
	}

	if err := query.Find(&custos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar custos no banco de dados"})
		return
	}
	c.JSON(http.StatusOK, custos)
}

// GetMatriz godoc
// @Summary      Listar matriz de ações e prioridades
// @Description  Retorna a matriz de ações filtrada por bacia. Útil para separar dados do Curu, Salgado, etc.
// @Tags         Financeiro
// @Accept       json
// @Produce      json
// @Param        basin_id   query     int  false  "ID da Bacia Hidrográfica"
// @Success      200  {array}   model.ActionMatrix
// @Failure      500  {object}  map[string]string "Erro interno ao buscar dados"
// @Router       /financeiro/matriz [get]
func (h *FinanceiroHandler) GetMatriz(c *gin.Context) {
	var matriz []model.ActionMatrix
	query := h.DB

	// Filtro por BasinID
	if basinID := c.Query("basin_id"); basinID != "" {
		query = query.Where("basin_id = ?", basinID)
	}

	if err := query.Find(&matriz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar matriz de ações"})
		return
	}
	c.JSON(http.StatusOK, matriz)
}
