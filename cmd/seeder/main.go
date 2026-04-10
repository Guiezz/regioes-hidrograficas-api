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

type BasinConfig struct {
	Name   string
	Folder string
}

var targetBasins = []BasinConfig{
	{"Alto Jaguaribe", "alto jaguaribe"},
	{"Baixo Jaguaribe", "baixo jaguaribe"},
	{"Banabuiú", "banabuiú"},
	{"Crateús", "crateús"},
	{"Coreaú", "coreau"},
	{"Curu", "curu"},
	{"Ibiapaba", "ibiapaba"},
	{"Litoral", "litoral"},
	{"Médio Jaguaribe", "médio jaguaribe"},
	{"Metropolitana", "metropolitana"},
	{"Salgado", "salgado"},
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}
	database := db.Init(cfg)

	fmt.Println("💥 DESTRUINDO TABELAS ANTIGAS (LIMPEZA TOTAL)...")
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
		&model.ResumoGeral{},
		&model.EixoAcao{},
		&model.PeriodoAcao{},
		&model.CustoVariavel{},
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
		&model.ResumoGeral{},
		&model.EixoAcao{},
		&model.PeriodoAcao{},
		&model.CustoVariavel{},
	)
	if err != nil {
		log.Fatalf("❌ Erro no AutoMigrate: %v", err)
	}

	for _, basin := range targetBasins {
		fmt.Printf("\n========================================\n")
		fmt.Printf("🌊 PROCESSANDO BACIA: %s\n", strings.ToUpper(basin.Name))
		fmt.Printf("========================================\n")

		var bacia model.Basin
		if err := database.FirstOrCreate(&bacia, model.Basin{Name: basin.Name}).Error; err != nil {
			log.Printf("❌ Erro ao criar bacia %s: %v", basin.Name, err)
			continue
		}

		folderPath := filepath.Join("dados_importacao", basin.Folder)

		seedSections(database, bacia, folderPath)
		seedMonitoring(database, bacia, folderPath)
		seedPlanoFinanceiro(database, bacia, folderPath)
	}
	fmt.Println("\n🚀 SEED COMPLETO PARA TODAS AS BACIAS!")
}
