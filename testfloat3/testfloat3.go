package testfloat3

// Wrapper functions on Go side.

// #include <stdlib.h>
// #include "testfloat3.h"
import "C"

func SetLevel(level int) {
	C.genCases_setLevel(C.int(level));
}

func InitF16(numops int) {
	C.srand(2024);
	switch numops {
	case 1:
		C.init_a_f16()
	case 2:
		C.init_ab_f16()
	case 3:
		C.init_abc_f16()
	}
}

func GenF16(numops int) []uint16 {
	var a, b, c uint16 = 0, 0, 0
	switch numops {
	case 1:
		C.gen_a_f16((*C.uint16_t)(&a))
		return []uint16{a}
	case 2:
		C.gen_ab_f16((*C.uint16_t)(&a), (*C.uint16_t)(&b))
		return []uint16{a, b}
	case 3:
		C.gen_abc_f16((*C.uint16_t)(&a), (*C.uint16_t)(&b), (*C.uint16_t)(&c))
		return []uint16{a, b, c}
	}

	return []uint16{}
}

func InitF32(numops int) {
	C.srand(2024);
	switch numops {
	case 1:
		C.init_a_f32()
	case 2:
		C.init_ab_f32()
	case 3:
		C.init_abc_f32()
	}
}

func GenF32(numops int) []float32 {
	var a, b, c float32 = 0, 0, 0
	switch numops {
	case 1:
		C.gen_a_f32((*C.float)(&a))
		return []float32{a}
	case 2:
		C.gen_ab_f32((*C.float)(&a), (*C.float)(&b))
		return []float32{a, b}
	case 3:
		C.gen_abc_f32((*C.float)(&a), (*C.float)(&b), (*C.float)(&c))
		return []float32{a, b, c}
	}

	return []float32{}
}

func InitF64(numops int) {
	C.srand(2024);
	switch numops {
	case 1:
		C.init_a_f64()
	case 2:
		C.init_ab_f64()
	case 3:
		C.init_abc_f64()
	}
}

func GenF64(numops int) []float64 {
	var a, b, c float64 = 0, 0, 0
	switch numops {
	case 1:
		C.gen_a_f64((*C.double)(&a))
		return []float64{a}
	case 2:
		C.gen_ab_f64((*C.double)(&a), (*C.double)(&b))
		return []float64{a, b}
	case 3:
		C.gen_abc_f64((*C.double)(&a), (*C.double)(&b), (*C.double)(&c))
		return []float64{a, b, c}
	}

	return []float64{}
}
