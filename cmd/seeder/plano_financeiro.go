package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	var rawData model.PlanoAcaoResponse
	if err := json.Unmarshal(fileContent, &rawData); err != nil {
		log.Printf("❌ Erro JSON custos: %v", err)
		return
	}

	// 1. Salvar o Resumo Geral (Upsert baseado no BasinID)
	resumo := rawData.ResumoGeral
	resumo.BasinID = basinID
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "basin_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"valor_total_previsto", "percentual"}),
	}).Create(&resumo)

	// 2. Processar Eixos
	for _, eixoJson := range rawData.PlanoAcao {
		eixo := model.EixoAcao{
			BasinID:             basinID,
			Eixo:                eixoJson.Eixo,
			ValorTotalProjetado: eixoJson.ValorTotalProjetado,
			Percentual:          eixoJson.Percentual,
			Periodos:            eixoJson.Periodos,
		}

		// Importante: habilitar FullSaveAssociations para salvar Periodos e CustosVariaveis
		err := db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&eixo).Error
		if err != nil {
			log.Printf("❌ Erro ao salvar eixo %s: %v", eixo.Eixo, err)
		}
	}
	fmt.Printf("✅ [Financeiro] Eixos e custos variáveis processados.\n")
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
		Solicitacao      string `json:"Solicitacao"`
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
