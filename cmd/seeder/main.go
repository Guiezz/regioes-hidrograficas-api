package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/guiezz/regioes-hidrograficas-api/config"
	"github.com/guiezz/regioes-hidrograficas-api/db"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
)

// Lista de bacias para importar (o nome deve bater com a pasta em dados_importacao/)
var targetBasins = []string{
	"Alto Jaguaribe",
	"Baixo Jaguaribe",
	"Banabuiú",
	"Crateús",
	"Coreau",
	"Curu",
	"Ibiapaba",
	"Litoral",
	"Médio Jaguaribe",
	"Metropolitana",
	"Salgado",
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}
	database := db.Init(cfg)

	fmt.Println("💥 DESTRUINDO TABELAS ANTIGAS (LIMPEZA TOTAL)...")
	// Dropa tabelas para garantir integridade
	err = database.Migrator().DropTable(
		&model.TypologyStats{},
		&model.ConsolidatedStats{},
		&model.Measurement{},
		&model.Action{},
		&model.Program{},
		&model.Axis{},
		&model.Section{},
		&model.Basin{},
		&model.ActionMatrix{},
		&model.ResumoGeral{},   // Novo
		&model.EixoAcao{},      // Substitui o antigo model.Cost
		&model.PeriodoAcao{},   // Novo
		&model.CustoVariavel{}, // Novo
	)
	if err != nil {
		log.Printf("⚠️ (Info) Drop Table: %v", err)
	}

	fmt.Println("🏗️ Recriando Schema do Banco...")
	err = database.AutoMigrate(
		&model.Section{},
		&model.Basin{},
		&model.Axis{},
		&model.Program{},
		&model.Action{},
		&model.Measurement{},
		&model.ConsolidatedStats{},
		&model.TypologyStats{},
		&model.ActionMatrix{},
		&model.ResumoGeral{},   // Novo
		&model.EixoAcao{},      // Substitui o antigo model.Cost
		&model.PeriodoAcao{},   // Novo
		&model.CustoVariavel{}, // Novo
	)
	if err != nil {
		log.Fatalf("❌ Erro no AutoMigrate: %v", err)
	}

	// Loop para processar cada bacia
	for _, basinName := range targetBasins {
		fmt.Printf("\n========================================\n")
		fmt.Printf("🌊 PROCESSANDO BACIA: %s\n", strings.ToUpper(basinName))
		fmt.Printf("========================================\n")

		// 1. Cria ou Garante a Bacia no Banco
		var bacia model.Basin
		// FirstOrCreate busca pelo nome, se não achar, cria.
		if err := database.FirstOrCreate(&bacia, model.Basin{Name: basinName}).Error; err != nil {
			log.Printf("❌ Erro ao criar bacia %s: %v", basinName, err)
			continue
		}

		// Define o caminho da pasta: ex: dados_importacao/curu ou dados_importacao/salgado
		// Importante: As pastas devem estar em minúsculo no sistema de arquivos
		folderPath := filepath.Join("dados_importacao", strings.ToLower(basinName))

		// 2. Executa os importadores passando o caminho específico
		seedSections(database, bacia, folderPath)
		seedMonitoring(database, bacia, folderPath)
		seedPlanoFinanceiro(database, bacia, folderPath)
	}

	fmt.Println("\n🚀 SEED COMPLETO PARA TODAS AS BACIAS!")
}
