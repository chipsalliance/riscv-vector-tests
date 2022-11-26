package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeRdVs2() []string {
	combinations := i.combinations([]LMUL{1}, allSEWs, []bool{false})

	res := make([]string, 0, len(combinations))
	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		vs2 := int(c.LMUL1) * 2
		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, c.SEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s s0, v%d\n", i.Name, vs2))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gMoveScalarToVector("s0", vd, c.SEW))

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}
	return res
}
