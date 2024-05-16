package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdRs1(pos int) []string {
	lmuls := iff(strings.HasSuffix(i.Name, ".s.x"), []LMUL{1}, allLMULs)
	combinations := i.combinations(lmuls, allSEWs, []bool{false}, i.vxrms())

	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		vd, _, _ := getVRegs(c.LMUL1, true, i.Name)

		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		cases := i.integerTestCases(c.SEW)

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		switch c.SEW {
		case 32:
			builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint32](cases[0][0])))
		case 64:
			builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint64](cases[0][0])))
		}

		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, s0\n", i.Name, vd))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd, c.LMUL1))

		res = append(res, builder.String())
	}

	return res
}
