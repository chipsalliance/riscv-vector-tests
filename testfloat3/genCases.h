/*============================================================================

This C header file is part of TestFloat, Release 3e, a package of programs for
testing the correctness of floating-point arithmetic complying with the IEEE
Standard for Floating-Point, by John R. Hauser.

Copyright 2011, 2012, 2013, 2014, 2015, 2017 The Regents of the University of
California.  All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions, and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions, and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the University nor the names of its contributors may
    be used to endorse or promote products derived from this software without
    specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS "AS IS", AND ANY
EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE, ARE
DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

=============================================================================*/

#include <stdbool.h>
#include <stdint.h>

typedef uint16_t float16_t;
typedef float float32_t;
typedef double float64_t;

extern int genCases_level;

void genCases_setLevel(int);

extern uint_fast64_t genCases_total;
extern bool genCases_done;

void genCases_f16_a_init(void);
void genCases_f16_a_next(void);
void genCases_f16_ab_init(void);
void genCases_f16_ab_next(void);
void genCases_f16_abc_init(void);
void genCases_f16_abc_next(void);
extern float16_t genCases_f16_a, genCases_f16_b, genCases_f16_c;

void genCases_f32_a_init(void);
void genCases_f32_a_next(void);
void genCases_f32_ab_init(void);
void genCases_f32_ab_next(void);
void genCases_f32_abc_init(void);
void genCases_f32_abc_next(void);
extern float32_t genCases_f32_a, genCases_f32_b, genCases_f32_c;

void genCases_f64_a_init(void);
void genCases_f64_a_next(void);
void genCases_f64_ab_init(void);
void genCases_f64_ab_next(void);
void genCases_f64_abc_init(void);
void genCases_f64_abc_next(void);
extern float64_t genCases_f64_a, genCases_f64_b, genCases_f64_c;
