package main

// EstadoFrames é o "contexto leve" que a MMU expõe ao substituidor no momento
// de escolher a vítima. É a fonte ÚNICA da verdade (a própria tabela de páginas):
// o algoritmo lê/zera os bits reais, em vez de manter cópias que poderiam
// dessincronizar.
type EstadoFrames interface {
	NumFrames() int
	PaginaEm(frame int) int       // página atualmente no frame
	BitReferencia(frame int) bool // lê o bit R da PTE
	LimparReferencia(frame int)   // zera o bit R (segunda chance do Clock)
	BitSujeira(frame int) bool    // lê o bit D da PTE (usado pelo Clock-NRU)
}

// Substituidor é a interface plugável dos algoritmos de substituição.
// O estado EXCLUSIVo de cada algoritmo (fila do FIFO, recência do LRU,
// ponteiro do Clock) é mantido internamente e alimentado por eventos.
type Substituidor interface {
	Nome() string
	AoAcessar(pagina int)              // todo acesso (hit ou pós-carga)
	AoCarregar(pagina, frame int)      // página recém-colocada num frame
	EscolherVitima(e EstadoFrames) int // devolve o frame a ser substituído
}

// Visualizavel é opcional: algoritmos com um ponteiro a exibir (Clock, NRU).
type Visualizavel interface {
	Ponteiro() int
}

// criarSubstituidor instancia o algoritmo pelo nome. O OPT precisa da sequência
// futura (só disponível no modo de comparação), por isso é tratado à parte.
func criarSubstituidor(nome string) Substituidor {
	switch nome {
	case "fifo":
		return NovoFIFO()
	case "lru":
		return NovoLRU()
	case "clock":
		return NovoClock()
	case "nru":
		return NovoNRU()
	default:
		return nil
	}
}
