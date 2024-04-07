package generator

import (
	"fmt"
	"math"
	"strings"
)

func (v *vtype) vtypeImm(XLEN int, VLEN int, curVtypeRaw int64, vl int64) []int64 {
	t := int64(math.Ilogb(float64(v.lmul)) & 7)
	t = t + int64(math.Ilogb(float64(v.sew/8))<<3)
	if v.vta {
		t = t + (1 << 6)
	}
	if v.vma {
		t = t + (1 << 7)
	}
	res := v.vtypeRaw(XLEN, VLEN, t, curVtypeRaw, vl)
	return res
}

func (i *Insn) genCodevsetvli(pos int) []string {
	combinations := i.vsetvlicombinations(
		allLMULs,
		allSEWs,
		[]bool{false, true},
		[]bool{false, true},
	)
	ncase := 3
	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		cases := i.testCases(false, 8)
		curvtype := int64(0)
		curvl := int64(0)
		for _, cs := range cases {
			for idx := 0; idx < 3; idx++ {
				builder.WriteString("# -------------- TEST BEGIN --------------\n")
				builder.WriteString(fmt.Sprintf("li t0, %d\n", cs[0]))
				builder.WriteString(fmt.Sprintf("li t1, %d\n", cs[1]))

				rd := "t0"
				rs := "t1"
				if idx > 0 {
					rs = "zero"
				}
				if idx == 2 {
					rd = "zero"
				}

				builder.WriteString(fmt.Sprintf("vsetvli %s, %s, %s,%s,%s,%s\n", rd, rs,
					c.SEW, c.LMUL, ta(c.vta), ma(c.vma)))

				v := vtype{float32(c.LMUL), int(c.SEW), c.vta, c.vma}
				t := v.vtypeImm(int(i.Option.XLEN), int(i.Option.VLEN), curvtype, int64(cs[1].(uint8)))

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
	}
	return res
}
