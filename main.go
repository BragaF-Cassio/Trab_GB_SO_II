package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"sync"
)

func main() {
	var (
		fProcessos = flag.Int("processos", 2, "número de processos leves (mín. 2)")
		fAcessos   = flag.Int("acessos", 25, "acessos por processo (modo aleatório)")
		fSemente   = flag.Int64("semente", 42, "semente do gerador (reprodutibilidade)")
		fEscrita   = flag.Float64("escrita", 0.30, "proporção de escritas (0..1)")
		fTrace     = flag.String("trace", "", "arquivo de trace (modo cenário controlado)")
		fCor       = flag.Bool("cor", true, "usa cores ANSI na saída")
		fSnapshot  = flag.Bool("snapshot", true, "imprime o estado dos frames a cada falta de página")
	)
	flag.Parse()

	if *fProcessos < 2 {
		*fProcessos = 2 // o enunciado exige no mínimo dois processos leves
	}

	alg := NovoClock()

	cabecalho(alg.Nome(), *fProcessos, *fTrace)

	canalAcessos := make(chan Acesso, 64)
	canalEventos := make(chan Evento, 64)

	// Consumidor da saída (renderização separada da lógica).
	saida := &Saida{cor: *fCor, snapshot: *fSnapshot}
	var wgSaida sync.WaitGroup
	wgSaida.Add(1)
	go saida.Consumir(canalEventos, &wgSaida)

	// Produtores: trace (1 alimentador, ordem exata) ou aleatório (N goroutines).
	var wgProd sync.WaitGroup
	if *fTrace != "" {
		acessos, err := carregarTrace(*fTrace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro no trace: %v\n", err)
			os.Exit(1)
		}
		wgProd.Add(1)
		go func() {
			defer wgProd.Done()
			for _, a := range acessos {
				canalAcessos <- a
			}
		}()
	} else {
		for id := 1; id <= *fProcessos; id++ {
			p := NovoProcesso(id, *fAcessos, *fSemente+int64(id), *fEscrita)
			wgProd.Add(1)
			go p.Produzir(canalAcessos, &wgProd)
		}
	}

	// Fecha o canal de acessos quando todos os produtores terminarem.
	go func() {
		wgProd.Wait()
		close(canalAcessos)
	}()

	// MMU: o único consumidor dos acessos (sem mutex sobre o estado).
	mmu := NovaMMU(alg, canalEventos)
	for a := range canalAcessos {
		mmu.Processar(a)
	}
	close(canalEventos)
	wgSaida.Wait()

	resumoFinal(mmu)
}

func cabecalho(alg string, processos int, trace string) {
	fonte := fmt.Sprintf("aleatório, %d processos leves", processos)
	if trace != "" {
		fonte = "trace: " + trace
	}
	fmt.Printf("%sSimulador de Paginação%s  —  virtual 1 MB (128 pág) · física 64 KB (8 frames) · bloco 8 KB\n",
		cNegr, cReset)
	fmt.Printf("algoritmo: %s%s%s   fonte: %s\n\n", cNegr, alg, cReset, fonte)
}

func resumoFinal(m *MMU) {
	st := m.stats()
	taxa := 0.0
	if st.Total > 0 {
		taxa = 100 * float64(st.Faltas) / float64(st.Total)
	}
	fmt.Printf("\n%sResumo%s  acessos=%d  hits=%d  faltas=%d  write-backs=%d  taxa de falta=%.1f%%\n",
		cNegr, cReset, st.Total, st.Hits, st.Faltas, st.WriteBacks, taxa)

	// Distribuição final dos frames.
	var ocupados []int
	for i := 0; i < NumFrames; i++ {
		if p := m.memoria.Frames[i].Pagina; p >= 0 {
			ocupados = append(ocupados, p)
		}
	}
	sort.Ints(ocupados)
	fmt.Printf("frames ao final contêm as páginas: %v\n", ocupados)
}
