package generator

import (
	"fmt"
	"strings"
)

func (i *insn) genCodeVdRs1m() string {
	builder := strings.Builder{}
	builder.WriteString(i.gTestDataAddr())
	builder.WriteString(i.gWriteRandomData(LMUL(1)))
	builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8)))

	for _, c := range i.combinations([]LMUL{1}, []SEW{8}, []bool{false}) {
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gWriteTestData(c.LMUL1, c.SEW, 0))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)\n", i.Name, vd))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))
	}
	return builder.String()
}
