package main

// FIFO substitui a página que está há mais tempo carregada (ordem de carga).
type FIFO struct {
	fila []int // frames na ordem em que receberam páginas
}

func NovoFIFO() *FIFO { return &FIFO{} }

func (f *FIFO) Nome() string         { return "FIFO" }
func (f *FIFO) AoAcessar(pagina int) {} // FIFO ignora acessos; só importa a carga

func (f *FIFO) AoCarregar(pagina, frame int) {
	f.fila = append(f.fila, frame)
}

func (f *FIFO) EscolherVitima(e EstadoFrames) int {
	frame := f.fila[0]
	f.fila = f.fila[1:]
	return frame
}
