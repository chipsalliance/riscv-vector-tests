package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdImm(pos int) []string {
	combinations := i.combinations(allLMULs, allSEWs, []bool{false}, i.vxrms())

	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		vd, _, _ := getVRegs(c.LMUL1, true, i.Name)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))

		cases := i.integerTestCases(c.SEW)
		for a := 0; a < len(cases); a++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			switch c.SEW {
			case 8:
				builder.WriteString(fmt.Sprintf("%s v%d, %d\n",
					i.Name, vd, (int8(convNum[uint8](cases[a][0]))<<3)>>3))
			case 16:
				builder.WriteString(fmt.Sprintf("%s v%d, %d\n",
					i.Name, vd, (int8(convNum[uint16](cases[a][0]))<<3)>>3))
			case 32:
				builder.WriteString(fmt.Sprintf("%s v%d, %d\n",
					i.Name, vd, (int8(convNum[uint32](cases[a][0]))<<3)>>3))
			case 64:
				builder.WriteString(fmt.Sprintf("%s v%d, %d\n",
					i.Name, vd, (int8(convNum[uint64](cases[a][0]))<<3)>>3))
			}
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd, c.LMUL1))
		}

		res = append(res, builder.String())
	}

	return res
}
