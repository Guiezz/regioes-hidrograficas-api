package db

import (
	"fmt"
	"log"

	"github.com/guiezz/regioes-hidrograficas-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Fortaleza",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Falha ao conectar no banco de dados: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Falha ao obter instância SQL: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("❌ Banco de dados não responde ao Ping: %v", err)
	}

	log.Println("✅ Conexão com Banco de Dados estabelecida com sucesso!")
	return db
}
