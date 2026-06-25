package main

// Clock (Segunda Chance) é o algoritmo protagonista.
//
// Mantém apenas um ponteiro circular sobre os frames. O bit de referência vive
// na tabela de páginas (a MMU o seta a cada acesso); aqui só o lemos e zeramos
// através do contexto EstadoFrames — sem cópia duplicada.
//
//	Ao procurar a vítima, o ponteiro avança:
//	  R == 0  -> esta é a vítima
//	  R == 1  -> zera o bit (segunda chance) e avança
type Clock struct {
	ponteiro int
}

func NovoClock() *Clock { return &Clock{} }

func (c *Clock) Nome() string                 { return "Clock" }
func (c *Clock) AoAcessar(pagina int)         {} // o bit R é setado pela MMU na PTE
func (c *Clock) AoCarregar(pagina, frame int) {}
func (c *Clock) Ponteiro() int                { return c.ponteiro }

func (c *Clock) EscolherVitima(e EstadoFrames) int {
	n := e.NumFrames()
	for {
		f := c.ponteiro
		c.ponteiro = (c.ponteiro + 1) % n
		if !e.BitReferencia(f) {
			return f
		}
		e.LimparReferencia(f) // segunda chance
	}
}
