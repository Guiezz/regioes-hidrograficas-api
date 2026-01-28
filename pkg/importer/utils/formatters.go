package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// CleanMoneyCuru separa o valor da unidade.
// Ex: "6500,00 R$/ha" -> retorna (6500.00, "R$/ha")
// Ex: "R$ 350.000,00" -> retorna (350000.00, "Global")
func CleanMoneyCuru(raw string) (float64, string) {
	if raw == "" {
		return 0.0, ""
	}

	// 1. Tenta extrair apenas números e vírgula para o valor
	// Regex: pega o primeiro grupo de dígitos com pontuação
	reNum := regexp.MustCompile(`[0-9]+[.,]?[0-9]*[.,]?[0-9]*`)
	numStr := reNum.FindString(raw) // Pega "6500,00" ou "350.000,00"

	// Limpa pontos de milhar e troca vírgula decimal por ponto
	cleanNum := strings.ReplaceAll(numStr, ".", "")   // Tira ponto de milhar
	cleanNum = strings.ReplaceAll(cleanNum, ",", ".") // Troca vírgula por ponto

	val, _ := strconv.ParseFloat(cleanNum, 64)

	// 2. Define a unidade
	unit := "Global"
	if strings.Contains(raw, "/") {
		// Pega o que sobra ou padrões comuns
		parts := strings.Split(raw, " ")
		for _, p := range parts {
			if strings.Contains(p, "/") {
				unit = p
				break
			}
		}
	}

	return val, unit
}

// SplitCuruAxis separa "EIXO DEMANDA HÍDRICA - Programa de..."
func SplitCuruAxis(fullString string) (string, string) {
	// O padrão do Curu parece usar " - " ou " – " (hífen ou travessão)
	normalized := strings.ReplaceAll(fullString, "–", "-")
	parts := strings.SplitN(normalized, "-", 2)

	eixo := strings.TrimSpace(parts[0])
	programa := "Geral"

	if len(parts) > 1 {
		programa = strings.TrimSpace(parts[1])
	}

	return eixo, programa
}
