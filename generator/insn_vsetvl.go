package generator

import (
	"fmt"
	"log"
	"math"
	"strings"
)

func getSew(vtype int64) int {
	sewMask := 0x38
	sew := vtype & int64(sewMask)
	if (sew >> 3) > 3 {
		log.Fatalln("illegal vsew")
	}
	return 8 << (sew >> 3)
}

func getLmul(vtype int64) LMUL {
	lmulmask := 0x7
	t := vtype & int64(lmulmask)
	lmul := LMUL(int(1) << t)
	if t >= 5 {
		lmul = 1 / LMUL(int(1)<<(8-t))
	}
	return LMUL(lmul)
}

func getVlmax(vtype int64, vlen int) int {
	sew := getSew(vtype)
	vlmax := LMUL((vlen / sew)) * getLmul(vtype)
	return int(vlmax)
}

func bits64(val int64, first int, last int) int64 {
	nbits := first - last + 1
	return (val >> last) & (1<<nbits - 1)
}

func (v *vtype) vtypeRaw(XLEN int, VLEN int, newVtypeRaw int64, curVtypeRaw int64, vl int64) []int64 {
	res := make([]int64, 0, 3)
	var vstart int64 = 0
	res = append(res, vstart)

	vlmax := int64(getVlmax(curVtypeRaw, VLEN))
	vtypeRaw := newVtypeRaw
	if vtypeRaw != curVtypeRaw {
		vlmax = int64(getVlmax(vtypeRaw, VLEN))
		newLmul := getLmul(vtypeRaw)
		newVsew := getSew(vtypeRaw)
		// We assume ELEN is consistent with XLEN.
		new_vill :=
			!(newLmul >= 0.125 && newLmul <= 8) ||
				float64(newVsew) > math.Min(float64(newLmul), float64(1))*float64(XLEN) ||
				bits64(vtypeRaw, XLEN-2, 8) != 0
		if new_vill {
			vlmax = 0
			vtypeRaw = 1 << (XLEN - 1)
		}
	}
	res = append(res, vtypeRaw)

	if vl < 0 || vl > vlmax || vlmax == 0 {
		vl = vlmax
	}
	res = append(res, vl)

	return res
}

func (i *Insn) genCodevsetvl(pos int) []string {
	res := make([]string, 0, 1)
	ncase := 3
	builder := strings.Builder{}
	cases := i.testCases(false, 8)

	curvtype := int64(0)
	curvl := int64(0)
	for _, cs := range cases {
		for idx := 0; idx < 3; idx++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")

			builder.WriteString(fmt.Sprintf("li t0, %d\n", cs[0]))
			builder.WriteString(fmt.Sprintf("li t1, %d\n", cs[1]))
			builder.WriteString(fmt.Sprintf("li t2, %d\n", cs[2]))

			rd := "t0"
			rs := "t1"
			if idx > 0 {
				rs = "zero"
			}
			if idx == 2 {
				rd = "zero"
			}

			builder.WriteString(fmt.Sprintf("vsetvl %s, %s, t2\n", rd, rs))

			var v vtype
			t := v.vtypeRaw(int(i.Option.XLEN), int(i.Option.VLEN), int64(cs[2].(uint8)), curvtype, int64(cs[1].(uint8)))

			builder.WriteString(fmt.Sprintf("csrr a4, vstart\n"))
			builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[0]))
			ncase = ncase + 1

			builder.WriteString(fmt.Sprintf("csrr a4, vtype\n"))
			if t[1] == 1<<(i.Option.XLEN-1) {
				builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, uint64(1<<(i.Option.XLEN-1))))
			} else {
				builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[1]))
			}
			ncase = ncase + 1

			switch idx {
			case 0:
				curvl = t[2]
			case 1:
				if t[1] == 1<<(i.Option.XLEN-1) { // vill
					curvl = t[2]
				} else {
					curvl = int64(getVlmax(curvtype, int(i.Option.VLEN)))
				}
			}
			curvtype = t[1]

			builder.WriteString(fmt.Sprintf("csrr a4, vl\n"))
			builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, curvl))
			ncase = ncase + 1
			builder.WriteString("# -------------- TEST END   --------------\n")
		}
	}
	res = append(res, builder.String())

	return res
}
