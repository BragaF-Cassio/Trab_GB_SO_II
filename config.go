package main

// Parâmetros do hardware fictício, conforme o enunciado.
//
//	Memória virtual ...... 1 MB
//	Memória principal .... 64 KB
//	Bloco (página/frame) . 8 KB
//
// Disso derivam:
//
//	128 páginas virtuais (1 MB / 8 KB)
//	  8 frames físicos   (64 KB / 8 KB)
//	 20 bits de endereço virtual = 7 (página) + 13 (offset)
//	 16 bits de endereço físico  = 3 (frame)  + 13 (offset)
const (
	TamMemoriaVirtual = 1 << 20                      // 1 MB = 1.048.576 bytes
	TamMemoriaFisica  = 64 * 1024                    // 64 KB = 65.536 bytes
	TamBloco          = 8 * 1024                     // 8 KB = 8.192 bytes
	NumPaginas        = TamMemoriaVirtual / TamBloco // 128
	NumFrames         = TamMemoriaFisica / TamBloco  // 8
	BitsOffset        = 13                           // log2(8192)
	MascaraOffset     = TamBloco - 1                 // 0x1FFF
)
