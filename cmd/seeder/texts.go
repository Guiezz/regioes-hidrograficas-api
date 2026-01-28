package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

// Agora recebe "basin" como argumento
func seedSections(db *gorm.DB, basin model.Basin) {
	fmt.Println("📖 [Textos] Iniciando importação...")

	fullPath := filepath.Join("dados_importacao/Salgado", "textos_plano.json") // Ajuste se mudou a pasta

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Tenta na pasta do Curu se não achar no Salgado
		fullPath = filepath.Join("dados_importacao/curu", "textos_plano.json")
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			log.Printf("⚠️ Arquivo de textos não encontrado.")
			return
		}
	}

	fileContent, _ := ioutil.ReadFile(fullPath)
	var sections []model.Section
	json.Unmarshal(fileContent, &sections)

	for _, s := range sections {
		// Vincula à bacia correta
		s.BasinID = &basin.ID

		var exists model.Section
		if err := db.Where("number = ? AND basin_id = ?", s.Number, basin.ID).First(&exists).Error; err == nil {
			exists.Title = s.Title
			exists.Content = s.Content
			exists.Level = s.Level
			db.Save(&exists)
		} else {
			db.Create(&s)
		}
	}
	fmt.Printf("✅ [Textos] Sucesso! %d seções vinculadas à bacia %s (ID %d).\n", len(sections), basin.Name, basin.ID)
}
