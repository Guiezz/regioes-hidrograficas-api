package utils

import (
	"strings"
)

// Mapa de Pesos (Regra de Negócio)
var pesosPorCategoria = map[string]int{
	"OBRA":    5,
	"SERVIÇO": 4,
	"ESTUDO":  3, "PROJETO": 3, "PROGRAMA": 3,
	"FISCALIZAÇÃO": 2, "CADASTRO": 2, "MONITORAMENTO": 2, "AQUISIÇÃO": 2, "AMPLIAÇÃO": 2,
	"CAPACITAÇÃO": 1, "GESTÃO": 4, "COMUNICAÇÃO": 1, "LEGISLAÇÃO": 3,
}

var classificador = map[string]string{
	// PESO 5
	"estrutural": "OBRA", "obra": "OBRA", "construção": "OBRA",
	"implantação": "OBRA", "barragem": "OBRA", "adutora": "OBRA",
	"sistema adutor": "OBRA", "tamponamento": "OBRA", "perfuracao": "OBRA",

	// PESO 4
	"manutenção": "SERVIÇO", "conservação": "SERVIÇO", "revitalização": "SERVIÇO",
	"redução de perdas": "SERVIÇO", "redução das perdas": "SERVIÇO",
	"preservação": "SERVIÇO", "limpeza": "SERVIÇO",
	"comissão gestora": "SERVIÇO", // Gestão de alta complexidade

	// PESO 3
	"estudo": "ESTUDO", "diagnóstico": "ESTUDO", "plano": "ESTUDO",
	"planejamento": "ESTUDO", "subterrânea": "ESTUDO",
	"projeto": "PROJETO", "mapeamento": "PROJETO", "banco de dados": "PROJETO",
	"sistema": "PROJETO", "programa": "PROGRAMA",

	// PESO 2
	"fiscalização": "FISCALIZAÇÃO", "ordenação": "FISCALIZAÇÃO",
	"monitoramento": "MONITORAMENTO", "medição": "MONITORAMENTO",
	"ampliação de pessoal": "AMPLIAÇÃO", // <--- O que você pediu
	"aquisição":            "AQUISIÇÃO", "equipamento": "AQUISIÇÃO",

	// PESO 1
	"capacitação": "CAPACITAÇÃO", "educação": "CAPACITAÇÃO",
	"comunicação": "COMUNICAÇÃO", "divulgação": "COMUNICAÇÃO", "campanha": "COMUNICAÇÃO",
	"gestão": "GESTÃO", "articulação": "GESTÃO", "apoio": "GESTÃO", "contratação": "GESTÃO",
	"Legislação": "LEGISLAÇÃO", "norma": "LEGISLAÇÃO", "regulamentação": "LEGISLAÇÃO",
}

func CalcularPeso(textoAcao, tipologia string) int {
	texto := strings.ToLower(textoAcao + " " + tipologia)
	tipologiaLower := strings.ToLower(tipologia)

	maiorPeso := 1

	// 1. Busca Peso por Palavra-Chave
	for palavra, categoria := range classificador {
		if strings.Contains(texto, palavra) {
			p := pesosPorCategoria[categoria]
			if p > maiorPeso {
				maiorPeso = p
			}
		}
	}

	// 2. Trava de Segurança da Tipologia (para não superestimar Estudos)
	// Se a tipologia for explicitamente "Estudo" ou "Projeto", o peso máximo é 3.
	isTipologiaTravada := strings.Contains(tipologiaLower, "estudo") ||
		strings.Contains(tipologiaLower, "projeto") ||
		strings.Contains(tipologiaLower, "diagnóstico")

	isExcecao := strings.Contains(texto, "manutenção") || strings.Contains(texto, "redução") // Manutenção vale 4

	if isTipologiaTravada && !isExcecao {
		if maiorPeso > 3 {
			return 3
		}
	}

	return maiorPeso
}
