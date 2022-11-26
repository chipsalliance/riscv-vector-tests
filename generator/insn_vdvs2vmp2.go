package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdVs2VmP2() []string {
	vdMask := strings.HasPrefix(i.Name, "vm")

	combinations := i.combinations([]LMUL{1}, []SEW{8}, []bool{false, true})
	if !vdMask {
		combinations = i.combinations(allLMULs, allSEWs, []bool{false, true})
	}

	res := make([]string, 0, len(combinations))
	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		vs2 := 2 * int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(LMUL(1) * 2))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, c.LMUL1, c.SEW))
		builder.WriteString(fmt.Sprintf("addi a0, a0, %d\n", 1*i.vlenb()))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, c.SEW))

		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, v%d%s\n",
			i.Name, vd, vs2, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}
	return res
}
