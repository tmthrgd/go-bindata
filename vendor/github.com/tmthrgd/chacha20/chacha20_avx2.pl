#!/usr/bin/env perl

##############################################################################
#                                                                            #
# Copyright 2014 Intel Corporation                                           #
#                                                                            #
# Licensed under the Apache License, Version 2.0 (the "License");            #
# you may not use this file except in compliance with the License.           #
# You may obtain a copy of the License at                                    #
#                                                                            #
#    http://www.apache.org/licenses/LICENSE-2.0                              #
#                                                                            #
# Unless required by applicable law or agreed to in writing, software        #
# distributed under the License is distributed on an "AS IS" BASIS,          #
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.   #
# See the License for the specific language governing permissions and        #
# limitations under the License.                                             #
#                                                                            #
##############################################################################
#                                                                            #
#  Developers and authors:                                                   #
#  Shay Gueron (1, 2), and Vlad Krasnov (1)                                  #
#  (1) Intel Corporation, Israel Development Center                          #
#  (2) University of Haifa                                                   #
#                                                                            #
# Related work:                                                              #
# M. Goll, S. Gueron, "Vectorization on ChaCha Stream Cipher", IEEE          #
#          Proceedings of 11th International Conference on Information       #
#          Technology: New Generations (ITNG 2014), 612-615 (2014).          #
# M. Goll, S. Gueron, "Vectorization on Poly1305 Message Authentication Code"#
#           to be published.                                                 #
# A. Langley, chacha20poly1305 for the AEAD head                             #
# https://git.openssl.org/gitweb/?p=openssl.git;a=commit;h=9a8646510b3d0a48e950748f7a2aaa12ed40d5e0  #
##############################################################################

$flavour = shift;
$output  = shift;
if ($flavour =~ /\./) { $output = $flavour; undef $flavour; }

$win64=0; $win64=1 if ($flavour =~ /[nm]asm|mingw64/ || $output =~ /\.asm$/);

$0 =~ m/(.*[\/\\])[^\/\\]+$/; $dir=$1;
( $xlate="${dir}x86_64-xlate.pl" and -f $xlate ) or
( $xlate="${dir}../../perlasm/x86_64-xlate.pl" and -f $xlate) or
die "can't locate x86_64-xlate.pl";

open OUT,"| \"$^X\" $xlate $flavour $output";
*STDOUT=*OUT;

if (`$ENV{CC} -Wa,-v -c -o /dev/null -x assembler /dev/null 2>&1`
    =~ /GNU assembler version ([2-9]\.[0-9]+)/) {
  $avx = ($1>=2.19) + ($1>=2.22);
}

if ($win64 && ($flavour =~ /nasm/ || $ENV{ASM} =~ /nasm/) &&
      `nasm -v 2>&1` =~ /NASM version ([2-9]\.[0-9]+)/) {
  $avx = ($1>=2.09) + ($1>=2.10);
}

if ($win64 && ($flavour =~ /masm/ || $ENV{ASM} =~ /ml64/) &&
      `ml64 2>&1` =~ /Version ([0-9]+)\./) {
  $avx = ($1>=10) + ($1>=11);
}

if (`$ENV{CC} -v 2>&1` =~ /(^clang version|based on LLVM) ([3-9])\.([0-9]+)/) {
  my $ver = $2 + $3/100.0;  # 3.1->3.01, 3.10->3.10
  $avx = ($ver>=3.0) + ($ver>=3.01);
}

$avx = 2 if ($flavour =~ /^golang/);

if ($avx>=2) {{

my ($state_4567, $state_89ab, $state_cdef, $tmp,
    $v0, $v1, $v2, $v3, $v4, $v5, $v6, $v7,
    $v8, $v9, $v10, $v11)=map("%ymm$_",(0..15));

sub chacha_qr {

my ($a,$b,$c,$d)=@_;

$code.=<<___;

  vpaddd  $b, $a, $a            # a += b
  vpxor   $a, $d, $d            # d ^= a
  vpshufb .rol16(%rip), $d, $d  # d <<<= 16

  vpaddd  $d, $c, $c            # c += d
  vpxor   $c, $b, $b            # b ^= c
  vpslld  \$12, $b, $tmp
  vpsrld  \$20, $b, $b
  vpxor   $tmp, $b, $b          # b <<<= 12

  vpaddd  $b, $a, $a            # a += b
  vpxor   $a, $d, $d            # d ^= a
  vpshufb .rol8(%rip), $d, $d   # d <<<= 8

  vpaddd  $d, $c, $c            # c += d
  vpxor   $c, $b, $b            # b ^= c

  vpslld  \$7, $b, $tmp
  vpsrld  \$25, $b, $b
  vpxor   $tmp, $b, $b          # b <<<= 7
___
}

if ($flavour =~ /^golang/) {
    $code.=<<___;
// Created by chacha20_avx2.pl - DO NOT EDIT
// perl chacha20_avx2.pl golang-no-avx chacha20_avx2_amd64.s

// +build amd64,!gccgo,!appengine

#include "textflag.h"

DATA chacha20_consts<>+0x00(SB)/1, \$"e"
DATA chacha20_consts<>+0x01(SB)/1, \$"x"
DATA chacha20_consts<>+0x02(SB)/1, \$"p"
DATA chacha20_consts<>+0x03(SB)/1, \$"a"
DATA chacha20_consts<>+0x04(SB)/1, \$"n"
DATA chacha20_consts<>+0x05(SB)/1, \$"d"
DATA chacha20_consts<>+0x06(SB)/1, \$" "
DATA chacha20_consts<>+0x07(SB)/1, \$"3"
DATA chacha20_consts<>+0x08(SB)/1, \$"2"
DATA chacha20_consts<>+0x09(SB)/1, \$"-"
DATA chacha20_consts<>+0x0a(SB)/1, \$"b"
DATA chacha20_consts<>+0x0b(SB)/1, \$"y"
DATA chacha20_consts<>+0x0c(SB)/1, \$"t"
DATA chacha20_consts<>+0x0d(SB)/1, \$"e"
DATA chacha20_consts<>+0x0e(SB)/1, \$" "
DATA chacha20_consts<>+0x0f(SB)/1, \$"k"
DATA chacha20_consts<>+0x10(SB)/1, \$"e"
DATA chacha20_consts<>+0x11(SB)/1, \$"x"
DATA chacha20_consts<>+0x12(SB)/1, \$"p"
DATA chacha20_consts<>+0x13(SB)/1, \$"a"
DATA chacha20_consts<>+0x14(SB)/1, \$"n"
DATA chacha20_consts<>+0x15(SB)/1, \$"d"
DATA chacha20_consts<>+0x16(SB)/1, \$" "
DATA chacha20_consts<>+0x17(SB)/1, \$"3"
DATA chacha20_consts<>+0x18(SB)/1, \$"2"
DATA chacha20_consts<>+0x19(SB)/1, \$"-"
DATA chacha20_consts<>+0x1a(SB)/1, \$"b"
DATA chacha20_consts<>+0x1b(SB)/1, \$"y"
DATA chacha20_consts<>+0x1c(SB)/1, \$"t"
DATA chacha20_consts<>+0x1d(SB)/1, \$"e"
DATA chacha20_consts<>+0x1e(SB)/1, \$" "
DATA chacha20_consts<>+0x1f(SB)/1, \$"k"
GLOBL chacha20_consts<>(SB), RODATA, \$32

DATA rol8<>+0x00(SB)/1, \$3
DATA rol8<>+0x01(SB)/1, \$0
DATA rol8<>+0x02(SB)/1, \$1
DATA rol8<>+0x03(SB)/1, \$2
DATA rol8<>+0x04(SB)/1, \$7
DATA rol8<>+0x05(SB)/1, \$4
DATA rol8<>+0x06(SB)/1, \$5
DATA rol8<>+0x07(SB)/1, \$6
DATA rol8<>+0x08(SB)/1, \$11
DATA rol8<>+0x09(SB)/1, \$8
DATA rol8<>+0x0a(SB)/1, \$9
DATA rol8<>+0x0b(SB)/1, \$10
DATA rol8<>+0x0c(SB)/1, \$15
DATA rol8<>+0x0d(SB)/1, \$12
DATA rol8<>+0x0e(SB)/1, \$13
DATA rol8<>+0x0f(SB)/1, \$14
DATA rol8<>+0x10(SB)/1, \$3
DATA rol8<>+0x11(SB)/1, \$0
DATA rol8<>+0x12(SB)/1, \$1
DATA rol8<>+0x13(SB)/1, \$2
DATA rol8<>+0x14(SB)/1, \$7
DATA rol8<>+0x15(SB)/1, \$4
DATA rol8<>+0x16(SB)/1, \$5
DATA rol8<>+0x17(SB)/1, \$6
DATA rol8<>+0x18(SB)/1, \$11
DATA rol8<>+0x19(SB)/1, \$8
DATA rol8<>+0x1a(SB)/1, \$9
DATA rol8<>+0x1b(SB)/1, \$10
DATA rol8<>+0x1c(SB)/1, \$15
DATA rol8<>+0x1d(SB)/1, \$12
DATA rol8<>+0x1e(SB)/1, \$13
DATA rol8<>+0x1f(SB)/1, \$14
GLOBL rol8<>(SB), RODATA, \$32

DATA rol16<>+0x00(SB)/1, \$2
DATA rol16<>+0x01(SB)/1, \$3
DATA rol16<>+0x02(SB)/1, \$0
DATA rol16<>+0x03(SB)/1, \$1
DATA rol16<>+0x04(SB)/1, \$6
DATA rol16<>+0x05(SB)/1, \$7
DATA rol16<>+0x06(SB)/1, \$4
DATA rol16<>+0x07(SB)/1, \$5
DATA rol16<>+0x08(SB)/1, \$10
DATA rol16<>+0x09(SB)/1, \$11
DATA rol16<>+0x0a(SB)/1, \$8
DATA rol16<>+0x0b(SB)/1, \$9
DATA rol16<>+0x0c(SB)/1, \$14
DATA rol16<>+0x0d(SB)/1, \$15
DATA rol16<>+0x0e(SB)/1, \$12
DATA rol16<>+0x0f(SB)/1, \$13
DATA rol16<>+0x10(SB)/1, \$2
DATA rol16<>+0x11(SB)/1, \$3
DATA rol16<>+0x12(SB)/1, \$0
DATA rol16<>+0x13(SB)/1, \$1
DATA rol16<>+0x14(SB)/1, \$6
DATA rol16<>+0x15(SB)/1, \$7
DATA rol16<>+0x16(SB)/1, \$4
DATA rol16<>+0x17(SB)/1, \$5
DATA rol16<>+0x18(SB)/1, \$10
DATA rol16<>+0x19(SB)/1, \$11
DATA rol16<>+0x1a(SB)/1, \$8
DATA rol16<>+0x1b(SB)/1, \$9
DATA rol16<>+0x1c(SB)/1, \$14
DATA rol16<>+0x1d(SB)/1, \$15
DATA rol16<>+0x1e(SB)/1, \$12
DATA rol16<>+0x1f(SB)/1, \$13
GLOBL rol16<>(SB), RODATA, \$32

DATA avx2Init<>+0x00(SB)/8, \$0x0
DATA avx2Init<>+0x08(SB)/8, \$0x0
DATA avx2Init<>+0x10(SB)/8, \$0x1
DATA avx2Init<>+0x18(SB)/8, \$0x0
GLOBL avx2Init<>(SB), RODATA, \$32

DATA avx2Inc<>+0x00(SB)/8, \$0x2
DATA avx2Inc<>+0x08(SB)/8, \$0x0
DATA avx2Inc<>+0x10(SB)/8, \$0x2
DATA avx2Inc<>+0x18(SB)/8, \$0x0
GLOBL avx2Inc<>(SB), RODATA, \$32

___
} else {
    $code.=<<___;
.text
.align 32
chacha20_consts:
.byte 'e','x','p','a','n','d',' ','3','2','-','b','y','t','e',' ','k'
.byte 'e','x','p','a','n','d',' ','3','2','-','b','y','t','e',' ','k'
.rol8:
.byte 3,0,1,2, 7,4,5,6, 11,8,9,10, 15,12,13,14
.byte 3,0,1,2, 7,4,5,6, 11,8,9,10, 15,12,13,14
.rol16:
.byte 2,3,0,1, 6,7,4,5, 10,11,8,9, 14,15,12,13
.byte 2,3,0,1, 6,7,4,5, 10,11,8,9, 14,15,12,13
.avx2Init:
.quad 0,0,1,0
.avx2Inc:
.quad 2,0,2,0
___
}

{

my $state_cdef_xmm=$state_cdef;

substr($state_cdef_xmm, 1, 1, "x");

my ($out, $in, $in_len, $key_ptr, $nr)
   =("%rdi", "%rsi", "%rdx", "%rbx", "%r8");

if ($flavour =~ /^golang/) {
    $code.=<<___;
TEXT ·chacha_20_core_avx2(SB),\$0-32
	movq	out+0(FP), DI
	movq	in+8(FP), SI
	movq	in_len+16(FP), DX
	movq	state+24(FP), BX

	movq	\$chacha20_consts<>(SB), R11
	movq	\$rol8<>(SB), R12
	movq	\$rol16<>(SB), R13
	movq	\$avx2Init<>(SB), R14
	movq	\$avx2Inc<>(SB), R15

___
} else {
    $code.=<<___;
.globl chacha_20_core_avx2
.type  chacha_20_core_avx2 ,\@function,2
.align 64
chacha_20_core_avx2:
___
}

$code.=<<___;
  vzeroupper

  # Init state
  vbroadcasti128  16*0($key_ptr), $state_4567
  vbroadcasti128  16*1($key_ptr), $state_89ab
  vbroadcasti128  16*2($key_ptr), $state_cdef
  vpaddq    .avx2Init(%rip), $state_cdef, $state_cdef

2:
  cmp  \$6*64, $in_len
  jb  2f

  vmovdqa  chacha20_consts(%rip), $v0
  vmovdqa  chacha20_consts(%rip), $v4
  vmovdqa  chacha20_consts(%rip), $v8

  vmovdqa  $state_4567, $v1
  vmovdqa  $state_4567, $v5
  vmovdqa  $state_4567, $v9

  vmovdqa  $state_89ab, $v2
  vmovdqa  $state_89ab, $v6
  vmovdqa  $state_89ab, $v10

  vmovdqa  $state_cdef, $v3
  vpaddq  .avx2Inc(%rip), $v3, $v7
  vpaddq  .avx2Inc(%rip), $v7, $v11

  mov  \$10, $nr

  1:
___

    &chacha_qr( $v0, $v1, $v2, $v3);
    &chacha_qr( $v4, $v5, $v6, $v7);
    &chacha_qr( $v8, $v9,$v10,$v11);

$code.=<<___;
    vpalignr  \$4,  $v1,  $v1,  $v1
    vpalignr  \$8,  $v2,  $v2,  $v2
    vpalignr \$12,  $v3,  $v3,  $v3
    vpalignr  \$4,  $v5,  $v5,  $v5
    vpalignr  \$8,  $v6,  $v6,  $v6
    vpalignr \$12,  $v7,  $v7,  $v7
    vpalignr  \$4,  $v9,  $v9,  $v9
    vpalignr  \$8, $v10, $v10, $v10
    vpalignr \$12, $v11, $v11, $v11
___

    &chacha_qr( $v0, $v1, $v2, $v3);
    &chacha_qr( $v4, $v5, $v6, $v7);
    &chacha_qr( $v8, $v9,$v10,$v11);

$code.=<<___;
    vpalignr \$12,  $v1,  $v1,  $v1
    vpalignr  \$8,  $v2,  $v2,  $v2
    vpalignr  \$4,  $v3,  $v3,  $v3
    vpalignr \$12,  $v5,  $v5,  $v5
    vpalignr  \$8,  $v6,  $v6,  $v6
    vpalignr  \$4,  $v7,  $v7,  $v7
    vpalignr \$12,  $v9,  $v9,  $v9
    vpalignr  \$8, $v10, $v10, $v10
    vpalignr  \$4, $v11, $v11, $v11

    dec  $nr

  jnz  1b

  vpaddd  chacha20_consts(%rip), $v0, $v0
  vpaddd  chacha20_consts(%rip), $v4, $v4
  vpaddd  chacha20_consts(%rip), $v8, $v8

  vpaddd  $state_4567, $v1, $v1
  vpaddd  $state_4567, $v5, $v5
  vpaddd  $state_4567, $v9, $v9

  vpaddd  $state_89ab, $v2, $v2
  vpaddd  $state_89ab, $v6, $v6
  vpaddd  $state_89ab, $v10, $v10

  vpaddd  $state_cdef, $v3, $v3
  vpaddq  .avx2Inc(%rip), $state_cdef, $state_cdef
  vpaddd  $state_cdef, $v7, $v7
  vpaddq  .avx2Inc(%rip), $state_cdef, $state_cdef
  vpaddd  $state_cdef, $v11, $v11
  vpaddq  .avx2Inc(%rip), $state_cdef, $state_cdef

  vperm2i128  \$0x02, $v0, $v1, $tmp
  vpxor  32*0($in), $tmp, $tmp
  vmovdqu  $tmp, 32*0($out)
  vperm2i128  \$0x02, $v2, $v3, $tmp
  vpxor  32*1($in), $tmp, $tmp
  vmovdqu  $tmp, 32*1($out)
  vperm2i128  \$0x13, $v0, $v1, $tmp
  vpxor  32*2($in), $tmp, $tmp
  vmovdqu  $tmp, 32*2($out)
  vperm2i128  \$0x13, $v2, $v3, $tmp
  vpxor  32*3($in), $tmp, $tmp
  vmovdqu  $tmp, 32*3($out)

  vperm2i128  \$0x02, $v4, $v5, $v0
  vperm2i128  \$0x02, $v6, $v7, $v1
  vperm2i128  \$0x13, $v4, $v5, $v2
  vperm2i128  \$0x13, $v6, $v7, $v3

  vpxor  32*4($in), $v0, $v0
  vpxor  32*5($in), $v1, $v1
  vpxor  32*6($in), $v2, $v2
  vpxor  32*7($in), $v3, $v3

  vmovdqu  $v0, 32*4($out)
  vmovdqu  $v1, 32*5($out)
  vmovdqu  $v2, 32*6($out)
  vmovdqu  $v3, 32*7($out)

  vperm2i128  \$0x02, $v8, $v9, $v0
  vperm2i128  \$0x02, $v10, $v11, $v1
  vperm2i128  \$0x13, $v8, $v9, $v2
  vperm2i128  \$0x13, $v10, $v11, $v3

  vpxor  32*8($in), $v0, $v0
  vpxor  32*9($in), $v1, $v1
  vpxor  32*10($in), $v2, $v2
  vpxor  32*11($in), $v3, $v3

  vmovdqu  $v0, 32*8($out)
  vmovdqu  $v1, 32*9($out)
  vmovdqu  $v2, 32*10($out)
  vmovdqu  $v3, 32*11($out)

  lea  64*6($in), $in
  lea  64*6($out), $out
  sub  \$64*6, $in_len

  jmp  2b

2:
  cmp  \$4*64, $in_len
  jb  2f

  vmovdqa  chacha20_consts(%rip), $v0
  vmovdqa  chacha20_consts(%rip), $v4
  vmovdqa  $state_4567, $v1
  vmovdqa  $state_4567, $v5
  vmovdqa  $state_89ab, $v2
  vmovdqa  $state_89ab, $v6
  vmovdqa  $state_89ab, $v10
  vmovdqa  $state_cdef, $v3
  vpaddq   .avx2Inc(%rip), $v3, $v7

  mov  \$10, $nr

  1:
___

    &chacha_qr($v0,$v1,$v2,$v3);
    &chacha_qr($v4,$v5,$v6,$v7);

$code.=<<___;
    vpalignr  \$4, $v1, $v1, $v1
    vpalignr  \$8, $v2, $v2, $v2
    vpalignr \$12, $v3, $v3, $v3
    vpalignr  \$4, $v5, $v5, $v5
    vpalignr  \$8, $v6, $v6, $v6
    vpalignr \$12, $v7, $v7, $v7
___

    &chacha_qr($v0,$v1,$v2,$v3);
    &chacha_qr($v4,$v5,$v6,$v7);

$code.=<<___;
    vpalignr \$12, $v1, $v1, $v1
    vpalignr  \$8, $v2, $v2, $v2
    vpalignr  \$4, $v3, $v3, $v3
    vpalignr \$12, $v5, $v5, $v5
    vpalignr  \$8, $v6, $v6, $v6
    vpalignr  \$4, $v7, $v7, $v7

    dec  $nr

  jnz  1b

  vpaddd  chacha20_consts(%rip), $v0, $v0
  vpaddd  chacha20_consts(%rip), $v4, $v4

  vpaddd  $state_4567, $v1, $v1
  vpaddd  $state_4567, $v5, $v5

  vpaddd  $state_89ab, $v2, $v2
  vpaddd  $state_89ab, $v6, $v6

  vpaddd  $state_cdef, $v3, $v3
  vpaddq  .avx2Inc(%rip), $state_cdef, $state_cdef
  vpaddd  $state_cdef, $v7, $v7
  vpaddq  .avx2Inc(%rip), $state_cdef, $state_cdef

  vperm2i128  \$0x02, $v0, $v1, $v8
  vperm2i128  \$0x02, $v2, $v3, $v9
  vperm2i128  \$0x13, $v0, $v1, $v10
  vperm2i128  \$0x13, $v2, $v3, $v11

  vpxor  32*0($in), $v8, $v8
  vpxor  32*1($in), $v9, $v9
  vpxor  32*2($in), $v10, $v10
  vpxor  32*3($in), $v11, $v11

  vmovdqu  $v8, 32*0($out)
  vmovdqu  $v9, 32*1($out)
  vmovdqu  $v10, 32*2($out)
  vmovdqu  $v11, 32*3($out)

  vperm2i128  \$0x02, $v4, $v5, $v0
  vperm2i128  \$0x02, $v6, $v7, $v1
  vperm2i128  \$0x13, $v4, $v5, $v2
  vperm2i128  \$0x13, $v6, $v7, $v3

  vpxor  32*4($in), $v0, $v0
  vpxor  32*5($in), $v1, $v1
  vpxor  32*6($in), $v2, $v2
  vpxor  32*7($in), $v3, $v3

  vmovdqu  $v0, 32*4($out)
  vmovdqu  $v1, 32*5($out)
  vmovdqu  $v2, 32*6($out)
  vmovdqu  $v3, 32*7($out)

  lea  64*4($in), $in
  lea  64*4($out), $out
  sub  \$64*4, $in_len

  jmp  2b
2:
  cmp  \$128, $in_len
  jb  2f

  vmovdqa  chacha20_consts(%rip), $v0
  vmovdqa  $state_4567, $v1
  vmovdqa  $state_89ab, $v2
  vmovdqa  $state_cdef, $v3

  mov  \$10, $nr

  1:
___

    &chacha_qr($v0,$v1,$v2,$v3);

$code.=<<___;
    vpalignr   \$4, $v1, $v1, $v1
    vpalignr   \$8, $v2, $v2, $v2
    vpalignr  \$12, $v3, $v3, $v3
___

    &chacha_qr($v0,$v1,$v2,$v3);
$code.=<<___;
    vpalignr  \$12, $v1, $v1, $v1
    vpalignr   \$8, $v2, $v2, $v2
    vpalignr   \$4, $v3, $v3, $v3

    dec  $nr
  jnz  1b

  vpaddd  chacha20_consts(%rip), $v0, $v0
  vpaddd  $state_4567, $v1, $v1
  vpaddd  $state_89ab, $v2, $v2
  vpaddd  $state_cdef, $v3, $v3
  vpaddq  .avx2Inc(%rip), $state_cdef, $state_cdef

  vperm2i128  \$0x02, $v0, $v1, $v8
  vperm2i128  \$0x02, $v2, $v3, $v9
  vperm2i128  \$0x13, $v0, $v1, $v10
  vperm2i128  \$0x13, $v2, $v3, $v11

  vpxor  32*0($in), $v8, $v8
  vpxor  32*1($in), $v9, $v9
  vpxor  32*2($in), $v10, $v10
  vpxor  32*3($in), $v11, $v11

  vmovdqu  $v8, 32*0($out)
  vmovdqu  $v9, 32*1($out)
  vmovdqu  $v10, 32*2($out)
  vmovdqu  $v11, 32*3($out)

  lea  64*2($in), $in
  lea  64*2($out), $out
  sub  \$64*2, $in_len
  jmp  2b

2:
  vmovdqu  $state_cdef_xmm, 16*2($key_ptr)

  vzeroupper
  ret
.size  chacha_20_core_avx2,.-chacha_20_core_avx2
___
}
}}


$code =~ s/\`([^\`]*)\`/eval($1)/gem;

if ($flavour =~ /^golang/) {
	$code =~ s/.chacha20_consts\(%rip\)/(%r11)/g;
	$code =~ s/.rol8\(%rip\)/(%r12)/g;
	$code =~ s/.rol16\(%rip\)/(%r13)/g;
	$code =~ s/.avx2Init\(%rip\)/(%r14)/g;
	$code =~ s/.avx2Inc\(%rip\)/(%r15)/g;
}

print $code;

close STDOUT;

