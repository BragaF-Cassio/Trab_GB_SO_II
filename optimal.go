package main

// Optimal (algoritmo de Belady) substitui a página cujo próximo uso está mais
// distante no futuro. É irrealizável num SO real (exige conhecer o futuro), mas
// como dispomos da sequência completa de acessos, serve de REFERÊNCIA teórica:
// nenhum outro algoritmo pode ter menos faltas que o OPT.
type Optimal struct {
	futuro []int // sequência de páginas de toda a execução
	pos    int   // índice do acesso atual
}

// NovoOptimal recebe a sequência completa de páginas acessadas.
func NovoOptimal(sequencia []int) *Optimal {
	return &Optimal{futuro: sequencia}
}

func (o *Optimal) Nome() string                 { return "OPT" }
func (o *Optimal) AoAcessar(pagina int)         { o.pos++ }
func (o *Optimal) AoCarregar(pagina, frame int) {}

func (o *Optimal) EscolherVitima(e EstadoFrames) int {
	melhorFrame := -1
	melhorDist := -1
	for i := 0; i < e.NumFrames(); i++ {
		pag := e.PaginaEm(i)
		if d := o.proximoUso(pag); d > melhorDist {
			melhorDist = d
			melhorFrame = i
		}
	}
	return melhorFrame
}

// proximoUso devolve a distância até o próximo acesso à página (quanto maior,
// melhor candidata a sair). Páginas nunca mais usadas recebem distância máxima.
func (o *Optimal) proximoUso(pagina int) int {
	for j := o.pos; j < len(o.futuro); j++ {
		if o.futuro[j] == pagina {
			return j - o.pos
		}
	}
	return 1 << 30
}
