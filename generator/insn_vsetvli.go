package generator

import (
	"fmt"
	"log"
	"strings"
)

func (i *Insn) getvset(curtype int, curvl int, sew SEW, lmul LMUL, vta bool, vma bool, rd int, rs1 int) []int {
	t := 0
	switch lmul {
	case LMUL(1) / 8:
		t = t + 5
	case LMUL(1) / 4:
		t = t + 6
	case LMUL(1) / 2:
		t = t + 7
	case 1:
		t = t + 0
	case 2:
		t = t + 1
	case 4:
		t = t + 2
	case 8:
		t = t + 3
	default:
		log.Fatalln("illegal vlmul")
	}
	switch sew {
	case 8:
		t = t + (0 << 3)
	case 16:
		t = t + (1 << 3)
	case 32:
		t = t + (2 << 3)
	case 64:
		t = t + (3 << 3)
	default:
		log.Fatalln("illegal vsew")
	}
	if vta {
		t = t + (1 << 6)
	}
	if vma {
		t = t + (1 << 7)
	}
	res := i.vtype(t, curtype, curvl, rd, rs1, t)
	return res
}

func (i *Insn) genCodevsetvli() []string {
	combinations := i.vsetvlicombinations(
		allLMULs,
		allSEWs,
		[]bool{false, true},
		[]bool{false, true},
	)
	ncase := 3
	res := make([]string, 0, len(combinations))
	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("# ------combination test begin---------\n"))
		builder.WriteString(c.comment())

		cases := i.testCases(false, 8)
		curvtype := 0
		curvl := 0
		for _, cs := range cases {
			builder.WriteString(fmt.Sprintf("# ------case test begin---------\n"))
			builder.WriteString(fmt.Sprintf("li t0, %d\n", cs[0]))
			builder.WriteString(fmt.Sprintf("li t1, %d\n", cs[1]))
			builder.WriteString(fmt.Sprintf("vsetvli t0, t1, %s,%s,%s,%s\n",
				c.SEW, c.LMUL, ta(c.vta), ma(c.vma)))
			t := i.getvset(curvtype, curvl, c.SEW, c.LMUL, c.vta, c.vma, int(cs[0].(uint8)), int(cs[1].(uint8)))

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
	}
	return res
}
