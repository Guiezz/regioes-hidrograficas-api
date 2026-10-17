package main

import "testing"

func TestNormalizeEixoMatriz(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Matriz de ações e programas do eixo de demanda hídrica da RH do Coreaú", "Eixo Demanda Hídrica"},
		{"Programas e ações do eixo oferta hídrica", "Eixo Oferta Hídrica"},
		{"Matriz de Ações e Programas do Eixo de Gerenciamento das Águas", "Eixo Gerenciamento das Águas"},
		{"Matriz de ações e programas do eixo político-institucional", "Eixo Político-Institucional"},
		{"Matriz de Ações e Programas do Eixo Politico-Institucional", "Eixo Político-Institucional"},
		{"Matriz de ações e programas do eixo Aspectos Ambientais", "Eixo Aspectos Ambientais"},
		{"Eixo Demanda Hídrica", "Eixo Demanda Hídrica"},
		{"  Texto qualquer  ", "Texto qualquer"},
	}
	for _, c := range cases {
		if got := normalizeEixoMatriz(c.in); got != c.want {
			t.Errorf("normalizeEixoMatriz(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizePrioridade(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Baixo", "Baixa"},
		{"baixa ", "Baixa"},
		{"Alta", "Alta"},
		{"Média", "Média"},
		{"media", "Média"},
		{"", ""},
	}
	for _, c := range cases {
		if got := normalizePrioridade(c.in); got != c.want {
			t.Errorf("normalizePrioridade(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
