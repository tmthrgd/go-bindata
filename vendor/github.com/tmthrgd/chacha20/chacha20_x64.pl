#!/usr/bin/env perl

##############################################################################
#                                                                            #
# Public Domain                                                              #
#                                                                            #
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

if ($flavour =~ /^golang/) {
    $code.=<<___;
// Created by chacha20_x64.pl - DO NOT EDIT
// perl chacha20_x64.pl golang-no-avx chacha20_x64_amd64.s

// +build amd64,!gccgo,!appengine

// This code was translated into a form compatible with 6a from the public
// domain sources in SUPERCOP: http://bench.cr.yp.to/supercop.html

#include "textflag.h"

___
}

{

if ($flavour =~ /^golang/) {
    $code.=<<___;
TEXT ·chacha_20_core_x64(SB),\$`512+64`-32
	movq	out+0(FP), DX
	movq	in+8(FP), SI
	movq	in_len+16(FP), BX
	movq	state+24(FP), DI

	movq	\$state-512(SP), R12
	andq	\$~63, %r12

___
} else {
    $code.=<<___;
.globl chacha_20_core_x64
.type  chacha_20_core_x64 ,\@function,2
.align 64
chacha_20_core_x64:
pushq %r13
pushq %r12
movq %rsp, %r12
andq \$~63, %r12
subq \$512, %r12

___
}

$code.=<<___;
movq \$0x3320646e61707865, %r8
movq \$0x6b20657479622d32, %r9
movd %r8, %xmm8
movd %r9, %xmm14
punpcklqdq %xmm14, %xmm8
movdqu 0(%rdi), %xmm9
movdqu 16(%rdi), %xmm10
movdqu 32(%rdi), %xmm11
movq \$20, %r11
movq \$1, %r9
movdqa %xmm8, 0(%r12)
movdqa %xmm9, 16(%r12)
movdqa %xmm10, 32(%r12)
movdqa %xmm11, 48(%r12)
cmpq \$256, %rbx
jb chacha_blocks_sse2_below256
pshufd \$0x00, %xmm8, %xmm0
pshufd \$0x55, %xmm8, %xmm1
pshufd \$0xaa, %xmm8, %xmm2
pshufd \$0xff, %xmm8, %xmm3
movdqa %xmm0, 128(%r12)
movdqa %xmm1, 144(%r12)
movdqa %xmm2, 160(%r12)
movdqa %xmm3, 176(%r12)
pshufd \$0x00, %xmm9, %xmm0
pshufd \$0x55, %xmm9, %xmm1
pshufd \$0xaa, %xmm9, %xmm2
pshufd \$0xff, %xmm9, %xmm3
movdqa %xmm0, 192(%r12)
movdqa %xmm1, 208(%r12)
movdqa %xmm2, 224(%r12)
movdqa %xmm3, 240(%r12)
pshufd \$0x00, %xmm10, %xmm0
pshufd \$0x55, %xmm10, %xmm1
pshufd \$0xaa, %xmm10, %xmm2
pshufd \$0xff, %xmm10, %xmm3
movdqa %xmm0, 256(%r12)
movdqa %xmm1, 272(%r12)
movdqa %xmm2, 288(%r12)
movdqa %xmm3, 304(%r12)
pshufd \$0xaa, %xmm11, %xmm0
pshufd \$0xff, %xmm11, %xmm1
movdqa %xmm0, 352(%r12)
movdqa %xmm1, 368(%r12)
.p2align 6,,63
chacha_blocks_sse2_atleast256:
movq 48(%r12), %rax
leaq 1(%rax), %r8
leaq 2(%rax), %r9
leaq 3(%rax), %r10
leaq 4(%rax), %r13
movl %eax, 320(%r12)
movl %r8d, 4+320(%r12)
movl %r9d, 8+320(%r12)
movl %r10d, 12+320(%r12)
shrq \$32, %rax
shrq \$32, %r8
shrq \$32, %r9
shrq \$32, %r10
movl %eax, 336(%r12)
movl %r8d, 4+336(%r12)
movl %r9d, 8+336(%r12)
movl %r10d, 12+336(%r12)
movq %r13, 48(%r12)
movq \$20, %r11
movdqa 128(%r12), %xmm0
movdqa 144(%r12), %xmm1
movdqa 160(%r12), %xmm2
movdqa 176(%r12), %xmm3
movdqa 192(%r12), %xmm4
movdqa 208(%r12), %xmm5
movdqa 224(%r12), %xmm6
movdqa 240(%r12), %xmm7
movdqa 256(%r12), %xmm8
movdqa 272(%r12), %xmm9
movdqa 288(%r12), %xmm10
movdqa 304(%r12), %xmm11
movdqa 320(%r12), %xmm12
movdqa 336(%r12), %xmm13
movdqa 352(%r12), %xmm14
movdqa 368(%r12), %xmm15
chacha_blocks_sse2_mainloop1:
paddd %xmm4, %xmm0
paddd %xmm5, %xmm1
pxor %xmm0, %xmm12
pxor %xmm1, %xmm13
paddd %xmm6, %xmm2
paddd %xmm7, %xmm3
movdqa %xmm6, 96(%r12)
pxor %xmm2, %xmm14
pxor %xmm3, %xmm15
pshuflw \$0xb1,%xmm12,%xmm12
pshufhw \$0xb1,%xmm12,%xmm12
pshuflw \$0xb1,%xmm13,%xmm13
pshufhw \$0xb1,%xmm13,%xmm13
pshuflw \$0xb1,%xmm14,%xmm14
pshufhw \$0xb1,%xmm14,%xmm14
pshuflw \$0xb1,%xmm15,%xmm15
pshufhw \$0xb1,%xmm15,%xmm15
paddd %xmm12, %xmm8
paddd %xmm13, %xmm9
paddd %xmm14, %xmm10
paddd %xmm15, %xmm11
movdqa %xmm12, 112(%r12)
pxor %xmm8, %xmm4
pxor %xmm9, %xmm5
movdqa 96(%r12), %xmm6
movdqa %xmm4, %xmm12
pslld \$12, %xmm4
psrld \$20, %xmm12
pxor %xmm12, %xmm4
movdqa %xmm5, %xmm12
pslld \$12, %xmm5
psrld \$20, %xmm12
pxor %xmm12, %xmm5
pxor %xmm10, %xmm6
pxor %xmm11, %xmm7
movdqa %xmm6, %xmm12
pslld \$12, %xmm6
psrld \$20, %xmm12
pxor %xmm12, %xmm6
movdqa %xmm7, %xmm12
pslld \$12, %xmm7
psrld \$20, %xmm12
pxor %xmm12, %xmm7
movdqa 112(%r12), %xmm12
paddd %xmm4, %xmm0
paddd %xmm5, %xmm1
pxor %xmm0, %xmm12
pxor %xmm1, %xmm13
paddd %xmm6, %xmm2
paddd %xmm7, %xmm3
movdqa %xmm6, 96(%r12)
pxor %xmm2, %xmm14
pxor %xmm3, %xmm15
movdqa %xmm12, %xmm6
pslld \$ 8, %xmm12
psrld \$24, %xmm6
pxor %xmm6, %xmm12
movdqa %xmm13, %xmm6
pslld \$ 8, %xmm13
psrld \$24, %xmm6
pxor %xmm6, %xmm13
paddd %xmm12, %xmm8
paddd %xmm13, %xmm9
movdqa %xmm14, %xmm6
pslld \$ 8, %xmm14
psrld \$24, %xmm6
pxor %xmm6, %xmm14
movdqa %xmm15, %xmm6
pslld \$ 8, %xmm15
psrld \$24, %xmm6
pxor %xmm6, %xmm15
paddd %xmm14, %xmm10
paddd %xmm15, %xmm11
movdqa %xmm12, 112(%r12)
pxor %xmm8, %xmm4
pxor %xmm9, %xmm5
movdqa 96(%r12), %xmm6
movdqa %xmm4, %xmm12
pslld \$ 7, %xmm4
psrld \$25, %xmm12
pxor %xmm12, %xmm4
movdqa %xmm5, %xmm12
pslld \$ 7, %xmm5
psrld \$25, %xmm12
pxor %xmm12, %xmm5
pxor %xmm10, %xmm6
pxor %xmm11, %xmm7
movdqa %xmm6, %xmm12
pslld \$ 7, %xmm6
psrld \$25, %xmm12
pxor %xmm12, %xmm6
movdqa %xmm7, %xmm12
pslld \$ 7, %xmm7
psrld \$25, %xmm12
pxor %xmm12, %xmm7
movdqa 112(%r12), %xmm12
paddd %xmm5, %xmm0
paddd %xmm6, %xmm1
pxor %xmm0, %xmm15
pxor %xmm1, %xmm12
paddd %xmm7, %xmm2
paddd %xmm4, %xmm3
movdqa %xmm7, 96(%r12)
pxor %xmm2, %xmm13
pxor %xmm3, %xmm14
pshuflw \$0xb1,%xmm15,%xmm15
pshufhw \$0xb1,%xmm15,%xmm15
pshuflw \$0xb1,%xmm12,%xmm12
pshufhw \$0xb1,%xmm12,%xmm12
pshuflw \$0xb1,%xmm13,%xmm13
pshufhw \$0xb1,%xmm13,%xmm13
pshuflw \$0xb1,%xmm14,%xmm14
pshufhw \$0xb1,%xmm14,%xmm14
paddd %xmm15, %xmm10
paddd %xmm12, %xmm11
paddd %xmm13, %xmm8
paddd %xmm14, %xmm9
movdqa %xmm15, 112(%r12)
pxor %xmm10, %xmm5
pxor %xmm11, %xmm6
movdqa 96(%r12), %xmm7
movdqa %xmm5, %xmm15
pslld \$ 12, %xmm5
psrld \$20, %xmm15
pxor %xmm15, %xmm5
movdqa %xmm6, %xmm15
pslld \$ 12, %xmm6
psrld \$20, %xmm15
pxor %xmm15, %xmm6
pxor %xmm8, %xmm7
pxor %xmm9, %xmm4
movdqa %xmm7, %xmm15
pslld \$ 12, %xmm7
psrld \$20, %xmm15
pxor %xmm15, %xmm7
movdqa %xmm4, %xmm15
pslld \$ 12, %xmm4
psrld \$20, %xmm15
pxor %xmm15, %xmm4
movdqa 112(%r12), %xmm15
paddd %xmm5, %xmm0
paddd %xmm6, %xmm1
pxor %xmm0, %xmm15
pxor %xmm1, %xmm12
paddd %xmm7, %xmm2
paddd %xmm4, %xmm3
movdqa %xmm7, 96(%r12)
pxor %xmm2, %xmm13
pxor %xmm3, %xmm14
movdqa %xmm15, %xmm7
pslld \$ 8, %xmm15
psrld \$24, %xmm7
pxor %xmm7, %xmm15
movdqa %xmm12, %xmm7
pslld \$ 8, %xmm12
psrld \$24, %xmm7
pxor %xmm7, %xmm12
paddd %xmm15, %xmm10
paddd %xmm12, %xmm11
movdqa %xmm13, %xmm7
pslld \$ 8, %xmm13
psrld \$24, %xmm7
pxor %xmm7, %xmm13
movdqa %xmm14, %xmm7
pslld \$ 8, %xmm14
psrld \$24, %xmm7
pxor %xmm7, %xmm14
paddd %xmm13, %xmm8
paddd %xmm14, %xmm9
movdqa %xmm15, 112(%r12)
pxor %xmm10, %xmm5
pxor %xmm11, %xmm6
movdqa 96(%r12), %xmm7
movdqa %xmm5, %xmm15
pslld \$ 7, %xmm5
psrld \$25, %xmm15
pxor %xmm15, %xmm5
movdqa %xmm6, %xmm15
pslld \$ 7, %xmm6
psrld \$25, %xmm15
pxor %xmm15, %xmm6
pxor %xmm8, %xmm7
pxor %xmm9, %xmm4
movdqa %xmm7, %xmm15
pslld \$ 7, %xmm7
psrld \$25, %xmm15
pxor %xmm15, %xmm7
movdqa %xmm4, %xmm15
pslld \$ 7, %xmm4
psrld \$25, %xmm15
pxor %xmm15, %xmm4
movdqa 112(%r12), %xmm15
subq \$2, %r11
jnz chacha_blocks_sse2_mainloop1
paddd 128(%r12), %xmm0
paddd 144(%r12), %xmm1
paddd 160(%r12), %xmm2
paddd 176(%r12), %xmm3
paddd 192(%r12), %xmm4
paddd 208(%r12), %xmm5
paddd 224(%r12), %xmm6
paddd 240(%r12), %xmm7
paddd 256(%r12), %xmm8
paddd 272(%r12), %xmm9
paddd 288(%r12), %xmm10
paddd 304(%r12), %xmm11
paddd 320(%r12), %xmm12
paddd 336(%r12), %xmm13
paddd 352(%r12), %xmm14
paddd 368(%r12), %xmm15
movdqa %xmm8, 384(%r12)
movdqa %xmm9, 400(%r12)
movdqa %xmm10, 416(%r12)
movdqa %xmm11, 432(%r12)
movdqa %xmm12, 448(%r12)
movdqa %xmm13, 464(%r12)
movdqa %xmm14, 480(%r12)
movdqa %xmm15, 496(%r12)
movdqa %xmm0, %xmm8
movdqa %xmm2, %xmm9
movdqa %xmm4, %xmm10
movdqa %xmm6, %xmm11
punpckhdq %xmm1, %xmm0
punpckhdq %xmm3, %xmm2
punpckhdq %xmm5, %xmm4
punpckhdq %xmm7, %xmm6
punpckldq %xmm1, %xmm8
punpckldq %xmm3, %xmm9
punpckldq %xmm5, %xmm10
punpckldq %xmm7, %xmm11
movdqa %xmm0, %xmm1
movdqa %xmm4, %xmm3
movdqa %xmm8, %xmm5
movdqa %xmm10, %xmm7
punpckhqdq %xmm2, %xmm0
punpckhqdq %xmm6, %xmm4
punpckhqdq %xmm9, %xmm8
punpckhqdq %xmm11, %xmm10
punpcklqdq %xmm2, %xmm1
punpcklqdq %xmm6, %xmm3
punpcklqdq %xmm9, %xmm5
punpcklqdq %xmm11, %xmm7
movdqu 0(%rsi), %xmm2
movdqu 16(%rsi), %xmm6
movdqu 64(%rsi), %xmm9
movdqu 80(%rsi), %xmm11
movdqu 128(%rsi), %xmm12
movdqu 144(%rsi), %xmm13
movdqu 192(%rsi), %xmm14
movdqu 208(%rsi), %xmm15
pxor %xmm2, %xmm5
pxor %xmm6, %xmm7
pxor %xmm9, %xmm8
pxor %xmm11, %xmm10
pxor %xmm12, %xmm1
pxor %xmm13, %xmm3
pxor %xmm14, %xmm0
pxor %xmm15, %xmm4
movdqu %xmm5, 0(%rdx)
movdqu %xmm7, 16(%rdx)
movdqu %xmm8, 64(%rdx)
movdqu %xmm10, 80(%rdx)
movdqu %xmm1, 128(%rdx)
movdqu %xmm3, 144(%rdx)
movdqu %xmm0, 192(%rdx)
movdqu %xmm4, 208(%rdx)
movdqa 384(%r12), %xmm0
movdqa 400(%r12), %xmm1
movdqa 416(%r12), %xmm2
movdqa 432(%r12), %xmm3
movdqa 448(%r12), %xmm4
movdqa 464(%r12), %xmm5
movdqa 480(%r12), %xmm6
movdqa 496(%r12), %xmm7
movdqa %xmm0, %xmm8
movdqa %xmm2, %xmm9
movdqa %xmm4, %xmm10
movdqa %xmm6, %xmm11
punpckldq %xmm1, %xmm8
punpckldq %xmm3, %xmm9
punpckhdq %xmm1, %xmm0
punpckhdq %xmm3, %xmm2
punpckldq %xmm5, %xmm10
punpckldq %xmm7, %xmm11
punpckhdq %xmm5, %xmm4
punpckhdq %xmm7, %xmm6
movdqa %xmm8, %xmm1
movdqa %xmm0, %xmm3
movdqa %xmm10, %xmm5
movdqa %xmm4, %xmm7
punpcklqdq %xmm9, %xmm1
punpcklqdq %xmm11, %xmm5
punpckhqdq %xmm9, %xmm8
punpckhqdq %xmm11, %xmm10
punpcklqdq %xmm2, %xmm3
punpcklqdq %xmm6, %xmm7
punpckhqdq %xmm2, %xmm0
punpckhqdq %xmm6, %xmm4
movdqu 32(%rsi), %xmm2
movdqu 48(%rsi), %xmm6
movdqu 96(%rsi), %xmm9
movdqu 112(%rsi), %xmm11
movdqu 160(%rsi), %xmm12
movdqu 176(%rsi), %xmm13
movdqu 224(%rsi), %xmm14
movdqu 240(%rsi), %xmm15
pxor %xmm2, %xmm1
pxor %xmm6, %xmm5
pxor %xmm9, %xmm8
pxor %xmm11, %xmm10
pxor %xmm12, %xmm3
pxor %xmm13, %xmm7
pxor %xmm14, %xmm0
pxor %xmm15, %xmm4
movdqu %xmm1, 32(%rdx)
movdqu %xmm5, 48(%rdx)
movdqu %xmm8, 96(%rdx)
movdqu %xmm10, 112(%rdx)
movdqu %xmm3, 160(%rdx)
movdqu %xmm7, 176(%rdx)
movdqu %xmm0, 224(%rdx)
movdqu %xmm4, 240(%rdx)
addq \$256, %rsi
addq \$256, %rdx
subq \$256, %rbx
cmp \$256, %rbx
jae chacha_blocks_sse2_atleast256
movdqa 0(%r12), %xmm8
movdqa 16(%r12), %xmm9
movdqa 32(%r12), %xmm10
movdqa 48(%r12), %xmm11
movq \$1, %r9
chacha_blocks_sse2_below256:
movq %r9, %xmm5
andq %rbx, %rbx
jz chacha_blocks_sse2_done
cmpq \$64, %rbx
jb chacha_blocks_sse2_done
chacha_blocks_sse2_above63:
movdqa %xmm8, %xmm0
movdqa %xmm9, %xmm1
movdqa %xmm10, %xmm2
movdqa %xmm11, %xmm3
movq \$20, %r11
chacha_blocks_sse2_mainloop2:
paddd %xmm1, %xmm0
pxor %xmm0, %xmm3
pshuflw \$0xb1,%xmm3,%xmm3
pshufhw \$0xb1,%xmm3,%xmm3
paddd %xmm3, %xmm2
pxor %xmm2, %xmm1
movdqa %xmm1,%xmm4
pslld \$12, %xmm1
psrld \$20, %xmm4
pxor %xmm4, %xmm1
paddd %xmm1, %xmm0
pxor %xmm0, %xmm3
movdqa %xmm3,%xmm4
pslld \$8, %xmm3
psrld \$24, %xmm4
pshufd \$0x93,%xmm0,%xmm0
pxor %xmm4, %xmm3
paddd %xmm3, %xmm2
pshufd \$0x4e,%xmm3,%xmm3
pxor %xmm2, %xmm1
pshufd \$0x39,%xmm2,%xmm2
movdqa %xmm1,%xmm4
pslld \$7, %xmm1
psrld \$25, %xmm4
pxor %xmm4, %xmm1
subq \$2, %r11
paddd %xmm1, %xmm0
pxor %xmm0, %xmm3
pshuflw \$0xb1,%xmm3,%xmm3
pshufhw \$0xb1,%xmm3,%xmm3
paddd %xmm3, %xmm2
pxor %xmm2, %xmm1
movdqa %xmm1,%xmm4
pslld \$12, %xmm1
psrld \$20, %xmm4
pxor %xmm4, %xmm1
paddd %xmm1, %xmm0
pxor %xmm0, %xmm3
movdqa %xmm3,%xmm4
pslld \$8, %xmm3
psrld \$24, %xmm4
pshufd \$0x39,%xmm0,%xmm0
pxor %xmm4, %xmm3
paddd %xmm3, %xmm2
pshufd \$0x4e,%xmm3,%xmm3
pxor %xmm2, %xmm1
pshufd \$0x93,%xmm2,%xmm2
movdqa %xmm1,%xmm4
pslld \$7, %xmm1
psrld \$25, %xmm4
pxor %xmm4, %xmm1
jnz chacha_blocks_sse2_mainloop2
paddd %xmm8, %xmm0
paddd %xmm9, %xmm1
paddd %xmm10, %xmm2
paddd %xmm11, %xmm3
movdqu 0(%rsi), %xmm12
movdqu 16(%rsi), %xmm13
movdqu 32(%rsi), %xmm14
movdqu 48(%rsi), %xmm15
pxor %xmm12, %xmm0
pxor %xmm13, %xmm1
pxor %xmm14, %xmm2
pxor %xmm15, %xmm3
addq \$64, %rsi
movdqu %xmm0, 0(%rdx)
movdqu %xmm1, 16(%rdx)
movdqu %xmm2, 32(%rdx)
movdqu %xmm3, 48(%rdx)
paddq %xmm5, %xmm11
cmpq \$64, %rbx
jbe chacha_blocks_sse2_done
addq \$64, %rdx
subq \$64, %rbx
jmp chacha_blocks_sse2_below256
chacha_blocks_sse2_done:
movdqu %xmm11, 32(%rdi)
___

if ($flavour !~ /^golang/) {
    $code.=<<___;
popq %r12
popq %r13
___
}

$code.=<<___;
ret
.size chacha_20_core_x64,.-chacha_20_core_x64
___
}

$code =~ s/\`([^\`]*)\`/eval($1)/gem;

print $code;

close STDOUT;
