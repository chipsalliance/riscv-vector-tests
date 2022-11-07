package generator

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *insn) genCodeVdRs1mRs2Vm() string {
	getEEW := func(name string) SEW {
		eew, _ := strconv.Atoi(
			strings.TrimSuffix(strings.TrimPrefix(i.Name, "vlse"), ".v"))
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
		builder.WriteString(fmt.Sprintf("%s v%d, (a0), zero%s\n", i.Name, vd, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")
		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		for _, stride := range []int{-1, 0, 1, 3} {
			stride = stride * int(c.SEW) / 8
			builder.WriteString(i.gWriteRandomData(c.LMUL1))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gWriteTestData(c.LMUL1*3, c.SEW, 0))
			builder.WriteString(fmt.Sprintf("addi a0, a0, -%d\n\n", i.vlenb()*int(c.LMUL1)))
			builder.WriteString(i.gWriteTestData(c.LMUL1, c.SEW, 0))
			builder.WriteString(fmt.Sprintf("addi a0, a0,  %d\n\n", i.vlenb()*int(c.LMUL1)))

			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(fmt.Sprintf("li s0, %d # stride\n", stride))
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, (a0), s0%s\n", i.Name, vd, v0t(c.Mask)))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd))
		}
	}
	return builder.String()
}
