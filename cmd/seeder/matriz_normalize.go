package main

import "strings"

var semAcento = strings.NewReplacer(
	"á", "a", "à", "a", "â", "a", "ã", "a",
	"é", "e", "ê", "e",
	"í", "i",
	"ó", "o", "ô", "o", "õ", "o",
	"ú", "u",
	"ç", "c",
)

// chave devolve o texto em minúsculas, sem acentos e sem espaços nas pontas.
func chave(raw string) string {
	return semAcento.Replace(strings.ToLower(strings.TrimSpace(raw)))
}

// normalizeEixoMatriz converte o título livre da matriz no nome canônico do eixo.
func normalizeEixoMatriz(raw string) string {
	k := chave(raw)
	switch {
	case strings.Contains(k, "demanda"):
		return "Eixo Demanda Hídrica"
	case strings.Contains(k, "oferta"):
		return "Eixo Oferta Hídrica"
	case strings.Contains(k, "gerenciamento"):
		return "Eixo Gerenciamento das Águas"
	case strings.Contains(k, "politico"):
		return "Eixo Político-Institucional"
	case strings.Contains(k, "ambient"):
		return "Eixo Aspectos Ambientais"
	}
	return strings.TrimSpace(raw)
}

// normalizePrioridade unifica variações de grafia da prioridade.
func normalizePrioridade(raw string) string {
	switch chave(raw) {
	case "baixo", "baixa":
		return "Baixa"
	case "alta", "alto":
		return "Alta"
	case "media", "medio":
		return "Média"
	}
	return strings.TrimSpace(raw)
}
