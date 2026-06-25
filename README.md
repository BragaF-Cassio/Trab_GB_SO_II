# Simulador de Paginação (Memória Virtual)

Simulador de **memória virtual paginada** escrito em Go para a disciplina de
Sistemas Operacionais: Análise e Aplicações (UNISINOS). 

Nome: Cássio Ferreira Braga

A MMU traduz endereços virtuais em físicos,
trata **faltas de página** e aplica um algoritmo de **substituição** plugável
(protagonista: **Clock / Segunda Chance**), com **write-back** de páginas sujas.

---

## Parâmetros do hardware fictício

| Parâmetro                 | Valor    | Derivação                       |
|---------------------------|----------|---------------------------------|
| Memória virtual           | 1 MB     | **128 páginas** de 8 KB         |
| Memória principal (RAM)   | 64 KB    | **8 frames** de 8 KB            |
| Bloco (página = frame)    | 8 KB     | offset de **13 bits**           |
| Endereço virtual          | 20 bits  | 7 (página) + 13 (offset)        |
| Endereço físico           | 16 bits  | 3 (frame) + 13 (offset)         |

A tradução é feita por **operações de bits** (deslocamento + máscara), não por
aritmética: `página = vaddr >> 13`, `offset = vaddr & 0x1FFF`.

---

## Como compilar e rodar

Requer Go 1.22+. Em ambientes sem `GOCACHE/GOPATH` configurados:

```bash
export GOCACHE=$HOME/.cache/go-build
export GOPATH=$HOME/go

go build -o simulador .
./simulador
```

Não há dependências externas (as cores usam códigos ANSI próprios).

### Exemplos

```bash
# Modo aleatório (padrão): 2 processos leves, algoritmo Clock
./simulador

# Mais acessos, semente fixa (reprodutível), proporção de escritas 40%
./simulador -acessos 60 -semente 7 -escrita 0.4

# Cenário controlado a partir de um trace (ordem exata dos acessos)
./simulador -trace exemplo.trace

# Comparar TODOS os algoritmos sobre a mesma sequência de acessos
./simulador -acessos 80 -comparar -snapshot=false

# Trocar o algoritmo protagonista da execução ao vivo
./simulador -alg lru        # clock | fifo | lru | nru
```

### Flags

| Flag         | Padrão  | Descrição                                            |
|--------------|---------|------------------------------------------------------|
| `-alg`       | `clock` | algoritmo da execução ao vivo: `clock fifo lru nru`  |
| `-processos` | `2`     | número de processos leves (mínimo 2)                 |
| `-acessos`   | `25`    | acessos por processo (modo aleatório)                |
| `-semente`   | `42`    | semente do gerador (reprodutibilidade)               |
| `-escrita`   | `0.30`  | proporção de escritas (0..1)                         |
| `-trace`     | `""`    | arquivo de trace (cenário controlado)                |
| `-comparar`  | `false` | ao final, compara todos os algoritmos no mesmo trace |
| `-cor`       | `true`  | cores ANSI na saída                                  |
| `-snapshot`  | `true`  | imprime o estado dos frames a cada falta de página   |

---

## Decisões de projeto

1. **Linguagem Go** — goroutines + channels modelam diretamente "processos leves"
   e o padrão **produtor/consumidor**.
2. **Concorrência sem mutex** — os processos *produzem* acessos num canal; a **MMU
   é a única consumidora** e a única a tocar o estado compartilhado. A saída é
   outra goroutine consumindo um canal de eventos.
3. **Espaço virtual único e compartilhado** — 1 MB / 128 páginas, com **uma**
   tabela de páginas global. Os processos colidem nas mesmas páginas (memória
   compartilhada), o que torna a substituição mais interessante.
4. **Tradução por bits** — deslocamento e máscara (ver `endereco.go`).
5. **Geração híbrida de acessos** — aleatória com semente (modelo de localidade)
   **ou** trace de arquivo. Alocação de frames é global.
6. **Substituição plugável** — interface `Substituidor`; **Clock** é o
   protagonista, com FIFO, LRU, Clock-NRU e OPT como comparáveis.
7. **Disco simulado em memória** (1 MB) com **conteúdo sintético derivado do
   endereço** (`[pag NNN]`), tornando a saída autoverificável.
8. **Leitura e escrita** com **bit de sujeira (D)** e **write-back** de páginas
   sujas na substituição.
9. **Saída em log + snapshots tabulares** (com cores) a cada falta de página.
10. **Interface híbrida do substituidor** — o algoritmo recebe *eventos*
    (`AoAcessar`, `AoCarregar`) para seu estado privado e um *contexto leve*
    (`EstadoFrames`) para escolher a vítima lendo/zerando os bits **reais** da
    tabela (fonte única da verdade).
11. **Clock simples** (só bit de referência) como protagonista.

---

## Estrutura dos arquivos

```
config.go        constantes derivadas do enunciado
endereco.go      decomposição de endereços (bitwise)
tabela.go        tabela de páginas global (PTEs: presente, frame, R, D)
memoria.go       memória física (8 frames de 8 KB)
disco.go         backing store de 1 MB com conteúdo sintético
processo.go      processos leves (produtores) com localidade
trace.go         leitor de trace de arquivo
substituidor.go  interfaces Substituidor / EstadoFrames / Visualizavel
fifo.go lru.go clock.go nru.go optimal.go   algoritmos
mmu.go           núcleo: tradução, falta de página, substituição, write-back
saida.go         goroutine de saída: log colorido + snapshots
comparar.go      execução silenciosa e tabela comparativa
main.go          flags e orquestração produtor/consumidor
exemplo.trace    cenário demonstrativo controlado
```

---

## Como ler a saída

Cada linha de acesso mostra: sequência, processo, tipo (R/W), endereço virtual e
sua **decomposição em bits** (`página·offset`), o resultado (**HIT** ou
**FALTA**), o frame resolvido, o endereço físico e uma amostra do **conteúdo**.

A cada **falta de página** é impresso um *snapshot* dos 8 frames com os bits
**R** e **D** e a posição do **ponteiro do Clock** (`← ●`). O write-back de uma
página suja aparece como `+write-back`.

O `exemplo.trace` foi montado para, em sequência: encher os 8 frames (deixando
páginas sujas), dar segunda chance via Clock ao acessar uma nova página, e
**remover uma página suja com write-back** — recarregando-a depois para mostrar
que o conteúdo escrito **persistiu** no disco.

---

## Comparação de algoritmos

`./simulador -comparar` roda a **mesma** sequência de acessos por FIFO, LRU,
Clock, Clock-NRU e OPT, imprimindo faltas, write-backs e taxas. O **OPT**
(algoritmo de Belady) é irrealizável num SO real — serve de **limite inferior**
de faltas para situar os demais.
