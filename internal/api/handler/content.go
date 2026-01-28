package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

type ContentHandler struct {
	db *gorm.DB
}

func NewContentHandler(db *gorm.DB) *ContentHandler {
	return &ContentHandler{db: db}
}

// GetSections godoc
// @Summary      Textos do Plano
// @Description  Retorna os textos hierárquicos (1.1, 1.2...) para montar a página de leitura
// @Tags         Content
// @Produce      json
// @Param        basin_id   query   int     false  "ID da Bacia"
// @Success      200        {array} model.Section
// @Router       /content [get]
func (h *ContentHandler) GetSections(c *gin.Context) {
	basinID := c.Query("basin_id") // Opcional, se tiver textos diferentes por bacia

	var sections []model.Section
	query := h.db.Order("id ASC") // A ordem de inserção do JSON geralmente é a correta

	if basinID != "" {
		query = query.Where("basin_id = ?", basinID)
	}

	query.Find(&sections)

	c.JSON(http.StatusOK, sections)
}
