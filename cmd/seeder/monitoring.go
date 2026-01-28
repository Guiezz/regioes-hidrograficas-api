package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/guiezz/regioes-hidrograficas-api/internal/domain/model"
	"github.com/guiezz/regioes-hidrograficas-api/pkg/importer/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

var desiredCategories = []string{
	"Ampliação de pessoal", "Estrutural", "Estudo", "Projeto", "Preservação",
	"Gestão", "Legislação", "Capacitação", "Planejamento", "Fiscalização",
	"Conservação de água", "Regulamentação", "Monitoramento", "Programa",
	"Articulação institucional", "Comunicação",
}

type statsTemp struct {
	Count        int
	SumPotential float64
}

func seedMonitoring(db *gorm.DB, basin model.Basin) {
	fmt.Printf("📊 Iniciando importação de Monitoramento para %s (ID %d)...\n", basin.Name, basin.ID)

	excelPath := findFirstExcel(basePath)
	if excelPath == "" {
		log.Println("⚠️ Nenhum arquivo Excel encontrado.")
		return
	}

	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		log.Printf("❌ Erro ao abrir Excel: %v", err)
		return
	}
	defer f.Close()

	rows, _ := f.GetRows(f.GetSheetList()[0])
	tx := db.Begin()

	typologyMap := make(map[string]*statsTemp)
	for _, cat := range desiredCategories {
		typologyMap[cat] = &statsTemp{Count: 0, SumPotential: 0.0}
	}

	var totalPotencialGlobal float64 = 0.0
	var totalRealizadoGlobal float64 = 0.0
	var countActionsGlobal int = 0

	for i, row := range rows {
		if i == 0 || len(row) < 8 {
			continue
		}

		// Leitura
		colEixo := safeGet(row, 1)
		colPrograma := safeGet(row, 2)
		colAcao := safeGet(row, 3)
		colTipologia := safeGet(row, 4)
		colFonte := safeGet(row, 5)
		colOrcamento := safeGet(row, 6)
		colCronograma := safeGet(row, 7)
		colMetrica := parseFloat(safeGet(row, 8))

		// Tratamento
		tipologiaNormalizada := normalizeTipologia(colTipologia)
		startYear, endYear := parseCleanYear(colCronograma)
		orcamentoVal, orcamentoUnit := cleanMoney(colOrcamento)
		peso := utils.CalcularPeso(colAcao, colTipologia)

		// Cálculos IEA
		ieaRealizado := float64(peso) * colMetrica
		ieaPotencial := float64(peso) * 1.0

		// Hierarquia
		var eixo model.Axis
		tx.FirstOrCreate(&eixo, model.Axis{Name: colEixo, BasinID: basin.ID})
		var prog model.Program
		tx.FirstOrCreate(&prog, model.Program{Name: colPrograma, AxisID: eixo.ID})

		// Salvar Ação (COM OS NOVOS CAMPOS PREENCHIDOS)
		var acao model.Action
		err := tx.Where("description = ? AND program_id = ?", colAcao, prog.ID).First(&acao).Error

		if err == gorm.ErrRecordNotFound {
			acao = model.Action{
				ProgramID:        prog.ID,
				ReservatorioNome: basin.Name,
				Description:      colAcao,
				Typology:         tipologiaNormalizada,
				Source:           colFonte,
				TotalBudget:      orcamentoVal,
				BudgetUnit:       orcamentoUnit,
				StartYear:        startYear,
				EndYear:          endYear,

				// --- CAMPOS NOVOS ---
				ExecutionPerc: colMetrica,
				PDPWeight:     peso,
				IEA:           ieaRealizado,
			}
			tx.Create(&acao)
		} else {
			// Se já existe (num re-seed parcial), atualiza os valores
			acao.ExecutionPerc = colMetrica
			acao.PDPWeight = peso
			acao.IEA = ieaRealizado
			tx.Save(&acao)
		}

		// Acumuladores e Stats (mantidos iguais)
		if endYear <= 2033 && endYear > 0 {
			totalPotencialGlobal += ieaPotencial
			totalRealizadoGlobal += ieaRealizado
			countActionsGlobal++

			if stats, ok := typologyMap[tipologiaNormalizada]; ok {
				stats.Count++
				stats.SumPotential += ieaPotencial
			} else {
				typologyMap[tipologiaNormalizada] = &statsTemp{Count: 1, SumPotential: ieaPotencial}
			}
		}

		// Medição (Histórico)
		m := &model.Measurement{
			ActionID:       acao.ID,
			ReferenceMonth: time.Now().Format("01/2006"),
			ExecutionPerc:  colMetrica,
			PDPWeight:      peso,
			IEA:            ieaRealizado,
			MeasuredAt:     time.Now(),
		}
		tx.Create(m)
	}

	// Relatórios e Salvamento Global
	fmt.Printf("\n📋 RELATÓRIO FINAL (%s)\n", basin.Name)
	printed := make(map[string]bool)
	for _, cat := range desiredCategories {
		stats := typologyMap[cat]
		saveStat(tx, basin.ID, cat, stats, totalPotencialGlobal)
		printed[cat] = true
	}
	for cat, stats := range typologyMap {
		if !printed[cat] && stats.Count > 0 {
			saveStat(tx, basin.ID, cat, stats, totalPotencialGlobal)
		}
	}

	indiceGlobal := 0.0
	if totalPotencialGlobal > 0 {
		indiceGlobal = totalRealizadoGlobal / totalPotencialGlobal
	}

	tx.Create(&model.ConsolidatedStats{
		BasinID:        basin.ID,
		ReferenceMonth: time.Now().Format("01/2006"),
		IEA_Max:        totalPotencialGlobal,
		IndiceAtual:    totalRealizadoGlobal,
		IndiceGlobal:   indiceGlobal,
		TotalActions:   countActionsGlobal,
	})

	tx.Commit()
	fmt.Println("✅ Monitoramento importado com sucesso!")
}

// --- Helpers ---

func saveStat(tx *gorm.DB, basinID uint, cat string, stats *statsTemp, totalGlobal float64) {
	perc := 0.0
	if totalGlobal > 0 {
		perc = (stats.SumPotential / totalGlobal) * 100
	}
	tx.Create(&model.TypologyStats{
		BasinID:        basinID,
		ReferenceMonth: time.Now().Format("01/2006"),
		Category:       cat,
		Count:          stats.Count,
		SumPotential:   stats.SumPotential,
		Percentage:     perc,
	})
	fmt.Printf("%-30s | %-5d | %-15.1f | %-10.2f%%\n", cat, stats.Count, stats.SumPotential, perc)
}

func normalizeTipologia(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 {
		return "Não Identificado"
	}
	lowerRaw := strings.ToLower(raw)
	for _, cat := range desiredCategories {
		if strings.ToLower(cat) == lowerRaw {
			return cat
		}
	}
	return strings.ToUpper(raw[:1]) + raw[1:]
}

func extractYear(val string) int {
	nums := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, val)
	if len(nums) > 4 {
		nums = nums[:4]
	}
	if len(nums) == 4 {
		y, _ := strconv.Atoi(nums)
		return y
	}
	return 0
}

func parseCleanYear(val string) (int, int) {
	if val == "" {
		return 0, 0
	}
	val = strings.ReplaceAll(val, "–", "-")
	val = strings.ReplaceAll(val, "—", "-")
	parts := strings.Split(val, "-")
	start := extractYear(parts[0])
	end := start
	if len(parts) > 1 {
		end = extractYear(parts[1])
	}
	return start, end
}

func findFirstExcel(dir string) string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".xlsx") && !strings.HasPrefix(f.Name(), "~$") {
			return filepath.Join(dir, f.Name())
		}
	}
	return ""
}

func safeGet(row []string, index int) string {
	if index < len(row) {
		return strings.TrimSpace(row[index])
	}
	return ""
}

func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func cleanMoney(raw string) (float64, string) {
	valStr := strings.ReplaceAll(raw, "R$", "")
	valStr = strings.TrimSpace(valStr)
	parts := strings.Split(valStr, " ")
	numPart := parts[0]
	numPart = strings.ReplaceAll(numPart, ".", "")
	numPart = strings.ReplaceAll(numPart, ",", ".")
	val, _ := strconv.ParseFloat(numPart, 64)
	unit := "Global"
	if strings.Contains(raw, "/") {
		split := strings.Split(raw, "/")
		if len(split) > 1 {
			unit = strings.TrimSpace(split[1])
		}
	}
	return val, unit
}
