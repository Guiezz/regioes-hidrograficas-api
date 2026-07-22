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

func seedKPIs(db *gorm.DB, basin model.Basin, folderPath string) {
	fmt.Printf("📊 [KPIs] Importando para %s...\n", basin.Name)

	fullPath := filepath.Join(folderPath, "kpis.json")

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("⚠️ Arquivo de KPIs não encontrado em: %s", fullPath)
		return
	}

	fileContent, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Printf("❌ Erro ao ler arquivo: %v", err)
		return
	}

	var kpis []model.BasinKPI
	if err := json.Unmarshal(fileContent, &kpis); err != nil {
		log.Printf("❌ Erro no JSON de KPIs: %v", err)
		return
	}

	count := 0
	for _, k := range kpis {
		k.BasinID = basin.ID

		var exists model.BasinKPI
		if err := db.
			Where("basin_id = ? AND tab = ? AND view_mode = ? AND \"order\" = ?",
				basin.ID, k.Tab, k.ViewMode, k.Order).
			First(&exists).Error; err == nil {
			exists.Value = k.Value
			exists.Unit = k.Unit
			exists.Label = k.Label
			exists.Sublabel = k.Sublabel
			exists.Icon = k.Icon
			exists.Severity = k.Severity
			db.Save(&exists)
		} else {
			db.Create(&k)
			count++
		}
	}
	fmt.Printf("✅ [KPIs] %d indicadores importados para bacia %s.\n", count, basin.Name)
}
