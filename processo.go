package main

import (
	"math/rand"
	"sync"
)

// TipoAcesso distingue leitura de escrita.
type TipoAcesso int

const (
	Leitura TipoAcesso = iota
	Escrita
)

func (t TipoAcesso) String() string {
	if t == Escrita {
		return "W"
	}
	return "R"
}

// Acesso é uma instrução gerada por um processo.
type Acesso struct {
	Processo int
	Tipo     TipoAcesso
	Endereco EnderecoVirtual
	Seq      int // nº de sequência global (preenchido pela MMU ao consumir)
}

// Processo é um "processo leve": uma goroutine produtora que gera acessos
// sobre o espaço virtual COMPARTILHADO e os envia ao canal (produtor/consumidor).
// Usa um modelo simples de localidade para que os algoritmos de substituição
// tenham comportamento interessante (e não puramente aleatório).
type Processo struct {
	ID          int
	rng         *rand.Rand
	qtd         int
	propEscrita float64
}

// NovoProcesso cria um processo com gerador próprio (semente derivada = reprodutível).
func NovoProcesso(id, qtd int, semente int64, propEscrita float64) *Processo {
	return &Processo{
		ID:          id,
		rng:         rand.New(rand.NewSource(semente)),
		qtd:         qtd,
		propEscrita: propEscrita,
	}
}

// Produzir gera os acessos e os envia ao canal, sinalizando o WaitGroup ao terminar.
func (p *Processo) Produzir(canal chan<- Acesso, wg *sync.WaitGroup) {
	defer wg.Done()
	base := p.rng.Intn(NumPaginas) // página-base da região de trabalho atual
	const janela = 6               // tamanho da localidade (em páginas)
	for i := 0; i < p.qtd; i++ {
		if p.rng.Float64() < 0.20 { // 20%: salta para outra região (compartilhada)
			base = p.rng.Intn(NumPaginas)
		}
		pag := (base + p.rng.Intn(janela)) % NumPaginas
		offset := p.rng.Intn(TamBloco)
		tipo := Leitura
		if p.rng.Float64() < p.propEscrita {
			tipo = Escrita
		}
		canal <- Acesso{
			Processo: p.ID,
			Tipo:     tipo,
			Endereco: EnderecoVirtual(pag*TamBloco + offset),
		}
	}
}
