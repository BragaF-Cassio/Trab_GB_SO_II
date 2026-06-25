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
