package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *insn) genCodeVdVs2Rs1Vm() []string {
	vdWidening := strings.HasPrefix(i.Name, "vw")
	vs2Widening := strings.HasSuffix(i.Name, ".wx")
	vdSize := iff(vdWidening, 2, 1)
	vs2Size := iff(vs2Widening, 2, 1)

	sews := iff(vdWidening, allSEWs[:len(allSEWs)-1], allSEWs)
	combinations := i.combinations(
		iff(vdWidening, wideningMULs, allLMULs),
		sews,
		[]bool{false, true},
	)
	res := make([]string, 0, len(combinations))

	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(i.gTestDataAddr())
		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8)))

		builder.WriteString(c.comment())

		vdEMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vdSize), 1))
		vs2EMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vs2Size), 1))
		vdEEW := c.SEW * SEW(vdSize)
		vs2EEW := c.SEW * SEW(vs2Size)
		vd := int(vdEMUL1)
		vs2 := vd * 2
		builder.WriteString(i.gWriteRandomData(vdEMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, vdEMUL1, SEW(8)))

		builder.WriteString(i.gWriteIntegerTestData(vs2EMUL1, vs2EEW, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, vs2EMUL1, vs2EEW))

		cases := i.integerTestCases(c.SEW)
		for a := 0; a < len(cases); a++ {
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			switch c.SEW {
			case 8:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint8](cases[a][0])))
			case 16:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint16](cases[a][0])))
			case 32:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint32](cases[a][0])))
			case 64:
				builder.WriteString(fmt.Sprintf("li s0, %d\n", convNum[uint64](cases[a][0])))
			}
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, v%d, s0%s\n",
				i.Name, vd, vs2, v0t(c.Mask)))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gStoreRegisterGroupIntoData(vd, vdEMUL1, vdEEW))
			builder.WriteString(i.gMagicInsn(vd))
		}

		res = append(res, builder.String())
	}

	return res
}
