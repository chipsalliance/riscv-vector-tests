package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdVs2Vs1V0() []string {
	combinations := i.combinations(allLMULs, allSEWs, []bool{false})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		vd := int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		for idx := 0; idx < 2; idx++ {
			builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, idx))
			builder.WriteString(i.gLoadDataIntoRegisterGroup((idx+2)*int(c.LMUL1), c.LMUL1, c.SEW))
		}

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))

		builder.WriteString(fmt.Sprintf("%s v%d, v%d, v%d, v0\n",
			i.Name, vd, 3*int(c.LMUL1), 2*int(c.LMUL1)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}

	return res
}
