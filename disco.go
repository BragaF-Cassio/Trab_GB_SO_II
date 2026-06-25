package main

import "fmt"

// Disco simula o backing store: 1 MB representando a "memória virtual em repouso".
// É um array em memória — fiel ao conceito de cópia disco→RAM sem I/O real.
type Disco struct {
	dados []byte
}

// NovoDisco cria o disco e o preenche com conteúdo sintético autoverificável.
func NovoDisco() *Disco {
	d := &Disco{dados: make([]byte, TamMemoriaVirtual)}
	d.inicializarConteudo()
	return d
}

// inicializarConteudo preenche cada página com uma marca derivada do seu número,
// de forma que, ao exibir o conteúdo, dê para conferir de qual página ele veio.
func (d *Disco) inicializarConteudo() {
	for pag := 0; pag < NumPaginas; pag++ {
		marca := []byte(fmt.Sprintf("[pag %03d]", pag))
		base := pag * TamBloco
		for i := 0; i < TamBloco; i++ {
			d.dados[base+i] = marca[i%len(marca)]
		}
	}
}

// LerBloco copia os 8 KB de uma página do disco (usado ao carregar na RAM).
func (d *Disco) LerBloco(pagina int) []byte {
	base := pagina * TamBloco
	bloco := make([]byte, TamBloco)
	copy(bloco, d.dados[base:base+TamBloco])
	return bloco
}

// EscreverBloco grava os 8 KB de um frame de volta no disco (write-back).
func (d *Disco) EscreverBloco(pagina int, dados []byte) {
	base := pagina * TamBloco
	copy(d.dados[base:base+TamBloco], dados)
}
