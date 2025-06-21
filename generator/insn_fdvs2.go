package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeFdVs2(pos int) []string {
	combinations := i.combinations([]LMUL{1}, i.floatSEWs(), []bool{false}, i.rms())

	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		vd, vs2, _ := getVRegs(c.LMUL1, true, i.Name)

		for r := 0; r < i.Option.Repeat; r += 1 {
			builder.WriteString(i.gWriteRandomData(LMUL(1)))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))
			builder.WriteString(i.gWriteTestData(true, !i.NoTestfloat3, r != 0, c.LMUL1, c.SEW, 0, 1))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, c.SEW))

			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s f0, v%d\n", i.Name, vs2))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gMoveScalarToVector("f0", vd, c.SEW))

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd, c.LMUL1))
		}

		res = append(res, builder.String())
	}
	return res
}
