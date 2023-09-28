package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdVs2ImmV0(pos int) []string {
	combinations := i.combinations(allLMULs, allSEWs, []bool{false})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		vd := int(c.LMUL1)
		vs2 := 2 * int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL1, c.SEW))

		cases := i.integerTestCases(c.SEW)
		for a := 0; a < len(cases); a++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			switch c.SEW {
			case 8:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d, v0\n",
					i.Name, vd, vs2, (int8(convNum[uint8](cases[a][0]))<<3)>>3))
			case 16:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d, v0\n",
					i.Name, vd, vs2, (int8(convNum[uint16](cases[a][0]))<<3)>>3))
			case 32:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d, v0\n",
					i.Name, vd, vs2, (int8(convNum[uint32](cases[a][0]))<<3)>>3))
			case 64:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d, v0\n",
					i.Name, vd, vs2, (int8(convNum[uint64](cases[a][0]))<<3)>>3))
			}
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
			builder.WriteString(i.gMagicInsn(vd))
		}

		res = append(res, builder.String())
	}

	return res
}
