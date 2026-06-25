package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// carregarTrace lê um arquivo de trace e devolve a sequência de acessos, em ordem.
//
// Formato (um acesso por linha; '#' inicia comentário):
//
//	P1 R 0x1A2B4
//	P2 W 74565
//	P1 R 12288
//
// Campos: <processo> <R|W> <endereço (hex 0x.. ou decimal)>.
func carregarTrace(caminho string) ([]Acesso, error) {
	f, err := os.Open(caminho)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var acessos []Acesso
	sc := bufio.NewScanner(f)
	linha := 0
	for sc.Scan() {
		linha++
		texto := strings.TrimSpace(sc.Text())
		if i := strings.IndexByte(texto, '#'); i >= 0 {
			texto = strings.TrimSpace(texto[:i])
		}
		if texto == "" {
			continue
		}
		campos := strings.Fields(texto)
		if len(campos) != 3 {
			return nil, fmt.Errorf("linha %d: esperados 3 campos, lidos %d", linha, len(campos))
		}

		proc, err := strconv.Atoi(strings.TrimPrefix(strings.ToUpper(campos[0]), "P"))
		if err != nil {
			return nil, fmt.Errorf("linha %d: processo inválido %q", linha, campos[0])
		}

		var tipo TipoAcesso
		switch strings.ToUpper(campos[1]) {
		case "R":
			tipo = Leitura
		case "W":
			tipo = Escrita
		default:
			return nil, fmt.Errorf("linha %d: tipo inválido %q (use R ou W)", linha, campos[1])
		}

		end, err := parseEndereco(campos[2])
		if err != nil {
			return nil, fmt.Errorf("linha %d: %v", linha, err)
		}
		if end >= TamMemoriaVirtual {
			return nil, fmt.Errorf("linha %d: endereço %d fora do espaço virtual (max %d)", linha, end, TamMemoriaVirtual-1)
		}

		acessos = append(acessos, Acesso{Processo: proc, Tipo: tipo, Endereco: EnderecoVirtual(end)})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(acessos) == 0 {
		return nil, fmt.Errorf("trace vazio")
	}
	return acessos, nil
}

func parseEndereco(s string) (int, error) {
	s = strings.ToLower(s)
	if strings.HasPrefix(s, "0x") {
		v, err := strconv.ParseInt(s[2:], 16, 64)
		return int(v), err
	}
	v, err := strconv.Atoi(s)
	return v, err
}
