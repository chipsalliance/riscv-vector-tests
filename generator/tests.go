package generator

import (
	"math"
	"strconv"
	"strings"
	"unsafe"

	f16 "github.com/x448/float16"
)

type num interface {
	uint8 | uint16 | uint32 | uint64 | float32 | float64 | string
}

func convNum[T num](n any) T {
	var res T
	var ok bool
	if res, ok = n.(T); !ok {
		res = T(n.(uint8))
	}
	return res
}

func parseCustomFloat16(n string) uint16 {
	switch n {
	case "nan":
		return 0x7e01
	case "-nan":
		return 0xfe01
	case "inf":
		return 0x7c00
	case "-inf":
		return 0xfc00
	case "quiet_nan":
		return 0x7e00
	case "signaling_nan":
		return 0x7d00
	case "smallest_nonzero_float":
		return 0x0001
	case "largest_subnormal_float":
		return 0x03ff
	case "smallest_normal_float":
		return 0x0400
	case "max_float":
		return 0x7bff
	case "-smallest_nonzero_float":
		return 0x8001
	case "-largest_subnormal_float":
		return 0x83ff
	case "-smallest_normal_float":
		return 0x8400
	case "-max_float":
		return 0xfbff
	default:
		v, _ := strconv.ParseFloat(n, 32)
		return f16.Fromfloat32(float32(v)).Bits()
	}
}

func parseCustomFloat32(n string) float32 {
	var val uint32 = 0x800000
	var smallestNormalFloat32 = *(*float32)(unsafe.Pointer(&val))
	switch n {
	case "nan":
		val = 0x7fc00001
		return *(*float32)(unsafe.Pointer(&val))
	case "-nan":
		val = 0xffc00001
		return *(*float32)(unsafe.Pointer(&val))
	case "inf":
		val = 0x7f800000
		return *(*float32)(unsafe.Pointer(&val))
	case "-inf":
		val = 0xff800000
		return *(*float32)(unsafe.Pointer(&val))
	case "quiet_nan":
		val = 0x7fc00000
		return *(*float32)(unsafe.Pointer(&val))
	case "signaling_nan":
		val = 0x7fa00000
		return *(*float32)(unsafe.Pointer(&val))
	case "smallest_nonzero_float":
		return math.SmallestNonzeroFloat32
	case "largest_subnormal_float":
		return smallestNormalFloat32 - math.SmallestNonzeroFloat32
	case "smallest_normal_float":
		return smallestNormalFloat32
	case "max_float":
		return math.MaxFloat32
	case "-smallest_nonzero_float":
		return -math.SmallestNonzeroFloat32
	case "-largest_subnormal_float":
		return math.SmallestNonzeroFloat32 - smallestNormalFloat32
	case "-smallest_normal_float":
		return -smallestNormalFloat32
	case "-max_float":
		return -math.MaxFloat32
	default:
		v, _ := strconv.ParseFloat(n, 32)
		return float32(v)
	}
}

func parseCustomFloat64(n string) float64 {
	var val uint64 = 0x10000000000000
	var smallestNormalFloat64 = *(*float64)(unsafe.Pointer(&val))
	switch n {
	case "nan":
		val = 0x7ff8000000000001
		return *(*float64)(unsafe.Pointer(&val))
	case "-nan":
		val = 0xfff8000000000001
		return *(*float64)(unsafe.Pointer(&val))
	case "inf":
		val = 0x7ff0000000000000
		return *(*float64)(unsafe.Pointer(&val))
	case "-inf":
		val = 0xfff0000000000000
		return *(*float64)(unsafe.Pointer(&val))
	case "quiet_nan":
		val = 0x7ff8000000000000
		return *(*float64)(unsafe.Pointer(&val))
	case "signaling_nan":
		val = 0x7ff4000000000000
		return *(*float64)(unsafe.Pointer(&val))
	case "smallest_nonzero_float":
		return math.SmallestNonzeroFloat64
	case "largest_subnormal_float":
		return smallestNormalFloat64 - math.SmallestNonzeroFloat64
	case "smallest_normal_float":
		return smallestNormalFloat64
	case "max_float":
		return math.MaxFloat64
	case "-smallest_nonzero_float":
		return -math.SmallestNonzeroFloat64
	case "-largest_subnormal_float":
		return math.SmallestNonzeroFloat64 - smallestNormalFloat64
	case "-smallest_normal_float":
		return -smallestNormalFloat64
	case "-max_float":
		return -math.MaxFloat64
	default:
		v, _ := strconv.ParseFloat(n, 64)
		return v
	}
}

type testCase[T num] []T

type tests struct {
	Base  []testCase[uint8]  `toml:"base"`
	SEW8  []testCase[uint8]  `toml:"sew8"`
	SEW16 []testCase[uint16] `toml:"sew16"`
	SEW32 []testCase[uint32] `toml:"sew32"`

	// Go toml cannot parse uint64/float32/float64 well, parse it ourself.
	SEW64_ []testCase[string] `toml:"sew64"`
	SEW64  []testCase[uint64] `toml:"-"`

	FSEW16_ []testCase[string] `toml:"fsew16"`
	FSEW16  []testCase[uint16] `toml:"-"`

	FSEW32_ []testCase[string]  `toml:"fsew32"`
	FSEW32  []testCase[float32] `toml:"-"`

	FSEW64_ []testCase[string]  `toml:"fsew64"`
	FSEW64  []testCase[float64] `toml:"-"`
}

func (t *tests) initialize() error {
	var err error
	for i, ss := range t.SEW64_ {
		t.SEW64 = append(t.SEW64, make([]uint64, len(ss)))
		for j, s := range ss {
			t.SEW64[i][j], err = strconv.ParseUint(
				strings.TrimPrefix(s, "0x"), 16, 64)
			if err != nil {
				return err
			}
		}
	}

	for i, ss := range t.FSEW16_ {
		t.FSEW16 = append(t.FSEW16, make([]uint16, len(ss)))
		for j, s := range ss {
			t.FSEW16[i][j] = parseCustomFloat16(s)
			if err != nil {
				return err
			}
		}
	}

	for i, ss := range t.FSEW32_ {
		t.FSEW32 = append(t.FSEW32, make([]float32, len(ss)))
		for j, s := range ss {
			t.FSEW32[i][j] = parseCustomFloat32(s)
			if err != nil {
				return err
			}
		}
	}

	for i, ss := range t.FSEW64_ {
		t.FSEW64 = append(t.FSEW64, make([]float64, len(ss)))
		for j, s := range ss {
			t.FSEW64[i][j] = parseCustomFloat64(s)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
