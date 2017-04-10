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
// Created by hchacha20_x64.pl - DO NOT EDIT
// perl hchacha20_x64.pl golang-no-avx hchacha20_x64_amd64.s

// +build amd64,!gccgo,!appengine

// This code was translated into a form compatible with 6a from the public
// domain sources in SUPERCOP: http://bench.cr.yp.to/supercop.html

#include "textflag.h"

___
}

{

if ($flavour =~ /^golang/) {
    $code.=<<___;
TEXT ·hchacha_20_x64(SB),\$0-24
	movq	key+0(FP), DI
	movq	nonce+8(FP), SI
	movq	out+16(FP), DX

___
} else {
    $code.=<<___;
.globl hchacha_20_x64
.type  hchacha_20_x64 ,\@function,2
.align 64
hchacha_20_x64:
___
}

$code.=<<___;
movq \$20, %rbx
movq \$0x3320646e61707865, %rax
movq \$0x6b20657479622d32, %r8
movd %rax, %xmm0
movd %r8, %xmm4
punpcklqdq %xmm4, %xmm0
movdqu 0(%rdi), %xmm1
movdqu 16(%rdi), %xmm2
movdqu 0(%rsi), %xmm3
hchacha_sse2_mainloop:
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
subq \$2, %rbx
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
ja hchacha_sse2_mainloop
movdqu %xmm0, 0(%rdx)
movdqu %xmm3, 16(%rdx)
ret
.size hchacha_20_x64,.-hchacha_20_x64
___
}

$code =~ s/\`([^\`]*)\`/eval($1)/gem;

print $code;

close STDOUT;
