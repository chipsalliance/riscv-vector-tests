package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdVs2Fs1V0(pos int) []string {
	combinations := i.combinations(allLMULs, floatSEWs, []bool{false}, i.vxrms())
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		vd := int(c.LMUL1)
		vs2 := 2 * int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(i.gWriteTestData(true, c.LMUL1, c.SEW, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, c.SEW))
		cases := i.testCases(true, c.SEW)
		for a := 0; a < len(cases); a++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			switch c.SEW {
			case 32:
				builder.WriteString(fmt.Sprintf("li s0, 0x%x\n", convNum[uint32](cases[a][0])))
				builder.WriteString(fmt.Sprintf("fmv.w.x f0, s0\n"))
			case 64:
				builder.WriteString(fmt.Sprintf("li s0, 0x%x\n", convNum[uint64](cases[a][0])))
				builder.WriteString(fmt.Sprintf("fmv.d.x f0, s0\n"))
			}
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, v%d, f0, v0\n",
				i.Name, vd, vs2))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd))
		}
		res = append(res, builder.String())
	}

	return res
}
