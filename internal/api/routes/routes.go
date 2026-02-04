package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/internal/api/handler"
	"gorm.io/gorm"

	// Swagger imports
	_ "github.com/guiezz/regioes-hidrograficas-api/docs" // <--- IMPORTANTE: Importa a pasta docs que será criada
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// ... (criação dos handlers igual) ...
	dashboardHandler := handler.NewDashboardHandler(db)
	contentHandler := handler.NewContentHandler(db)
	actionHandler := handler.NewActionHandler(db)

	// Rota da Documentação
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/assets", "./assets")

	api := r.Group("/api/v1")
	{
		api.GET("/content", contentHandler.GetSections)
		api.GET("/dashboard/radar", dashboardHandler.GetRadarData)
		api.GET("/dashboard/consolidated", dashboardHandler.GetConsolidated)
		api.GET("/actions", actionHandler.GetActions)
		api.GET("/actions/filters", actionHandler.GetFilters)
	}
}
