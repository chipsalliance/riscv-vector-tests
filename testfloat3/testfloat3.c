#include "testfloat3.h"

void init_a_f32(void) { genCases_f32_a_init(); }
void gen_a_f32(float32_t *a) {
  genCases_f32_a_next();
  *a = genCases_f32_a;
}

void init_ab_f32(void) { genCases_f32_ab_init(); }
void gen_ab_f32(float32_t *a, float32_t *b) {
  genCases_f32_ab_next();
  *a = genCases_f32_a;
  *b = genCases_f32_b;
}

void init_abc_f32(void) { genCases_f32_abc_init(); }
void gen_abc_f32(float32_t *a, float32_t *b, float32_t *c) {
  genCases_f32_abc_next();
  *a = genCases_f32_a;
  *b = genCases_f32_b;
  *c = genCases_f32_c;
}

void init_a_f64(void) { genCases_f64_a_init(); }
void gen_a_f64(float64_t *a) {
  genCases_f64_a_next();
  *a = genCases_f64_a;
}

void init_ab_f64(void) { genCases_f64_ab_init(); }
void gen_ab_f64(float64_t *a, float64_t *b) {
  genCases_f64_ab_next();
  *a = genCases_f64_a;
  *b = genCases_f64_b;
}

void init_abc_f64(void) { genCases_f64_abc_init(); }
void gen_abc_f64(float64_t *a, float64_t *b, float64_t *c) {
  genCases_f64_abc_next();
  *a = genCases_f64_a;
  *b = genCases_f64_b;
  *c = genCases_f64_c;
}
