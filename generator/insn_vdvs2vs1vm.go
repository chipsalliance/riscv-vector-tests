package generator

import (
	"fmt"
	"strings"
)

func (i *insn) genCodeVdVs2Vs1Vm() string {
	builder := strings.Builder{}
	builder.WriteString(i.genCodeTestDataAddr())
	builder.WriteString(i.genCodeWriteRandomData(LMUL(1)))
	builder.WriteString(i.genCodeLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8)))

	for _, c := range i.combinations(allSEWs) {
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		builder.WriteString(i.genCodeWriteRandomData(c.LMUL1))
		builder.WriteString(i.genCodeLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		for idx := 0; idx < 2; idx++ {
			builder.WriteString(i.genCodeWriteTestData(c.LMUL1, c.SEW, idx))
			builder.WriteString(i.genCodeLoadDataIntoRegisterGroup((idx+2)*int(c.LMUL1), c.LMUL1, c.SEW))
		}

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.genCodeVsetvli(c.Vl, c.SEW, c.LMUL))

		builder.WriteString(fmt.Sprintf("%s v%d, v%d, v%d%s\n",
			i.Name, vd, 3*int(c.LMUL1), 2*int(c.LMUL1), v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.genCodeStoreRegisterGroupIntoDataArea(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.genCodeMagicInsn(int(c.LMUL1)))
	}
	return builder.String()
}
