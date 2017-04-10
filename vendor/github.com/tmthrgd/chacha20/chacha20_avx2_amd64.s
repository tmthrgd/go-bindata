// Created by chacha20_avx2.pl - DO NOT EDIT
// perl chacha20_avx2.pl golang-no-avx chacha20_avx2_amd64.s

// +build amd64,!gccgo,!appengine

#include "textflag.h"

DATA chacha20_consts<>+0x00(SB)/1, $"e"
DATA chacha20_consts<>+0x01(SB)/1, $"x"
DATA chacha20_consts<>+0x02(SB)/1, $"p"
DATA chacha20_consts<>+0x03(SB)/1, $"a"
DATA chacha20_consts<>+0x04(SB)/1, $"n"
DATA chacha20_consts<>+0x05(SB)/1, $"d"
DATA chacha20_consts<>+0x06(SB)/1, $" "
DATA chacha20_consts<>+0x07(SB)/1, $"3"
DATA chacha20_consts<>+0x08(SB)/1, $"2"
DATA chacha20_consts<>+0x09(SB)/1, $"-"
DATA chacha20_consts<>+0x0a(SB)/1, $"b"
DATA chacha20_consts<>+0x0b(SB)/1, $"y"
DATA chacha20_consts<>+0x0c(SB)/1, $"t"
DATA chacha20_consts<>+0x0d(SB)/1, $"e"
DATA chacha20_consts<>+0x0e(SB)/1, $" "
DATA chacha20_consts<>+0x0f(SB)/1, $"k"
DATA chacha20_consts<>+0x10(SB)/1, $"e"
DATA chacha20_consts<>+0x11(SB)/1, $"x"
DATA chacha20_consts<>+0x12(SB)/1, $"p"
DATA chacha20_consts<>+0x13(SB)/1, $"a"
DATA chacha20_consts<>+0x14(SB)/1, $"n"
DATA chacha20_consts<>+0x15(SB)/1, $"d"
DATA chacha20_consts<>+0x16(SB)/1, $" "
DATA chacha20_consts<>+0x17(SB)/1, $"3"
DATA chacha20_consts<>+0x18(SB)/1, $"2"
DATA chacha20_consts<>+0x19(SB)/1, $"-"
DATA chacha20_consts<>+0x1a(SB)/1, $"b"
DATA chacha20_consts<>+0x1b(SB)/1, $"y"
DATA chacha20_consts<>+0x1c(SB)/1, $"t"
DATA chacha20_consts<>+0x1d(SB)/1, $"e"
DATA chacha20_consts<>+0x1e(SB)/1, $" "
DATA chacha20_consts<>+0x1f(SB)/1, $"k"
GLOBL chacha20_consts<>(SB), RODATA, $32

DATA rol8<>+0x00(SB)/1, $3
DATA rol8<>+0x01(SB)/1, $0
DATA rol8<>+0x02(SB)/1, $1
DATA rol8<>+0x03(SB)/1, $2
DATA rol8<>+0x04(SB)/1, $7
DATA rol8<>+0x05(SB)/1, $4
DATA rol8<>+0x06(SB)/1, $5
DATA rol8<>+0x07(SB)/1, $6
DATA rol8<>+0x08(SB)/1, $11
DATA rol8<>+0x09(SB)/1, $8
DATA rol8<>+0x0a(SB)/1, $9
DATA rol8<>+0x0b(SB)/1, $10
DATA rol8<>+0x0c(SB)/1, $15
DATA rol8<>+0x0d(SB)/1, $12
DATA rol8<>+0x0e(SB)/1, $13
DATA rol8<>+0x0f(SB)/1, $14
DATA rol8<>+0x10(SB)/1, $3
DATA rol8<>+0x11(SB)/1, $0
DATA rol8<>+0x12(SB)/1, $1
DATA rol8<>+0x13(SB)/1, $2
DATA rol8<>+0x14(SB)/1, $7
DATA rol8<>+0x15(SB)/1, $4
DATA rol8<>+0x16(SB)/1, $5
DATA rol8<>+0x17(SB)/1, $6
DATA rol8<>+0x18(SB)/1, $11
DATA rol8<>+0x19(SB)/1, $8
DATA rol8<>+0x1a(SB)/1, $9
DATA rol8<>+0x1b(SB)/1, $10
DATA rol8<>+0x1c(SB)/1, $15
DATA rol8<>+0x1d(SB)/1, $12
DATA rol8<>+0x1e(SB)/1, $13
DATA rol8<>+0x1f(SB)/1, $14
GLOBL rol8<>(SB), RODATA, $32

DATA rol16<>+0x00(SB)/1, $2
DATA rol16<>+0x01(SB)/1, $3
DATA rol16<>+0x02(SB)/1, $0
DATA rol16<>+0x03(SB)/1, $1
DATA rol16<>+0x04(SB)/1, $6
DATA rol16<>+0x05(SB)/1, $7
DATA rol16<>+0x06(SB)/1, $4
DATA rol16<>+0x07(SB)/1, $5
DATA rol16<>+0x08(SB)/1, $10
DATA rol16<>+0x09(SB)/1, $11
DATA rol16<>+0x0a(SB)/1, $8
DATA rol16<>+0x0b(SB)/1, $9
DATA rol16<>+0x0c(SB)/1, $14
DATA rol16<>+0x0d(SB)/1, $15
DATA rol16<>+0x0e(SB)/1, $12
DATA rol16<>+0x0f(SB)/1, $13
DATA rol16<>+0x10(SB)/1, $2
DATA rol16<>+0x11(SB)/1, $3
DATA rol16<>+0x12(SB)/1, $0
DATA rol16<>+0x13(SB)/1, $1
DATA rol16<>+0x14(SB)/1, $6
DATA rol16<>+0x15(SB)/1, $7
DATA rol16<>+0x16(SB)/1, $4
DATA rol16<>+0x17(SB)/1, $5
DATA rol16<>+0x18(SB)/1, $10
DATA rol16<>+0x19(SB)/1, $11
DATA rol16<>+0x1a(SB)/1, $8
DATA rol16<>+0x1b(SB)/1, $9
DATA rol16<>+0x1c(SB)/1, $14
DATA rol16<>+0x1d(SB)/1, $15
DATA rol16<>+0x1e(SB)/1, $12
DATA rol16<>+0x1f(SB)/1, $13
GLOBL rol16<>(SB), RODATA, $32

DATA avx2Init<>+0x00(SB)/8, $0x0
DATA avx2Init<>+0x08(SB)/8, $0x0
DATA avx2Init<>+0x10(SB)/8, $0x1
DATA avx2Init<>+0x18(SB)/8, $0x0
GLOBL avx2Init<>(SB), RODATA, $32

DATA avx2Inc<>+0x00(SB)/8, $0x2
DATA avx2Inc<>+0x08(SB)/8, $0x0
DATA avx2Inc<>+0x10(SB)/8, $0x2
DATA avx2Inc<>+0x18(SB)/8, $0x0
GLOBL avx2Inc<>(SB), RODATA, $32

TEXT ·chacha_20_core_avx2(SB),$0-32
	MOVQ	out+0(FP),DI
	MOVQ	in+8(FP),SI
	MOVQ	in_len+16(FP),DX
	MOVQ	state+24(FP),BX

	MOVQ	$chacha20_consts<>(SB),R11
	MOVQ	$rol8<>(SB),R12
	MOVQ	$rol16<>(SB),R13
	MOVQ	$avx2Init<>(SB),R14
	MOVQ	$avx2Inc<>(SB),R15

	VZEROUPPER


	// VBROADCASTI128	16*0(BX),Y0
	BYTE $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x5a; BYTE $0x03
	// VBROADCASTI128	16*1(BX),Y1
	BYTE $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x5a; BYTE $0x4b; BYTE $0x10
	// VBROADCASTI128	16*2(BX),Y2
	BYTE $0xc4; BYTE $0xe2; BYTE $0x7d; BYTE $0x5a; BYTE $0x53; BYTE $0x20
	// VPADDQ	(R14),Y2,Y2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x6d; BYTE $0xd4; BYTE $0x16

label2a:
	CMPQ	DX,$384
	JB	label2b

	// VMOVDQA	(R11),Y4
	BYTE $0xc4; BYTE $0xc1; BYTE $0x7d; BYTE $0x6f; BYTE $0x23
	// VMOVDQA	(R11),Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x7d; BYTE $0x6f; BYTE $0x03
	// VMOVDQA	(R11),Y12
	BYTE $0xc4; BYTE $0x41; BYTE $0x7d; BYTE $0x6f; BYTE $0x23

	// VMOVDQA	Y0,Y5
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xe8
	// VMOVDQA	Y0,Y9
	BYTE $0xc5; BYTE $0x7d; BYTE $0x6f; BYTE $0xc8
	// VMOVDQA	Y0,Y13
	BYTE $0xc5; BYTE $0x7d; BYTE $0x6f; BYTE $0xe8

	// VMOVDQA	Y1,Y6
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xf1
	// VMOVDQA	Y1,Y10
	BYTE $0xc5; BYTE $0x7d; BYTE $0x6f; BYTE $0xd1
	// VMOVDQA	Y1,Y14
	BYTE $0xc5; BYTE $0x7d; BYTE $0x6f; BYTE $0xf1

	// VMOVDQA	Y2,Y7
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xfa
	// VPADDQ	(R15),Y7,Y11
	BYTE $0xc4; BYTE $0x41; BYTE $0x45; BYTE $0xd4; BYTE $0x1f
	// VPADDQ	(R15),Y11,Y15
	BYTE $0xc4; BYTE $0x41; BYTE $0x25; BYTE $0xd4; BYTE $0x3f

	MOVQ	$10,R8

label1a:

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R13),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5
	// VPSLLD	$12,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R12),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5

	// VPSLLD	$7,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R13),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x5d; BYTE $0x00

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9
	// VPSLLD	$12,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x0c
	// VPSRLD	$20,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x14
	VPXOR	Y3,Y9,Y9

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R12),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x1c; BYTE $0x24

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9

	// VPSLLD	$7,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x07
	// VPSRLD	$25,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x19
	VPXOR	Y3,Y9,Y9

	// VPADDD	Y13,Y12,Y12
	BYTE $0xc4; BYTE $0x41; BYTE $0x1d; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y12,Y15,Y15
	// VPSHUFB	(R13),Y15,Y15
	BYTE $0xc4; BYTE $0x42; BYTE $0x05; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y15,Y14,Y14
	BYTE $0xc4; BYTE $0x41; BYTE $0x0d; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y14,Y13,Y13
	// VPSLLD	$12,Y13,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y13,Y13
	BYTE $0xc4; BYTE $0xc1; BYTE $0x15; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y13,Y13

	// VPADDD	Y13,Y12,Y12
	BYTE $0xc4; BYTE $0x41; BYTE $0x1d; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y12,Y15,Y15
	// VPSHUFB	(R12),Y15,Y15
	BYTE $0xc4; BYTE $0x42; BYTE $0x05; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y15,Y14,Y14
	BYTE $0xc4; BYTE $0x41; BYTE $0x0d; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y14,Y13,Y13

	// VPSLLD	$7,Y13,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y13,Y13
	BYTE $0xc4; BYTE $0xc1; BYTE $0x15; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y13,Y13
	// VPALIGNR	$4,Y5,Y5,Y5
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x0f; BYTE $0xed; BYTE $0x04
	// VPALIGNR	$8,Y6,Y6,Y6
	BYTE $0xc4; BYTE $0xe3; BYTE $0x4d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$12,Y7,Y7,Y7
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x0f; BYTE $0xff; BYTE $0x0c
	// VPALIGNR	$4,Y9,Y9,Y9
	BYTE $0xc4; BYTE $0x43; BYTE $0x35; BYTE $0x0f; BYTE $0xc9; BYTE $0x04
	// VPALIGNR	$8,Y10,Y10,Y10
	BYTE $0xc4; BYTE $0x43; BYTE $0x2d; BYTE $0x0f; BYTE $0xd2; BYTE $0x08
	// VPALIGNR	$12,Y11,Y11,Y11
	BYTE $0xc4; BYTE $0x43; BYTE $0x25; BYTE $0x0f; BYTE $0xdb; BYTE $0x0c
	// VPALIGNR	$4,Y13,Y13,Y13
	BYTE $0xc4; BYTE $0x43; BYTE $0x15; BYTE $0x0f; BYTE $0xed; BYTE $0x04
	// VPALIGNR	$8,Y14,Y14,Y14
	BYTE $0xc4; BYTE $0x43; BYTE $0x0d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$12,Y15,Y15,Y15
	BYTE $0xc4; BYTE $0x43; BYTE $0x05; BYTE $0x0f; BYTE $0xff; BYTE $0x0c

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R13),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5
	// VPSLLD	$12,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R12),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5

	// VPSLLD	$7,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R13),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x5d; BYTE $0x00

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9
	// VPSLLD	$12,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x0c
	// VPSRLD	$20,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x14
	VPXOR	Y3,Y9,Y9

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R12),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x1c; BYTE $0x24

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9

	// VPSLLD	$7,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x07
	// VPSRLD	$25,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x19
	VPXOR	Y3,Y9,Y9

	// VPADDD	Y13,Y12,Y12
	BYTE $0xc4; BYTE $0x41; BYTE $0x1d; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y12,Y15,Y15
	// VPSHUFB	(R13),Y15,Y15
	BYTE $0xc4; BYTE $0x42; BYTE $0x05; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y15,Y14,Y14
	BYTE $0xc4; BYTE $0x41; BYTE $0x0d; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y14,Y13,Y13
	// VPSLLD	$12,Y13,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y13,Y13
	BYTE $0xc4; BYTE $0xc1; BYTE $0x15; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y13,Y13

	// VPADDD	Y13,Y12,Y12
	BYTE $0xc4; BYTE $0x41; BYTE $0x1d; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y12,Y15,Y15
	// VPSHUFB	(R12),Y15,Y15
	BYTE $0xc4; BYTE $0x42; BYTE $0x05; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y15,Y14,Y14
	BYTE $0xc4; BYTE $0x41; BYTE $0x0d; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y14,Y13,Y13

	// VPSLLD	$7,Y13,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y13,Y13
	BYTE $0xc4; BYTE $0xc1; BYTE $0x15; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y13,Y13
	// VPALIGNR	$12,Y5,Y5,Y5
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x0f; BYTE $0xed; BYTE $0x0c
	// VPALIGNR	$8,Y6,Y6,Y6
	BYTE $0xc4; BYTE $0xe3; BYTE $0x4d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$4,Y7,Y7,Y7
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x0f; BYTE $0xff; BYTE $0x04
	// VPALIGNR	$12,Y9,Y9,Y9
	BYTE $0xc4; BYTE $0x43; BYTE $0x35; BYTE $0x0f; BYTE $0xc9; BYTE $0x0c
	// VPALIGNR	$8,Y10,Y10,Y10
	BYTE $0xc4; BYTE $0x43; BYTE $0x2d; BYTE $0x0f; BYTE $0xd2; BYTE $0x08
	// VPALIGNR	$4,Y11,Y11,Y11
	BYTE $0xc4; BYTE $0x43; BYTE $0x25; BYTE $0x0f; BYTE $0xdb; BYTE $0x04
	// VPALIGNR	$12,Y13,Y13,Y13
	BYTE $0xc4; BYTE $0x43; BYTE $0x15; BYTE $0x0f; BYTE $0xed; BYTE $0x0c
	// VPALIGNR	$8,Y14,Y14,Y14
	BYTE $0xc4; BYTE $0x43; BYTE $0x0d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$4,Y15,Y15,Y15
	BYTE $0xc4; BYTE $0x43; BYTE $0x05; BYTE $0x0f; BYTE $0xff; BYTE $0x04

	DECQ	R8

	JNZ	label1a

	// VPADDD	(R11),Y4,Y4
	BYTE $0xc4; BYTE $0xc1; BYTE $0x5d; BYTE $0xfe; BYTE $0x23
	// VPADDD	(R11),Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0x03
	// VPADDD	(R11),Y12,Y12
	BYTE $0xc4; BYTE $0x41; BYTE $0x1d; BYTE $0xfe; BYTE $0x23

	// VPADDD	Y0,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0xfe; BYTE $0xe8
	// VPADDD	Y0,Y9,Y9
	BYTE $0xc5; BYTE $0x35; BYTE $0xfe; BYTE $0xc8
	// VPADDD	Y0,Y13,Y13
	BYTE $0xc5; BYTE $0x15; BYTE $0xfe; BYTE $0xe8

	// VPADDD	Y1,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf1
	// VPADDD	Y1,Y10,Y10
	BYTE $0xc5; BYTE $0x2d; BYTE $0xfe; BYTE $0xd1
	// VPADDD	Y1,Y14,Y14
	BYTE $0xc5; BYTE $0x0d; BYTE $0xfe; BYTE $0xf1

	// VPADDD	Y2,Y7,Y7
	BYTE $0xc5; BYTE $0xc5; BYTE $0xfe; BYTE $0xfa
	// VPADDQ	(R15),Y2,Y2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x6d; BYTE $0xd4; BYTE $0x17
	// VPADDD	Y2,Y11,Y11
	BYTE $0xc5; BYTE $0x25; BYTE $0xfe; BYTE $0xda
	// VPADDQ	(R15),Y2,Y2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x6d; BYTE $0xd4; BYTE $0x17
	// VPADDD	Y2,Y15,Y15
	BYTE $0xc5; BYTE $0x05; BYTE $0xfe; BYTE $0xfa
	// VPADDQ	(R15),Y2,Y2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x6d; BYTE $0xd4; BYTE $0x17

	// VPERM2I128	$2,Y4,Y5,Y3
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x46; BYTE $0xdc; BYTE $0x02
	VPXOR	32*0(SI),Y3,Y3
	VMOVDQU	Y3,32*0(DI)
	// VPERM2I128	$2,Y6,Y7,Y3
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x46; BYTE $0xde; BYTE $0x02
	VPXOR	32*1(SI),Y3,Y3
	VMOVDQU	Y3,32*1(DI)
	// VPERM2I128	$19,Y4,Y5,Y3
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x46; BYTE $0xdc; BYTE $0x13
	VPXOR	32*2(SI),Y3,Y3
	VMOVDQU	Y3,32*2(DI)
	// VPERM2I128	$19,Y6,Y7,Y3
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x46; BYTE $0xde; BYTE $0x13
	VPXOR	32*3(SI),Y3,Y3
	VMOVDQU	Y3,32*3(DI)

	// VPERM2I128	$2,Y8,Y9,Y4
	BYTE $0xc4; BYTE $0xc3; BYTE $0x35; BYTE $0x46; BYTE $0xe0; BYTE $0x02
	// VPERM2I128	$2,Y10,Y11,Y5
	BYTE $0xc4; BYTE $0xc3; BYTE $0x25; BYTE $0x46; BYTE $0xea; BYTE $0x02
	// VPERM2I128	$19,Y8,Y9,Y6
	BYTE $0xc4; BYTE $0xc3; BYTE $0x35; BYTE $0x46; BYTE $0xf0; BYTE $0x13
	// VPERM2I128	$19,Y10,Y11,Y7
	BYTE $0xc4; BYTE $0xc3; BYTE $0x25; BYTE $0x46; BYTE $0xfa; BYTE $0x13

	VPXOR	32*4(SI),Y4,Y4
	VPXOR	32*5(SI),Y5,Y5
	VPXOR	32*6(SI),Y6,Y6
	VPXOR	32*7(SI),Y7,Y7

	VMOVDQU	Y4,32*4(DI)
	VMOVDQU	Y5,32*5(DI)
	VMOVDQU	Y6,32*6(DI)
	VMOVDQU	Y7,32*7(DI)

	// VPERM2I128	$2,Y12,Y13,Y4
	BYTE $0xc4; BYTE $0xc3; BYTE $0x15; BYTE $0x46; BYTE $0xe4; BYTE $0x02
	// VPERM2I128	$2,Y14,Y15,Y5
	BYTE $0xc4; BYTE $0xc3; BYTE $0x05; BYTE $0x46; BYTE $0xee; BYTE $0x02
	// VPERM2I128	$19,Y12,Y13,Y6
	BYTE $0xc4; BYTE $0xc3; BYTE $0x15; BYTE $0x46; BYTE $0xf4; BYTE $0x13
	// VPERM2I128	$19,Y14,Y15,Y7
	BYTE $0xc4; BYTE $0xc3; BYTE $0x05; BYTE $0x46; BYTE $0xfe; BYTE $0x13

	VPXOR	32*8(SI),Y4,Y4
	VPXOR	32*9(SI),Y5,Y5
	VPXOR	32*10(SI),Y6,Y6
	VPXOR	32*11(SI),Y7,Y7

	VMOVDQU	Y4,32*8(DI)
	VMOVDQU	Y5,32*9(DI)
	VMOVDQU	Y6,32*10(DI)
	VMOVDQU	Y7,32*11(DI)

	LEAQ	64*6(SI),SI
	LEAQ	64*6(DI),DI
	SUBQ	$384,DX

	JMP	label2a

label2b:
	CMPQ	DX,$256
	JB	label2c

	// VMOVDQA	(R11),Y4
	BYTE $0xc4; BYTE $0xc1; BYTE $0x7d; BYTE $0x6f; BYTE $0x23
	// VMOVDQA	(R11),Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x7d; BYTE $0x6f; BYTE $0x03
	// VMOVDQA	Y0,Y5
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xe8
	// VMOVDQA	Y0,Y9
	BYTE $0xc5; BYTE $0x7d; BYTE $0x6f; BYTE $0xc8
	// VMOVDQA	Y1,Y6
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xf1
	// VMOVDQA	Y1,Y10
	BYTE $0xc5; BYTE $0x7d; BYTE $0x6f; BYTE $0xd1
	// VMOVDQA	Y1,Y14
	BYTE $0xc5; BYTE $0x7d; BYTE $0x6f; BYTE $0xf1
	// VMOVDQA	Y2,Y7
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xfa
	// VPADDQ	(R15),Y7,Y11
	BYTE $0xc4; BYTE $0x41; BYTE $0x45; BYTE $0xd4; BYTE $0x1f

	MOVQ	$10,R8

label1b:

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R13),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5
	// VPSLLD	$12,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R12),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5

	// VPSLLD	$7,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R13),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x5d; BYTE $0x00

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9
	// VPSLLD	$12,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x0c
	// VPSRLD	$20,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x14
	VPXOR	Y3,Y9,Y9

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R12),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x1c; BYTE $0x24

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9

	// VPSLLD	$7,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x07
	// VPSRLD	$25,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x19
	VPXOR	Y3,Y9,Y9
	// VPALIGNR	$4,Y5,Y5,Y5
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x0f; BYTE $0xed; BYTE $0x04
	// VPALIGNR	$8,Y6,Y6,Y6
	BYTE $0xc4; BYTE $0xe3; BYTE $0x4d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$12,Y7,Y7,Y7
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x0f; BYTE $0xff; BYTE $0x0c
	// VPALIGNR	$4,Y9,Y9,Y9
	BYTE $0xc4; BYTE $0x43; BYTE $0x35; BYTE $0x0f; BYTE $0xc9; BYTE $0x04
	// VPALIGNR	$8,Y10,Y10,Y10
	BYTE $0xc4; BYTE $0x43; BYTE $0x2d; BYTE $0x0f; BYTE $0xd2; BYTE $0x08
	// VPALIGNR	$12,Y11,Y11,Y11
	BYTE $0xc4; BYTE $0x43; BYTE $0x25; BYTE $0x0f; BYTE $0xdb; BYTE $0x0c

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R13),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5
	// VPSLLD	$12,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R12),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5

	// VPSLLD	$7,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R13),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x5d; BYTE $0x00

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9
	// VPSLLD	$12,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x0c
	// VPSRLD	$20,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x14
	VPXOR	Y3,Y9,Y9

	// VPADDD	Y9,Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0xc1
	VPXOR	Y8,Y11,Y11
	// VPSHUFB	(R12),Y11,Y11
	BYTE $0xc4; BYTE $0x42; BYTE $0x25; BYTE $0x00; BYTE $0x1c; BYTE $0x24

	// VPADDD	Y11,Y10,Y10
	BYTE $0xc4; BYTE $0x41; BYTE $0x2d; BYTE $0xfe; BYTE $0xd3
	VPXOR	Y10,Y9,Y9

	// VPSLLD	$7,Y9,Y3
	BYTE $0xc4; BYTE $0xc1; BYTE $0x65; BYTE $0x72; BYTE $0xf1; BYTE $0x07
	// VPSRLD	$25,Y9,Y9
	BYTE $0xc4; BYTE $0xc1; BYTE $0x35; BYTE $0x72; BYTE $0xd1; BYTE $0x19
	VPXOR	Y3,Y9,Y9
	// VPALIGNR	$12,Y5,Y5,Y5
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x0f; BYTE $0xed; BYTE $0x0c
	// VPALIGNR	$8,Y6,Y6,Y6
	BYTE $0xc4; BYTE $0xe3; BYTE $0x4d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$4,Y7,Y7,Y7
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x0f; BYTE $0xff; BYTE $0x04
	// VPALIGNR	$12,Y9,Y9,Y9
	BYTE $0xc4; BYTE $0x43; BYTE $0x35; BYTE $0x0f; BYTE $0xc9; BYTE $0x0c
	// VPALIGNR	$8,Y10,Y10,Y10
	BYTE $0xc4; BYTE $0x43; BYTE $0x2d; BYTE $0x0f; BYTE $0xd2; BYTE $0x08
	// VPALIGNR	$4,Y11,Y11,Y11
	BYTE $0xc4; BYTE $0x43; BYTE $0x25; BYTE $0x0f; BYTE $0xdb; BYTE $0x04

	DECQ	R8

	JNZ	label1b

	// VPADDD	(R11),Y4,Y4
	BYTE $0xc4; BYTE $0xc1; BYTE $0x5d; BYTE $0xfe; BYTE $0x23
	// VPADDD	(R11),Y8,Y8
	BYTE $0xc4; BYTE $0x41; BYTE $0x3d; BYTE $0xfe; BYTE $0x03

	// VPADDD	Y0,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0xfe; BYTE $0xe8
	// VPADDD	Y0,Y9,Y9
	BYTE $0xc5; BYTE $0x35; BYTE $0xfe; BYTE $0xc8

	// VPADDD	Y1,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf1
	// VPADDD	Y1,Y10,Y10
	BYTE $0xc5; BYTE $0x2d; BYTE $0xfe; BYTE $0xd1

	// VPADDD	Y2,Y7,Y7
	BYTE $0xc5; BYTE $0xc5; BYTE $0xfe; BYTE $0xfa
	// VPADDQ	(R15),Y2,Y2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x6d; BYTE $0xd4; BYTE $0x17
	// VPADDD	Y2,Y11,Y11
	BYTE $0xc5; BYTE $0x25; BYTE $0xfe; BYTE $0xda
	// VPADDQ	(R15),Y2,Y2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x6d; BYTE $0xd4; BYTE $0x17

	// VPERM2I128	$2,Y4,Y5,Y12
	BYTE $0xc4; BYTE $0x63; BYTE $0x55; BYTE $0x46; BYTE $0xe4; BYTE $0x02
	// VPERM2I128	$2,Y6,Y7,Y13
	BYTE $0xc4; BYTE $0x63; BYTE $0x45; BYTE $0x46; BYTE $0xee; BYTE $0x02
	// VPERM2I128	$19,Y4,Y5,Y14
	BYTE $0xc4; BYTE $0x63; BYTE $0x55; BYTE $0x46; BYTE $0xf4; BYTE $0x13
	// VPERM2I128	$19,Y6,Y7,Y15
	BYTE $0xc4; BYTE $0x63; BYTE $0x45; BYTE $0x46; BYTE $0xfe; BYTE $0x13

	VPXOR	32*0(SI),Y12,Y12
	VPXOR	32*1(SI),Y13,Y13
	VPXOR	32*2(SI),Y14,Y14
	VPXOR	32*3(SI),Y15,Y15

	VMOVDQU	Y12,32*0(DI)
	VMOVDQU	Y13,32*1(DI)
	VMOVDQU	Y14,32*2(DI)
	VMOVDQU	Y15,32*3(DI)

	// VPERM2I128	$2,Y8,Y9,Y4
	BYTE $0xc4; BYTE $0xc3; BYTE $0x35; BYTE $0x46; BYTE $0xe0; BYTE $0x02
	// VPERM2I128	$2,Y10,Y11,Y5
	BYTE $0xc4; BYTE $0xc3; BYTE $0x25; BYTE $0x46; BYTE $0xea; BYTE $0x02
	// VPERM2I128	$19,Y8,Y9,Y6
	BYTE $0xc4; BYTE $0xc3; BYTE $0x35; BYTE $0x46; BYTE $0xf0; BYTE $0x13
	// VPERM2I128	$19,Y10,Y11,Y7
	BYTE $0xc4; BYTE $0xc3; BYTE $0x25; BYTE $0x46; BYTE $0xfa; BYTE $0x13

	VPXOR	32*4(SI),Y4,Y4
	VPXOR	32*5(SI),Y5,Y5
	VPXOR	32*6(SI),Y6,Y6
	VPXOR	32*7(SI),Y7,Y7

	VMOVDQU	Y4,32*4(DI)
	VMOVDQU	Y5,32*5(DI)
	VMOVDQU	Y6,32*6(DI)
	VMOVDQU	Y7,32*7(DI)

	LEAQ	64*4(SI),SI
	LEAQ	64*4(DI),DI
	SUBQ	$256,DX

	JMP	label2b
label2c:
	CMPQ	DX,$128
	JB	label2d

	// VMOVDQA	(R11),Y4
	BYTE $0xc4; BYTE $0xc1; BYTE $0x7d; BYTE $0x6f; BYTE $0x23
	// VMOVDQA	Y0,Y5
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xe8
	// VMOVDQA	Y1,Y6
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xf1
	// VMOVDQA	Y2,Y7
	BYTE $0xc5; BYTE $0xfd; BYTE $0x6f; BYTE $0xfa

	MOVQ	$10,R8

label1c:

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R13),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5
	// VPSLLD	$12,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R12),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5

	// VPSLLD	$7,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y5,Y5
	// VPALIGNR	$4,Y5,Y5,Y5
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x0f; BYTE $0xed; BYTE $0x04
	// VPALIGNR	$8,Y6,Y6,Y6
	BYTE $0xc4; BYTE $0xe3; BYTE $0x4d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$12,Y7,Y7,Y7
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x0f; BYTE $0xff; BYTE $0x0c

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R13),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x7d; BYTE $0x00

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5
	// VPSLLD	$12,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x0c
	// VPSRLD	$20,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x14
	VPXOR	Y3,Y5,Y5

	// VPADDD	Y5,Y4,Y4
	BYTE $0xc5; BYTE $0xdd; BYTE $0xfe; BYTE $0xe5
	VPXOR	Y4,Y7,Y7
	// VPSHUFB	(R12),Y7,Y7
	BYTE $0xc4; BYTE $0xc2; BYTE $0x45; BYTE $0x00; BYTE $0x3c; BYTE $0x24

	// VPADDD	Y7,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf7
	VPXOR	Y6,Y5,Y5

	// VPSLLD	$7,Y5,Y3
	BYTE $0xc5; BYTE $0xe5; BYTE $0x72; BYTE $0xf5; BYTE $0x07
	// VPSRLD	$25,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0x72; BYTE $0xd5; BYTE $0x19
	VPXOR	Y3,Y5,Y5
	// VPALIGNR	$12,Y5,Y5,Y5
	BYTE $0xc4; BYTE $0xe3; BYTE $0x55; BYTE $0x0f; BYTE $0xed; BYTE $0x0c
	// VPALIGNR	$8,Y6,Y6,Y6
	BYTE $0xc4; BYTE $0xe3; BYTE $0x4d; BYTE $0x0f; BYTE $0xf6; BYTE $0x08
	// VPALIGNR	$4,Y7,Y7,Y7
	BYTE $0xc4; BYTE $0xe3; BYTE $0x45; BYTE $0x0f; BYTE $0xff; BYTE $0x04

	DECQ	R8
	JNZ	label1c

	// VPADDD	(R11),Y4,Y4
	BYTE $0xc4; BYTE $0xc1; BYTE $0x5d; BYTE $0xfe; BYTE $0x23
	// VPADDD	Y0,Y5,Y5
	BYTE $0xc5; BYTE $0xd5; BYTE $0xfe; BYTE $0xe8
	// VPADDD	Y1,Y6,Y6
	BYTE $0xc5; BYTE $0xcd; BYTE $0xfe; BYTE $0xf1
	// VPADDD	Y2,Y7,Y7
	BYTE $0xc5; BYTE $0xc5; BYTE $0xfe; BYTE $0xfa
	// VPADDQ	(R15),Y2,Y2
	BYTE $0xc4; BYTE $0xc1; BYTE $0x6d; BYTE $0xd4; BYTE $0x17

	// VPERM2I128	$2,Y4,Y5,Y12
	BYTE $0xc4; BYTE $0x63; BYTE $0x55; BYTE $0x46; BYTE $0xe4; BYTE $0x02
	// VPERM2I128	$2,Y6,Y7,Y13
	BYTE $0xc4; BYTE $0x63; BYTE $0x45; BYTE $0x46; BYTE $0xee; BYTE $0x02
	// VPERM2I128	$19,Y4,Y5,Y14
	BYTE $0xc4; BYTE $0x63; BYTE $0x55; BYTE $0x46; BYTE $0xf4; BYTE $0x13
	// VPERM2I128	$19,Y6,Y7,Y15
	BYTE $0xc4; BYTE $0x63; BYTE $0x45; BYTE $0x46; BYTE $0xfe; BYTE $0x13

	VPXOR	32*0(SI),Y12,Y12
	VPXOR	32*1(SI),Y13,Y13
	VPXOR	32*2(SI),Y14,Y14
	VPXOR	32*3(SI),Y15,Y15

	VMOVDQU	Y12,32*0(DI)
	VMOVDQU	Y13,32*1(DI)
	VMOVDQU	Y14,32*2(DI)
	VMOVDQU	Y15,32*3(DI)

	LEAQ	64*2(SI),SI
	LEAQ	64*2(DI),DI
	SUBQ	$128,DX
	JMP	label2c

label2d:
	VMOVDQU	X2,16*2(BX)

	VZEROUPPER
	RET

