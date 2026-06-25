package main

// EntradaTabela é uma entrada (PTE) da tabela de páginas global.
// Como o espaço virtual é único e compartilhado, basta UMA tabela
// indexada pelo número da página (0..127).
type EntradaTabela struct {
	Presente   bool // a página está em algum frame da RAM?
	Frame      int  // qual frame (válido apenas se Presente)
	Referencia bool // bit R — setado a cada acesso; usado pelo Clock
	Suja       bool // bit D — setado em escritas; exige write-back na saída
}

// TabelaPaginas é a tabela global de 128 entradas.
type TabelaPaginas [NumPaginas]EntradaTabela

// NovaTabela cria a tabela com todas as páginas ausentes.
func NovaTabela() *TabelaPaginas {
	t := &TabelaPaginas{}
	for i := range t {
		t[i].Frame = -1
	}
	return t
}
