package main

import "fmt"

// MMU é a unidade de gerência de memória e o ÚNICO consumidor do canal de
// acessos. Por ser a única goroutine a tocar a tabela de páginas, os frames e o
// disco, não há necessidade de mutex.
type MMU struct {
	tabela  *TabelaPaginas
	memoria *MemoriaFisica
	disco   *Disco
	alg     Substituidor
	eventos chan<- Evento // nil no modo de comparação (silencioso)

	// estatísticas
	seq        int
	hits       int
	faltas     int
	writebacks int

	consumidos []Acesso // ordem real de consumo (para comparação posterior)
}

// NovaMMU monta a MMU com um algoritmo de substituição e um canal de eventos.
func NovaMMU(alg Substituidor, eventos chan<- Evento) *MMU {
	return &MMU{
		tabela:  NovaTabela(),
		memoria: NovaMemoriaFisica(),
		disco:   NovoDisco(),
		alg:     alg,
		eventos: eventos,
	}
}

// --- EstadoFrames: a MMU é a fonte da verdade exposta ao substituidor ---

func (m *MMU) NumFrames() int         { return NumFrames }
func (m *MMU) PaginaEm(frame int) int { return m.memoria.Frames[frame].Pagina }
func (m *MMU) BitReferencia(f int) bool {
	p := m.memoria.Frames[f].Pagina
	return p >= 0 && m.tabela[p].Referencia
}
func (m *MMU) BitSujeira(f int) bool {
	p := m.memoria.Frames[f].Pagina
	return p >= 0 && m.tabela[p].Suja
}
func (m *MMU) LimparReferencia(f int) {
	if p := m.memoria.Frames[f].Pagina; p >= 0 {
		m.tabela[p].Referencia = false
	}
}

// Processar traduz um acesso, tratando hit ou falta de página, e emite um evento.
func (m *MMU) Processar(a Acesso) {
	m.seq++
	a.Seq = m.seq
	m.consumidos = append(m.consumidos, a)

	pag := a.Endereco.Pagina()
	off := a.Endereco.Offset()
	m.alg.AoAcessar(pag)

	ev := Evento{Acesso: a, Pagina: pag, Offset: off}
	pte := &m.tabela[pag]

	if pte.Presente {
		// ---- HIT ----
		m.hits++
		ev.Hit = true
		ev.Frame = pte.Frame
		pte.Referencia = true
		if a.Tipo == Escrita {
			pte.Suja = true
			m.escrever(pte.Frame, off, a)
		}
	} else {
		// ---- FALTA DE PÁGINA ----
		m.faltas++
		ev.Falta = true
		frame := m.memoria.FrameLivre()
		if frame == -1 {
			// Sem frame livre: substitui.
			frame = m.alg.EscolherVitima(m)
			vit := m.memoria.Frames[frame].Pagina
			ev.Substituiu = true
			ev.VitimaPagina = vit
			ev.VitimaFrame = frame
			if m.tabela[vit].Suja { // write-back de página suja
				m.disco.EscreverBloco(vit, m.memoria.Frames[frame].Dados)
				m.writebacks++
				ev.WriteBack = true
			}
			m.tabela[vit] = EntradaTabela{Frame: -1} // invalida a vítima
		}
		// Carrega o bloco do disco para o frame.
		copy(m.memoria.Frames[frame].Dados, m.disco.LerBloco(pag))
		m.memoria.Frames[frame].Pagina = pag
		*pte = EntradaTabela{Presente: true, Frame: frame, Referencia: true}
		m.alg.AoCarregar(pag, frame)
		ev.Frame = frame
		if a.Tipo == Escrita {
			pte.Suja = true
			m.escrever(frame, off, a)
		}
	}

	ev.EnderecoFisico = EnderecoFisico(ev.Frame, off)
	ev.Conteudo = m.ler(ev.Frame, off)

	if m.eventos != nil {
		ev.Snapshot = m.snapshot()
		ev.Stats = m.stats()
		m.eventos <- ev
	}
}

// escrever altera o conteúdo da página na RAM (deixa uma marca visível), de modo
// que o write-back e a persistência fiquem demonstráveis na saída.
func (m *MMU) escrever(frame, offset int, a Acesso) {
	marca := []byte(fmt.Sprintf("<W:P%d>", a.Processo))
	dados := m.memoria.Frames[frame].Dados
	for i := 0; i < len(marca) && offset+i < TamBloco; i++ {
		dados[offset+i] = marca[i]
	}
}

// ler devolve uma amostra legível do conteúdo no endereço físico resolvido.
func (m *MMU) ler(frame, offset int) string {
	dados := m.memoria.Frames[frame].Dados
	fim := offset + 12
	if fim > TamBloco {
		fim = TamBloco
	}
	return string(dados[offset:fim])
}

func (m *MMU) snapshot() []FrameInfo {
	ponteiro := -1
	if v, ok := m.alg.(Visualizavel); ok {
		ponteiro = v.Ponteiro()
	}
	s := make([]FrameInfo, NumFrames)
	for i := 0; i < NumFrames; i++ {
		pag := m.memoria.Frames[i].Pagina
		fi := FrameInfo{Indice: i, Pagina: pag, Ponteiro: i == ponteiro}
		if pag >= 0 {
			fi.R = m.tabela[pag].Referencia
			fi.D = m.tabela[pag].Suja
		}
		s[i] = fi
	}
	return s
}

func (m *MMU) stats() Stats {
	return Stats{Total: m.seq, Hits: m.hits, Faltas: m.faltas, WriteBacks: m.writebacks}
}
