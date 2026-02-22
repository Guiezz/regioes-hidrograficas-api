package db

import (
	"fmt"
	"log"

	"github.com/guiezz/regioes-hidrograficas-api/config"
	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model" // Importante: importar seus modelos
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg config.Config) *gorm.DB {
	var dsn string

	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Fortaleza",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Falha ao conectar no banco de dados: %v", err)
	}

	// --- ADICIONE ESTE BLOCO AQUI ---
	// O AutoMigrate verifica as structs e cria/atualiza as tabelas no Postgres automaticamente
	err = db.AutoMigrate(
		&model.Basin{},         // Suas tabelas existentes
		&model.ActionMatrix{},  // Suas tabelas existentes
		&model.ResumoGeral{},   // Nova tabela de resumo financeiro
		&model.EixoAcao{},      // Nova tabela de eixos (substitui model.Cost antigo)
		&model.PeriodoAcao{},   // Nova tabela de períodos (2021-2025, etc)
		&model.CustoVariavel{}, // Nova tabela de métricas unitárias
	)
	if err != nil {
		log.Fatalf("❌ Falha ao realizar migração automática: %v", err)
	}
	// -------------------------------

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Falha ao obter instância SQL: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("❌ Banco de dados não responde ao Ping: %v", err)
	}

	log.Println("✅ Conexão e Migrações estabelecidas com sucesso!")
	return db
}
