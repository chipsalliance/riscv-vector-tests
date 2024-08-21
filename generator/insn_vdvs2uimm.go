package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdVs2Uimm(pos int) []string {
	sew32Only_insn := strings.HasPrefix(i.Name, "vaes") || strings.HasPrefix(i.Name, "vsm4")
	sews := iff(sew32Only_insn, []SEW{32}, allSEWs)
	vs2Size := 1
	vdSize := 1

	combinations := i.combinations(
		allLMULs,
		sews,
		[]bool{false, true},
		i.vxrms(),
	)
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		if sew32Only_insn && c.Vl % 4 != 0 {
			c.Vl = (c.Vl + 3) / 4 * 4 
		}
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		vdEMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vdSize), 1))
		vs2EMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vs2Size), 1))
		vdEEW := c.SEW * SEW(vdSize)
		vs2EEW := c.SEW * SEW(vs2Size)
		if vdEEW > SEW(i.Option.XLEN) || vs2EEW > SEW(i.Option.XLEN) {
			res = append(res, "")
			continue
		}

		vd := int(vdEMUL1)
		vs2 := vd * 2

		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		builder.WriteString(i.gWriteIntegerTestData(vs2EMUL1, vs2EEW, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, vs2EMUL1, vs2EEW))

		cases := i.integerTestCases(c.SEW)
		for a := 0; a < len(cases); a++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			switch c.SEW {
			case 8:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d\n",
					i.Name, vd, vs2, convNum[uint8](cases[a][0]) % 32))
			case 16:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d\n",
					i.Name, vd, vs2, convNum[uint16](cases[a][0]) % 32))
			case 32:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d\n",
					i.Name, vd, vs2, convNum[uint32](cases[a][0]) % 32))
			case 64:
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, %d\n",
					i.Name, vd, vs2, convNum[uint64](cases[a][0]) % 32))
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
