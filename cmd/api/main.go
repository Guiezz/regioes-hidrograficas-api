package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guiezz/regioes-hidrograficas-api/config"
	"github.com/guiezz/regioes-hidrograficas-api/db"
	"github.com/guiezz/regioes-hidrograficas-api/internal/api/routes"
	// IMPORTANTE: Isso será gerado automaticamente depois, mas precisamos importar
	// _ "github.com/guiezz/regioes-hidrograficas-api/docs"
)

// @title           Regiões Hidrográficas API
// @version         1.0
// @description     API para monitoramento, gestão e visualização dos Planos de Recursos Hídricos (Curu, Salgado, etc).
// @termsOfService  http://swagger.io/terms/

// @contact.name    Suporte Técnico
// @contact.email   suporte@exemplo.com.br

// @host            localhost:8080
// @BasePath        /api/v1

func main() {

	files, err := os.ReadDir("./assets")
	if err != nil {
		log.Println("❌ ERRO: O Go não conseguiu ler a pasta ./assets:", err)
	} else {
		log.Println("📂 CONTEÚDO DA PASTA ASSETS VISTO PELO GO:")
		for _, file := range files {
			log.Println("   📄 Encontrado:", file.Name())
		}
	}
	// 1. Configuração
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Erro ao carregar config: %v", err)
	}

	// 2. Banco de Dados
	dbConnection := db.Init(cfg)

	// 3. Servidor Web (Gin)
	r := gin.Default()

	// 4. Configurar CORS (Para o Frontend acessar)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://api-hidrografica.onrender.com",
			"https://regioes-hidrograficas-front.vercel.app",
		}, AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 5. Registrar Rotas
	routes.RegisterRoutes(r, dbConnection)

	// 6. Rodar
	log.Printf("🚀 Servidor rodando na porta %s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
