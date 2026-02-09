package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

type BasinHandler struct {
	db *gorm.DB
}

func NewBasinHandler(db *gorm.DB) *BasinHandler {
	return &BasinHandler{db: db}
}

// GetBasins lista todas as bacias disponíveis
func (h *BasinHandler) GetBasins(c *gin.Context) {
	var basins []model.Basin

	// Busca todas as bacias ordenadas por ID
	if err := h.db.Order("id asc").Find(&basins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar bacias"})
		return
	}

	c.JSON(http.StatusOK, basins)
}

// GetBasinByID retorna os detalhes de uma bacia específica
func (h *BasinHandler) GetBasinByID(c *gin.Context) {
	id := c.Param("id")
	var basin model.Basin

	if err := h.db.First(&basin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bacia não encontrada"})
		return
	}

	c.JSON(http.StatusOK, basin)
}
