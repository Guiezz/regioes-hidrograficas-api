package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

type KPIHandler struct {
	db *gorm.DB
}

func NewKPIHandler(db *gorm.DB) *KPIHandler {
	return &KPIHandler{db: db}
}

// GetKPIs godoc
// @Summary      KPIs da Situação Hídrica
// @Description  Retorna indicadores-chave por aba (infraestrutura, oferta, demanda, balanco) e modo de visão (atual, futuro)
// @Tags         KPIs
// @Produce      json
// @Param        basin_id   query   int     false  "ID da Bacia (Padrão: 1)"
// @Success      200  {object}  map[string]map[string][]model.BasinKPI
// @Router       /kpis [get]
func (h *KPIHandler) GetKPIs(c *gin.Context) {
	basinID := c.Query("basin_id")
	if basinID == "" {
		basinID = "1"
	}

	var kpis []model.BasinKPI
	if err := h.db.
		Where("basin_id = ?", basinID).
		Order("\"tab\", \"view_mode\", \"order\"").
		Find(&kpis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar KPIs"})
		return
	}

	result := make(map[string]map[string][]model.BasinKPI)
	for _, k := range kpis {
		if result[k.Tab] == nil {
			result[k.Tab] = make(map[string][]model.BasinKPI)
		}
		result[k.Tab][k.ViewMode] = append(result[k.Tab][k.ViewMode], k)
	}

	c.JSON(http.StatusOK, result)
}
