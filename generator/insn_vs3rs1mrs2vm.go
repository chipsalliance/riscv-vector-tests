package generator

import (
	"fmt"
	"strings"
)

func (i *insn) genCodeVs3Rs1mRs2Vm() []string {
	combinations := i.combinations(
		allLMULs,
		[]SEW{getEEW(i.Name)},
		[]bool{false, true})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(i.gTestDataAddr())
		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8)))

		builder.WriteString(c.comment())

		vs3 := int(c.LMUL1)
		for _, s := range []int{minStride, 0, 1, maxStride} {
			stride := s * int(c.SEW) / 8
			builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 0))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))
			builder.WriteString(i.gWriteRandomData(c.LMUL1 * strides))

			builder.WriteString(fmt.Sprintf("li a5, %d\n", -minStride*i.vlenb()*int(c.LMUL1)))
			builder.WriteString("add a0, a0, a5\n")

			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(fmt.Sprintf("li s0, %d # stride\n", stride))
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, (a0), s0%s\n", i.Name, vs3, v0t(c.Mask)))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(fmt.Sprintf("li a5, %d\n", minStride*i.vlenb()*int(c.LMUL1)))
			builder.WriteString("add a0, a0, a5\n")

			builder.WriteString("mv a4, a0\n")
			for a := 0; a < strides; a++ {
				builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))
				builder.WriteString(i.gMagicInsn(vs3))
				builder.WriteString(fmt.Sprintf("li a5, %d\n", i.vlenb()*int(c.LMUL1)))
				builder.WriteString(fmt.Sprintf("add a4, a4, a5\n"))
				builder.WriteString("mv a0, a4\n")
			}

			builder.WriteString(fmt.Sprintf("li a5, %d\n", -strides*i.vlenb()*int(c.LMUL1)))
			builder.WriteString("add a0, a0, a5\n")
		}

		res = append(res, builder.String())
	}
	return res
}
