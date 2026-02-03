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
// @Success      200  {object}  map[string]float64
// @Router       /dashboard/radar [get]
func (h *DashboardHandler) GetRadarData(c *gin.Context) {
	basinID := c.Query("basin_id")
	if basinID == "" {
		basinID = "1"
	}

	var stats []model.TypologyStats

	// Busca os dados da bacia
	result := h.db.Where("basin_id = ?", basinID).Find(&stats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados"})
		return
	}

	// MUDANÇA AQUI:
	// Criamos um mapa simples: Chave (Nome da Categoria) -> Valor (Potencial)
	// Isso garante que cada tipologia seja única e o Frontend consiga mapear facilmente.
	response := make(map[string]float64)

	for _, s := range stats {
		// Se houver duplicatas no banco, isso pega o último valor encontrado,
		// garantindo que não tenhamos "Obra" duas vezes no gráfico.
		response[s.Category] = s.SumPotential
	}

	// O JSON padrão do Go ordena as chaves do mapa alfabeticamente (A-Z),
	// o que é ótimo para manter o desenho do Radar consistente.
	c.JSON(http.StatusOK, response)
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
