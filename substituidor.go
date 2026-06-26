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

// Visualizavel é opcional: algoritmos com um ponteiro a exibir (Clock).
type Visualizavel interface {
	Ponteiro() int
}
