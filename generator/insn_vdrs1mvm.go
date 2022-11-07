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
	builder.WriteString(i.gTestDataAddr())
	builder.WriteString(i.gWriteRandomData(LMUL(1)))
	builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8)))

	for _, c := range i.combinations(allLMULs, []SEW{getEEW(i.Name)}, []bool{false, true}) {
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gWriteTestData(c.LMUL1, c.SEW, 0))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)%s\n", i.Name, vd, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))
	}
	return builder.String()
}
