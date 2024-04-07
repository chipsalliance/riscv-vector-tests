package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodevsetivli(pos int) []string {
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
		for _, cs := range cases {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(fmt.Sprintf("li t0, %d\n", cs[0]))
			builder.WriteString(fmt.Sprintf("vsetivli t0, %d, %s,%s,%s,%s\n",
				cs[1], c.SEW, c.LMUL, ta(c.vta), ma(c.vma)))

			v := vtype{float32(c.LMUL), int(c.SEW), c.vta, c.vma}
			t := v.vtypeImm(int(i.Option.XLEN), int(i.Option.VLEN), curvtype, int64(cs[1].(uint8)))

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

			builder.WriteString("# -------------- TEST BEGIN --------------\n")
		}
		res = append(res, builder.String())
	}
	return res
}
