package generator

import (
	"fmt"
	"strings"
)

func (i *insn) genCodeVs3Rs1m() string {
	builder := strings.Builder{}
	builder.WriteString(i.gTestDataAddr())
	for _, c := range i.combinations([]LMUL{1}, []SEW{8}, []bool{false}) {
		builder.WriteString(c.comment())

		vs3 := int(c.LMUL1)
		builder.WriteString(i.gWriteTestData(c.LMUL1, c.SEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))
		builder.WriteString(i.gWriteRandomData(c.LMUL1))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)\n", i.Name, vs3))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vs3))
	}
	return builder.String()
}
