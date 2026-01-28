package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetRadarData godoc
// @Summary      Dados do Gráfico Radar
// @Description  Retorna soma dos pesos agrupados por tipologia (Obra, Estudo, etc)
// @Tags         Dashboard
// @Produce      json
// @Param        basin_id   query      int  false  "ID da Bacia (Padrão: 1)"
// @Success      200  {object}  map[string]interface{}
// @Router       /dashboard/radar [get]
func (h *DashboardHandler) GetRadarData(c *gin.Context) {
	basinID := c.Query("basin_id")
	if basinID == "" {
		// Se não passar ID, tenta pegar o primeiro do banco ou erro
		basinID = "1"
	}

	var stats []model.TypologyStats

	// Busca as estatísticas mais recentes para a bacia
	// Ordenamos por categoria para garantir que o gráfico sempre desenhe igual
	result := h.db.Where("basin_id = ?", basinID).
		Order("category ASC").
		Find(&stats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados"})
		return
	}

	// Retorna JSON otimizado para gráficos (Arrays de labels e valores)
	var labels []string
	var values []float64

	for _, s := range stats {
		labels = append(labels, s.Category)
		values = append(values, s.SumPotential) // Ou s.Percentage se preferir o gráfico em %
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": labels, // Ex: ["OBRA", "ESTUDO"...]
		"data":   values, // Ex: [20.0, 15.0...]
		"raw":    stats,  // Manda os objetos completos caso o front queira detalhes
	})
}

// GetConsolidated godoc
// @Summary      Dados Consolidados (Cards)
// @Description  Retorna IEA Máximo, IEA Atual e Índice Global
// @Tags         Dashboard
// @Produce      json
// @Param        basin_id   query      int  false  "ID da Bacia (Padrão: 1)"
// @Success      200  {object}  model.ConsolidatedStats
// @Router       /dashboard/consolidated [get]
func (h *DashboardHandler) GetConsolidated(c *gin.Context) {
	basinID := c.Query("basin_id")
	if basinID == "" {
		basinID = "1"
	}

	var stats model.ConsolidatedStats
	// Pega o último consolidado gerado
	h.db.Where("basin_id = ?", basinID).Last(&stats)

	c.JSON(http.StatusOK, stats)
}
