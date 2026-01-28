package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseFloatBlindado converte qualquer string suja do Excel em float64 seguro
// Ex: "R$ 1.200,50" -> 1200.50 | "" -> 0.0 | "#VALUE!" -> 0.0
func ParseFloatBlindado(value string) float64 {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\"", "") // Remove aspas extras do CSV

	// Trata erros de fórmula ou vazios
	if value == "" || value == "-" || strings.HasPrefix(value, "#") {
		return 0.0
	}

	// Remove "R$", "%" e letras
	re := regexp.MustCompile(`[^\d,\.-]`) // Mantém digitos, vírgula, ponto e sinal negativo
	value = re.ReplaceAllString(value, "")

	// Padroniza separador decimal (Brasil ',' para Internacional '.')
	// Cuidado: Se tiver ponto de milhar (1.000,00), remove o ponto primeiro
	if strings.Contains(value, ",") && strings.Contains(value, ".") {
		value = strings.ReplaceAll(value, ".", "") // Remove ponto de milhar
	}
	value = strings.ReplaceAll(value, ",", ".")

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return f
}

// ParseIntBlindado usa a lógica float e converte para int
func ParseIntBlindado(value string) int {
	f := ParseFloatBlindado(value)
	return int(f)
}

// ParsePeriod extrai "2024" e "2033" de "2024 - 2033"
func ParsePeriod(value string) (int, int) {
	parts := strings.Split(value, "-")
	if len(parts) < 2 {
		return 0, 0
	}
	s := ParseIntBlindado(parts[0])
	e := ParseIntBlindado(parts[1])
	return s, e
}
