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

// Agora recebe folderPath para saber onde buscar o arquivo
func seedSections(db *gorm.DB, basin model.Basin, folderPath string) {
	fmt.Printf("📖 [Textos] Importando para %s...\n", basin.Name)

	fullPath := filepath.Join(folderPath, "textos_plano.json")

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("⚠️ Arquivo de textos não encontrado em: %s", fullPath)
		return
	}

	fileContent, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Printf("❌ Erro ao ler arquivo: %v", err)
		return
	}

	var sections []model.Section
	if err := json.Unmarshal(fileContent, &sections); err != nil {
		log.Printf("❌ Erro no JSON de textos: %v", err)
		return
	}

	count := 0
	for _, s := range sections {
		// Vincula à bacia correta
		s.BasinID = &basin.ID

		var exists model.Section
		// Verifica duplicidade baseada no número E na bacia
		if err := db.Where("number = ? AND basin_id = ?", s.Number, basin.ID).First(&exists).Error; err == nil {
			exists.Title = s.Title
			exists.Content = s.Content
			exists.Level = s.Level
			exists.Image = s.Image
			db.Save(&exists)
		} else {
			db.Create(&s)
			count++
		}
	}
	fmt.Printf("✅ [Textos] %d seções novas importadas para bacia %s.\n", count, basin.Name)
}
