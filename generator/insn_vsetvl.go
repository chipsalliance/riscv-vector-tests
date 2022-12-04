package generator

import (
	"fmt"
	"log"
	"math"
	"strings"
)

func (i *Insn) getsew(vtype int) int {
	vsewmask := 0x38
	vsew := vtype & vsewmask
	vsew = vsew >> 3
	switch vsew {
	case 0x0:
		return 8
	case 0x1:
		return 16
	case 0x2:
		return 32
	case 0x3:
		return 64
	default:
		log.Fatalln("illegal vsew")
	}
	return -1
}

func (i *Insn) getlmul(vtype int) LMUL {
	lmulmask := 0x7
	lmul := vtype & lmulmask
	switch lmul {
	case 0:
		return LMUL(1)
	case 1:
		return LMUL(2)
	case 2:
		return LMUL(4)
	case 3:
		return LMUL(8)
	case 5:
		return LMUL(LMUL(1) / 8)
	case 6:
		return LMUL(LMUL(1) / 4)
	case 7:
		return LMUL(LMUL(1) / 2)
	default:
		log.Fatalln("illegal vlmul")
	}
	return -1
}

func (i *Insn) getVlmax(vtype int, vlen int) int {
	sew := i.getsew(vtype)
	vlmax := LMUL((vlen / sew)) * i.getlmul(vtype)
	return int(vlmax)
}

func (i *Insn) bits32(vtype int, first int, last int) int {
	nbits := first - last + 1
	nbitsmask := -1
	if nbits < 32 {
		nbitsmask = 1<<nbits - 1
	}
	return (vtype >> last) & nbitsmask
}

func (i *Insn) vtype(newtype int, curtype int, curvl int, rd int, rs1 int, rs2 int) []int {
	vlen := int(i.Option.VLEN)
	res := make([]int, 0, 3)
	// vstart
	res = append(res, 0)

	vlmax := i.getVlmax(curtype, vlen)
	new_type := newtype
	if newtype != curtype {
		vlmax = i.getVlmax(newtype, vlen)
		newlmul := i.getlmul(newtype)
		newsew := i.getsew(newtype)
		// We assume ELEN is consistent with XLEN.
		new_vill :=
			!(newlmul >= 0.125 && newlmul <= 8) ||
				float64(newsew) > math.Min(float64(newlmul), float64(1))*float64(i.Option.XLEN) ||
				i.bits32(newtype, int(i.Option.XLEN)-2, 8) != 0
		if new_vill {
			vlmax = 0
			// vtype flag when vill is set.
			new_type = 1 << 31
		}
	}
	// vtype
	res = append(res, new_type)

	new_vl := curvl
	if vlmax == 0 {
		new_vl = 0
	} else if rd == 0 && rs1 == 0 {
		if curvl > int(vlmax) {
			new_vl = int(vlmax)
		} else {
			new_vl = curvl
		}
	} else if rd != 0 && rs1 == 0 {
		new_vl = int(vlmax)
	} else if rs1 != 0 {
		if rs1 > int(vlmax) {
			new_vl = int(vlmax)
		} else {
			new_vl = rs1
		}
	}
	// vl
	res = append(res, new_vl)

	return res
}

func (i *Insn) genCodevsetvl() []string {
	res := make([]string, 0, 1)
	ncase := 3
	builder := strings.Builder{}
	cases := i.testCases(false, 8)

	curvtype := 0
	curvl := 0
	for _, cs := range cases {
		builder.WriteString(fmt.Sprintf("# ------case test begin---------\n"))

		builder.WriteString(fmt.Sprintf("li t0, %d\n", cs[0]))
		builder.WriteString(fmt.Sprintf("li t1, %d\n", cs[1]))
		builder.WriteString(fmt.Sprintf("li t2, %d\n", cs[2]))
		builder.WriteString(fmt.Sprintf("vsetvl t0, t1, t2\n"))

		t := i.vtype(int(cs[2].(uint8)), curvtype, curvl, int(cs[0].(uint8)), int(cs[1].(uint8)), int(cs[2].(uint8)))

		builder.WriteString(fmt.Sprintf("csrr a4, vstart\n"))
		builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[0]))
		ncase = ncase + 1

		builder.WriteString(fmt.Sprintf("csrr a4, vtype\n"))
		if t[1] == 1<<31 {
			builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, uint64(1<<(i.Option.XLEN-1))))
		} else {
			builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[1]))
		}
		ncase = ncase + 1

		builder.WriteString(fmt.Sprintf("csrr a4, vl\n"))
		builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[2]))
		ncase = ncase + 1

		curvtype = t[1]
		curvl = t[2]

		builder.WriteString(fmt.Sprintf("# ------case test end---------\n"))
	}
	res = append(res, builder.String())

	return res
}
