package main

// NRU (Clock melhorado / Segunda Chance aprimorado) considera o par (R, D) e
// prefere remover páginas nesta ordem, economizando write-backs:
//
//	(R=0,D=0) não usada e limpa   -> ideal, sem write-back
//	(R=0,D=1) não usada mas suja
//	(R=1,D=0) usada e limpa
//	(R=1,D=1) usada e suja
//
// É feito em até quatro varreduras; na segunda, zera os bits R no caminho.
type NRU struct {
	ponteiro int
}

func NovoNRU() *NRU { return &NRU{} }

func (c *NRU) Nome() string                 { return "Clock-NRU" }
func (c *NRU) AoAcessar(pagina int)         {}
func (c *NRU) AoCarregar(pagina, frame int) {}
func (c *NRU) Ponteiro() int                { return c.ponteiro }

func (c *NRU) EscolherVitima(e EstadoFrames) int {
	n := e.NumFrames()
	inicio := c.ponteiro

	achar := func(refAlvo, sujaAlvo bool, limparR bool) int {
		for i := 0; i < n; i++ {
			f := (inicio + i) % n
			if e.BitReferencia(f) == refAlvo && e.BitSujeira(f) == sujaAlvo {
				c.ponteiro = (f + 1) % n
				return f
			}
			if limparR {
				e.LimparReferencia(f)
			}
		}
		return -1
	}

	// Passo 1: (0,0) sem alterar bits.
	if f := achar(false, false, false); f >= 0 {
		return f
	}
	// Passo 2: (0,1) zerando R no caminho.
	if f := achar(false, true, true); f >= 0 {
		return f
	}
	// Passo 3: agora todos têm R=0; reprocura (0,0).
	if f := achar(false, false, false); f >= 0 {
		return f
	}
	// Passo 4: (0,1).
	if f := achar(false, true, false); f >= 0 {
		return f
	}
	// Fallback (não deve ocorrer).
	f := inicio
	c.ponteiro = (f + 1) % n
	return f
}
