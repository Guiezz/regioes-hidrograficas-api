package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
)

func seedPlanoFinanceiro(db *gorm.DB, basin model.Basin, folderPath string) {
	fmt.Printf("📈 [Financeiro] Importando custos e matriz para: %s\n", basin.Name)

	importarCustosJSON(db, basin.ID, folderPath)
	importarMatrizJSON(db, basin.ID, folderPath)
}

func importarCustosJSON(db *gorm.DB, basinID uint, folderPath string) {
	fullPath := filepath.Join(folderPath, "custos.json")
	fileContent, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Printf("⚠️ Arquivo custos.json não encontrado em %s", fullPath)
		return
	}

	type CostJSON struct {
		Eixo       string  `json:"Eixo"`
		ValorTotal string  `json:"Valor Total"`
		Percentual float64 `json:"Percentual"`
		P2021_2025 string  `json:"2021-2025"`
		P2025_2030 string  `json:"2025-2030"`
		P2030_2035 string  `json:"2030-2035"`
		P2035_2040 string  `json:"2035-2040"`
		P2040_2045 string  `json:"2040-2045"`
		P2045_2050 string  `json:"2045-2050"`
	}

	var rawData []CostJSON
	if err := json.Unmarshal(fileContent, &rawData); err != nil {
		log.Printf("❌ Erro JSON custos: %v", err)
		return
	}

	for _, item := range rawData {
		c := model.Cost{
			BasinID:    basinID,
			Eixo:       item.Eixo,
			ValorTotal: item.ValorTotal,
			Percentual: item.Percentual,
			P2021_2025: item.P2021_2025,
			P2025_2030: item.P2025_2030, // Adicionado
			P2030_2035: item.P2030_2035, // Adicionado
			P2035_2040: item.P2035_2040, // Adicionado
			P2040_2045: item.P2040_2045, // Adicionado
			P2045_2050: item.P2045_2050,
		}
		db.Where("eixo = ? AND basin_id = ?", c.Eixo, basinID).FirstOrCreate(&c)
	}
	fmt.Printf("✅ [Financeiro] %d custos processados.\n", len(rawData))
}

func importarMatrizJSON(db *gorm.DB, basinID uint, folderPath string) {
	fullPath := filepath.Join(folderPath, "matriz_acao.json")
	fileContent, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Printf("⚠️ Arquivo matriz_acao.json não encontrado em %s", fullPath)
		return
	}

	type MatrizJSON struct {
		Matriz           string `json:"Matriz"`
		Solicitacao      string `json:"Solicitacao"` // Atenção ao acento no JSON se houver
		AcoesEspecificas string `json:"Ações específicas"`
		Programa         string `json:"Programa"`
		Instituicoes     string `json:"Instituições envolvidas"`
		Prioridade       string `json:"Prioridade"`
	}

	var rawData []MatrizJSON
	if err := json.Unmarshal(fileContent, &rawData); err != nil {
		log.Printf("❌ Erro JSON matriz: %v", err)
		return
	}

	for _, item := range rawData {
		m := model.ActionMatrix{
			BasinID:          basinID,
			TipoMatriz:       item.Matriz,
			SolicitacoesCBH:  item.Solicitacao,
			AcoesEspecificas: item.AcoesEspecificas,
			Programa:         item.Programa,
			Instituicoes:     item.Instituicoes,
			Prioridade:       item.Prioridade,
		}
		db.Where("programa = ? AND acoes_especificas = ? AND basin_id = ?", m.Programa, m.AcoesEspecificas, basinID).FirstOrCreate(&m)
	}
	fmt.Printf("✅ [Financeiro] %d ações da matriz processadas.\n", len(rawData))
}
