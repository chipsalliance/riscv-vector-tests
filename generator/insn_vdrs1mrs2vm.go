package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdRs1mRs2Vm() []string {
	combinations := i.combinations(
		allLMULs,
		[]SEW{getEEW(i.Name)},
		[]bool{false, true})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(i.gTestDataAddr())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(64)))

		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 0))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0), zero%s\n", i.Name, vd, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")
		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		for _, stride := range []int{minStride, 0, 1, maxStride} {
			stride = stride * int(c.SEW) / 8
			builder.WriteString(i.gWriteRandomData(c.LMUL1))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gWriteIntegerTestData(c.LMUL1*strides, c.SEW, 0))

			builder.WriteString(fmt.Sprintf("li a5, %d\n", -minStride*i.vlenb()*int(c.LMUL1)))
			builder.WriteString("add a0, a0, a5\n")

			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(fmt.Sprintf("li s0, %d # stride\n", stride))
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, (a0), s0%s\n", i.Name, vd, v0t(c.Mask)))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd))
		}
		res = append(res, builder.String())
	}
	return res
}
