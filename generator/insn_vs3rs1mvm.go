package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVs3Rs1mVm() []string {
	combinations := i.combinations(
		allLMULs,
		[]SEW{getEEW(i.Name)},
		[]bool{false, true})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(i.gTestDataAddr())
		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		builder.WriteString(c.comment())

		vs3 := int(c.LMUL1)
		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))
		builder.WriteString(i.gWriteRandomData(c.LMUL1))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)%s\n", i.Name, vs3, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vs3))

		res = append(res, builder.String())
	}
	return res
}
