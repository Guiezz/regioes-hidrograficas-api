package main

import "strings"

// hasMetric informa se a célula "Métrica de Indicador de execução" da planilha
// está preenchida. Só então uma medição deve ser criada.
func hasMetric(cell string) bool {
	return strings.TrimSpace(cell) != ""
}
