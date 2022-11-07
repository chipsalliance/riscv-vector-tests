package generator

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *insn) genCodeVdRs1mVm() string {
	getEEW := func(name string) SEW {
		eew, _ := strconv.Atoi(
			strings.TrimSuffix(strings.TrimPrefix(i.Name, "vle"), ".v"))
		return SEW(eew)
	}

	builder := strings.Builder{}
	builder.WriteString(i.genCodeTestDataAddr())
	builder.WriteString(i.genCodeWriteRandomData(LMUL(1)))
	builder.WriteString(i.genCodeLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8)))

	for _, c := range i.combinations([]SEW{getEEW(i.Name)}) {
		builder.WriteString(c.comment())

		vd := 1 * int(c.LMUL1)

		builder.WriteString(i.genCodeWriteRandomData(c.LMUL1))
		builder.WriteString(i.genCodeLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.genCodeWriteTestData(c.LMUL1, c.SEW, 0))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.genCodeVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)%s\n", i.Name, vd, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.genCodeStoreRegisterGroupIntoDataArea(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.genCodeMagicInsn(int(c.LMUL1)))
	}
	return builder.String()
}
