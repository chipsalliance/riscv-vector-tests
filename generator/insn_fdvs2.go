package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeFdVs2(pos int) []string {
	combinations := i.combinations([]LMUL{1}, floatSEWs, []bool{false}, i.vxrms())

	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		vd := int(c.LMUL1)
		vs2 := int(c.LMUL1) * 2
		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(i.gWriteTestData(true, c.LMUL1, c.SEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, c.SEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s f0, v%d\n", i.Name, vs2))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gMoveScalarToVector("f0", vd, c.SEW))

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}
	return res
}
