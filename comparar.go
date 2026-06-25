package main

import "fmt"

// rodarSilencioso processa toda a sequência com um algoritmo, sem emitir eventos,
// e devolve as estatísticas finais (estado de memória novo a cada chamada).
func rodarSilencioso(alg Substituidor, acessos []Acesso) Stats {
	mmu := NovaMMU(alg, nil)
	for _, a := range acessos {
		mmu.Processar(a)
	}
	return mmu.stats()
}

// compararAlgoritmos roda a MESMA sequência de acessos por todos os algoritmos
// (incluindo o OPT como referência teórica) e imprime uma tabela comparativa.
func compararAlgoritmos(acessos []Acesso, cor bool) {
	paginas := make([]int, len(acessos))
	for i, a := range acessos {
		paginas[i] = a.Endereco.Pagina()
	}

	tipo := func(a Substituidor) Stats { return rodarSilencioso(a, acessos) }

	linhas := []struct {
		nome string
		st   Stats
	}{
		{"FIFO", tipo(NovoFIFO())},
		{"LRU", tipo(NovoLRU())},
		{"Clock", tipo(NovoClock())},
		{"Clock-NRU", tipo(NovoNRU())},
		{"OPT (ótimo)", tipo(NovoOptimal(paginas))},
	}

	c := func(s, code string) string {
		if !cor {
			return s
		}
		return code + s + cReset
	}

	fmt.Println()
	fmt.Println(c("══ Comparação de algoritmos (mesmo trace) ══", cNegr))
	fmt.Printf("%-14s %8s %8s %8s %12s %12s\n", "algoritmo", "acessos", "faltas", "wb", "taxa falta", "taxa hit")
	fmt.Println("────────────────────────────────────────────────────────────────────")
	for _, l := range linhas {
		taxaFalta := 100 * float64(l.st.Faltas) / float64(l.st.Total)
		taxaHit := 100 * float64(l.st.Hits) / float64(l.st.Total)
		fmt.Printf("%-14s %8d %8d %8d %11.1f%% %11.1f%%\n",
			l.nome, l.st.Total, l.st.Faltas, l.st.WriteBacks, taxaFalta, taxaHit)
	}
	fmt.Println("────────────────────────────────────────────────────────────────────")
	fmt.Println(c("Nota: o OPT é irrealizável num SO real (exige o futuro); serve de limite inferior de faltas.", cCinza))
}
