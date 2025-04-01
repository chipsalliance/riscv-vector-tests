package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdRs1mRs2Vm(pos int) []string {
	nfields := getNfields(i.Name)
	combinations := i.combinations(
		nfieldsLMULs(nfields),
		[]SEW{getEEW(i.Name)},
		[]bool{false, true},
		i.rms(),
	)
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		vd := int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		for nf := 0; nf < nfields; nf++ {
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vd+(nf*int(c.LMUL1)), c.LMUL1, c.SEW))
		}
		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1*LMUL(nfields), c.SEW, 0))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0), zero%s\n", i.Name, vd, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd, c.LMUL1))

		for _, stride := range []int{minStride, 0, 1, maxStride} {
			stride = stride * int(c.SEW) / 8
			builder.WriteString(i.gWriteRandomData(c.LMUL1))
			for nf := 0; nf < nfields; nf++ {
				builder.WriteString(i.gLoadDataIntoRegisterGroup(vd+(nf*int(c.LMUL1)), c.LMUL1, c.SEW))
			}
			builder.WriteString(i.gWriteIntegerTestData(c.LMUL1*LMUL(nfields*stride), c.SEW, 0))
			builder.WriteString(fmt.Sprintf("li a5, %d\n", -minStride*i.vlenb()*int(c.LMUL1)))
			builder.WriteString("add a0, a0, a5\n")

			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(fmt.Sprintf("li s0, %d # stride\n", stride))
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, (a0), s0%s\n", i.Name, vd, v0t(c.Mask)))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gResultDataAddr())
			for nf := 0; nf < nfields; nf++ {
				builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd+(nf*int(c.LMUL1)), c.LMUL1, c.SEW))
				builder.WriteString(i.gMagicInsn(vd+(nf*int(c.LMUL1)), c.LMUL1))
			}
		}
		res = append(res, builder.String())
	}
	return res
}
