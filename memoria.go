package main

// Frame é um quadro físico de 8 KB na memória principal.
type Frame struct {
	Pagina int    // página atualmente carregada (-1 = frame livre)
	Dados  []byte // 8 KB de conteúdo
}

// MemoriaFisica é a RAM de 64 KB = 8 frames.
type MemoriaFisica struct {
	Frames [NumFrames]Frame
}

// NovaMemoriaFisica cria a RAM com todos os frames livres.
func NovaMemoriaFisica() *MemoriaFisica {
	m := &MemoriaFisica{}
	for i := range m.Frames {
		m.Frames[i].Pagina = -1
		m.Frames[i].Dados = make([]byte, TamBloco)
	}
	return m
}

// FrameLivre devolve o índice do primeiro frame livre, ou -1 se a RAM estiver cheia.
func (m *MemoriaFisica) FrameLivre() int {
	for i := range m.Frames {
		if m.Frames[i].Pagina == -1 {
			return i
		}
	}
	return -1
}
