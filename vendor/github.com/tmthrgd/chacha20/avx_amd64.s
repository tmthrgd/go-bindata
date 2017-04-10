// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

// +build amd64,!gccgo,!appengine

#include "textflag.h"

// func hasAVX() bool
// returns whether AVX or AVX2 is supported
TEXT ·hasAVX(SB),NOSPLIT,$0
	MOVB runtime·support_avx(SB), CX
	MOVB CX, avx+0(FP)
	MOVB runtime·support_avx2(SB), CX
	MOVB CX, avx2+1(FP)
	RET
