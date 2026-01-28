package main

import (
	"fmt"
	"log"

	"github.com/guiezz/regioes-hidrograficas-api/config"
	"github.com/guiezz/regioes-hidrograficas-api/db"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
)

var basePath = "dados_importacao/curu"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}
	database := db.Init(cfg)

	fmt.Println("💥 DESTRUINDO TABELAS ANTIGAS (LIMPEZA TOTAL)...")

	// 1. Dropa tudo para garantir que não sobrem tabelas com nomes errados (ex: actions vs acoes)
	// A ordem importa: apagar tabelas filhas (que têm FK) antes das pais
	err = database.Migrator().DropTable(
		&model.TypologyStats{},
		&model.ConsolidatedStats{},
		&model.Measurement{}, // nome tabela: medicoes
		&model.Action{},      // nome tabela: acoes
		&model.Program{},     // nome tabela: programas
		&model.Axis{},        // nome tabela: eixos
		&model.Section{},     // nome tabela: sections (ou secoes se vc mudou)
		&model.Basin{},       // nome tabela: bacias
	)
	if err != nil {
		log.Printf("⚠️ (Info) Drop Table: %v", err)
	}

	fmt.Println("🏗️ Recriando Schema do Banco...")
	// 2. Recria as tabelas com os nomes forçados nas Models
	err = database.AutoMigrate(
		&model.Section{},
		&model.Basin{},
		&model.Axis{},
		&model.Program{},
		&model.Action{},
		&model.Measurement{},
		&model.ConsolidatedStats{},
		&model.TypologyStats{},
	)
	if err != nil {
		log.Fatalf("❌ Erro no AutoMigrate: %v", err)
	}

	// 3. Garante a Bacia Curu (ID 1)
	var bacia model.Basin
	database.FirstOrCreate(&bacia, model.Basin{Name: "Curu"})
	log.Printf("🌊 Bacia Principal: %s (ID: %d)", bacia.Name, bacia.ID)

	// 4. Executa os importadores
	seedSections(database, bacia)
	seedMonitoring(database, bacia)
}
