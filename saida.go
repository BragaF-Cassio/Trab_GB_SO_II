package main

import (
	"fmt"
	"strings"
	"sync"
)

// Códigos ANSI para coloração do terminal.
const (
	cReset = "\033[0m"
	cVerm  = "\033[31m"
	cVerde = "\033[32m"
	cAmar  = "\033[33m"
	cAzul  = "\033[34m"
	cCiano = "\033[36m"
	cCinza = "\033[90m"
	cNegr  = "\033[1m"
)

// FrameInfo é um retrato (cópia) do estado de um frame para a saída.
type FrameInfo struct {
	Indice   int
	Pagina   int
	R, D     bool
	Ponteiro bool
}

// Stats agrega os contadores da execução.
type Stats struct {
	Total, Hits, Faltas, WriteBacks int
}

// Evento descreve tudo que aconteceu num acesso, para a goroutine de saída.
type Evento struct {
	Acesso         Acesso
	Pagina, Offset int
	Frame          int
	EnderecoFisico uint32
	Hit, Falta     bool
	Substituiu     bool
	VitimaPagina   int
	VitimaFrame    int
	WriteBack      bool
	Conteudo       string
	Snapshot       []FrameInfo
	Stats          Stats
}

// Saida é a goroutine consumidora dos eventos (separa renderização da lógica).
type Saida struct {
	cor      bool
	snapshot bool
}

func (s *Saida) c(texto, codigo string) string {
	if !s.cor {
		return texto
	}
	return codigo + texto + cReset
}

// Consumir lê eventos e imprime log + snapshots até o canal fechar.
func (s *Saida) Consumir(eventos <-chan Evento, wg *sync.WaitGroup) {
	defer wg.Done()
	for ev := range eventos {
		s.imprimirLinha(ev)
		if ev.Falta && s.snapshot {
			s.imprimirSnapshot(ev)
		}
	}
}

func (s *Saida) imprimirLinha(ev Evento) {
	cab := fmt.Sprintf("#%05d P%d %s", ev.Acesso.Seq, ev.Acesso.Processo, ev.Acesso.Tipo)
	virt := fmt.Sprintf("v=0x%05X [%s]", uint32(ev.Acesso.Endereco), ev.Acesso.Endereco.DecomposicaoBits())

	if ev.Hit {
		res := s.c("HIT  ", cVerde)
		fis := fmt.Sprintf("→ f=%d  fís=0x%04X", ev.Frame, ev.EnderecoFisico)
		fmt.Printf("%s  %s  %s %s  %s\n",
			cab, virt, res, fis, s.c(fmt.Sprintf("conteúdo=%q", ev.Conteudo), cCinza))
		return
	}

	// Falta de página
	res := s.c("FALTA", cVerm)
	var detalhe string
	if ev.Substituiu {
		sub := fmt.Sprintf("substitui pág %d (f=%d)", ev.VitimaPagina, ev.VitimaFrame)
		if ev.WriteBack {
			sub += s.c(" +write-back", cCiano)
		}
		detalhe = s.c(sub, cAmar)
	} else {
		detalhe = s.c("frame livre", cAmar)
	}
	fmt.Printf("%s  %s  %s  carrega pág %d → f=%d  [%s]  fís=0x%04X  %s\n",
		cab, virt, res, ev.Pagina, ev.Frame, detalhe, ev.EnderecoFisico,
		s.c(fmt.Sprintf("conteúdo=%q", ev.Conteudo), cCinza))
}

const boxIndent = "        " // recuo do snapshot
const boxInner = 54          // largura interna da moldura (colunas visíveis)

// rowFmt fixa as colunas do snapshot: frame · faixa física · página · faixa
// virtual · bits R/D. Cabeçalho e linhas usam o MESMO formato → alinhamento.
const rowFmt = " %-4s %-13s %4s  %-15s  %s"

func (s *Saida) imprimirSnapshot(ev Evento) {
	var b strings.Builder

	// borda superior com título "frames"
	titulo := "─ frames "
	b.WriteString(boxIndent)
	b.WriteString(s.c("┌"+titulo+strings.Repeat("─", boxInner-len([]rune(titulo)))+"┐\n", cCinza))

	// cabeçalho das colunas
	cab := fmt.Sprintf(rowFmt, "frm", "faixa física", "pág", "faixa virtual", "R D")
	b.WriteString(s.linhaBox(cab, len([]rune(cab))))

	for _, fi := range ev.Snapshot {
		fisica := FaixaFisicaFrame(fi.Indice)
		var conteudo string
		if fi.Pagina >= 0 {
			conteudo = fmt.Sprintf(rowFmt, fmt.Sprintf("[%d]", fi.Indice), fisica,
				fmt.Sprintf("%d", fi.Pagina), FaixaVirtualPagina(fi.Pagina),
				fmt.Sprintf("%d %d", b2i(fi.R), b2i(fi.D)))
		} else {
			conteudo = fmt.Sprintf(rowFmt, fmt.Sprintf("[%d]", fi.Indice), fisica,
				"--", "(livre)", "")
		}
		vis := len([]rune(conteudo))
		if fi.Ponteiro { // ponteiro do Clock
			conteudo += "  " + s.c("← ●", cAmar)
			vis += 5
		}
		b.WriteString(s.linhaBox(conteudo, vis))
	}

	// rodapé com contadores
	taxa := 0.0
	if ev.Stats.Total > 0 {
		taxa = 100 * float64(ev.Stats.Faltas) / float64(ev.Stats.Total)
	}
	rod := fmt.Sprintf(" hits=%d  faltas=%d  wb=%d  falta=%.1f%%",
		ev.Stats.Hits, ev.Stats.Faltas, ev.Stats.WriteBacks, taxa)
	b.WriteString(s.linhaBox(s.c(rod, cCinza), len([]rune(rod))))

	// borda inferior
	b.WriteString(boxIndent)
	b.WriteString(s.c("└"+strings.Repeat("─", boxInner)+"┘\n", cCinza))

	fmt.Print(b.String())
}

// linhaBox imprime uma linha da moldura. visivel é a largura aparente do
// conteúdo (descontando códigos de cor), usada para alinhar a borda direita.
func (s *Saida) linhaBox(conteudo string, visivel int) string {
	preenche := boxInner - visivel
	if preenche < 0 {
		preenche = 0
	}
	return boxIndent + s.c("│", cCinza) + conteudo + strings.Repeat(" ", preenche) + s.c("│", cCinza) + "\n"
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
