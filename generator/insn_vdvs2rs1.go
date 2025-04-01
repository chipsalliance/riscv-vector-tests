package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdVs2Rs1(pos int) []string {
	combinations := i.combinations(allLMULs, allSEWs, []bool{false}, i.rms())
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		vd, vs2, _ := getVRegs(c.LMUL1, true, i.Name)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, c.SEW))

		cases := i.integerTestCases(c.SEW)
		for a := 0; a < len(cases); a++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			switch c.SEW {
			case 8:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint8](cases[a][0])))
			case 16:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint16](cases[a][0])))
			case 32:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint32](cases[a][0])))
			case 64:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint64](cases[a][0])))
			}
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, v%d, s0\n",
				i.Name, vd, vs2))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd, c.LMUL1))
		}

		res = append(res, builder.String())
	}

	return res
}
