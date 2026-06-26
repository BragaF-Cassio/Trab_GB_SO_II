package main

import "fmt"

// EnderecoVirtual é um endereço de 20 bits no espaço virtual de 1 MB.
type EnderecoVirtual uint32

// Pagina extrai os 7 bits altos (número da página) via deslocamento.
func (e EnderecoVirtual) Pagina() int { return int(e >> BitsOffset) }

// Offset extrai os 13 bits baixos (deslocamento dentro do bloco) via máscara.
func (e EnderecoVirtual) Offset() int { return int(e) & MascaraOffset }

// EnderecoFisico monta o endereço físico de 16 bits: o offset é preservado,
// apenas o número da página é trocado pelo número do frame.
func EnderecoFisico(frame, offset int) uint32 {
	return uint32(frame)<<BitsOffset | uint32(offset)
}

// DecomposicaoBits mostra os 20 bits do endereço virtual fatiados em
// 7 bits de página + 13 bits de offset — útil para a explicação no vídeo.
func (e EnderecoVirtual) DecomposicaoBits() string {
	return fmt.Sprintf("%07b·%013b", e.Pagina(), e.Offset())
}

// FaixaFisicaFrame devolve a faixa de endereços físicos (16 bits) coberta por um
// frame: do byte base até o último byte do bloco de 8 KB.
func FaixaFisicaFrame(frame int) string {
	base := frame * TamBloco
	return fmt.Sprintf("0x%04X-0x%04X", base, base+TamBloco-1)
}

// FaixaVirtualPagina devolve a faixa de endereços virtuais (20 bits) coberta por
// uma página: do byte base até o último byte do bloco de 8 KB.
func FaixaVirtualPagina(pagina int) string {
	base := pagina * TamBloco
	return fmt.Sprintf("0x%05X-0x%05X", base, base+TamBloco-1)
}
