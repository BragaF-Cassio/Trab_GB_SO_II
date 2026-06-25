package main

// LRU substitui a página usada há mais tempo. Mantém o "relógio lógico" do
// último acesso de cada página.
type LRU struct {
	relogio int
	ultimo  map[int]int // página -> instante do último acesso
}

func NovoLRU() *LRU { return &LRU{ultimo: make(map[int]int)} }

func (l *LRU) Nome() string { return "LRU" }

func (l *LRU) AoAcessar(pagina int) {
	l.relogio++
	l.ultimo[pagina] = l.relogio
}

func (l *LRU) AoCarregar(pagina, frame int) {
	l.relogio++
	l.ultimo[pagina] = l.relogio
}

func (l *LRU) EscolherVitima(e EstadoFrames) int {
	melhorFrame := -1
	melhorTempo := int(^uint(0) >> 1) // maxInt
	for i := 0; i < e.NumFrames(); i++ {
		pag := e.PaginaEm(i)
		if t := l.ultimo[pag]; t < melhorTempo {
			melhorTempo = t
			melhorFrame = i
		}
	}
	return melhorFrame
}
