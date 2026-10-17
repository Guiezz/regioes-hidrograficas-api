package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

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

	importadas, ignoradas := 0, 0
	// Delete + inserts na mesma transação: reimportação idempotente e atômica.
	// Em Postgres um erro de linha aborta a transação, então ele também faz rollback.
	err = db.Transaction(func(tx *gorm.DB) error {
		importadas, ignoradas = 0, 0
		if err := tx.Where("basin_id = ?", basinID).Delete(&model.ActionMatrix{}).Error; err != nil {
			return err
		}
		for _, item := range rawData {
			if strings.TrimSpace(item.AcoesEspecificas) == "" {
				ignoradas++
				continue
			}
			m := model.ActionMatrix{
				BasinID:            basinID,
				TipoMatriz:         normalizeEixoMatriz(item.Matriz),
				TipoMatrizOriginal: item.Matriz,
				SolicitacoesCBH:    item.Solicitacao,
				AcoesEspecificas:   item.AcoesEspecificas,
				Programa:           item.Programa,
				Instituicoes:       item.Instituicoes,
				Prioridade:         normalizePrioridade(item.Prioridade),
			}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
			importadas++
		}
		return nil
	})
	if err != nil {
		log.Printf("❌ Erro ao importar matriz (nada foi alterado para a bacia %d): %v", basinID, err)
		return
	}
	fmt.Printf("✅ [Financeiro] %d de %d linhas da matriz importadas (%d sem ação específica ignoradas)\n", importadas, len(rawData), ignoradas)
}
