package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *insn) genCodeVdVs2Vs1Vm() []string {
	float := strings.HasPrefix(i.Name, "vf")
	vdWidening := strings.HasPrefix(i.Name, "vw")
	vs2Widening := strings.HasSuffix(i.Name, "wv")
	vdSize := iff(vdWidening, 2, 1)
	vs2Size := iff(vs2Widening, 2, 1)

	sews := iff(float, floatSEWs, allSEWs)
	sews = iff(vdWidening || vs2Widening, sews[:len(sews)-1], sews)
	combinations := i.combinations(
		iff(vdWidening || vs2Widening, wideningMULs, allLMULs),
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
		vss := []int{
			vd * 2,
			vd*2 + int(vs2EMUL1),
		}
		builder.WriteString(i.gWriteRandomData(vdEMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, vdEMUL1, SEW(8)))

		builder.WriteString(i.gWriteTestData(float, c.LMUL1, c.SEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vss[0], c.LMUL1, c.SEW))

		builder.WriteString(i.gWriteTestData(float, vs2EMUL1, vs2EEW, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vss[1], vs2EMUL1, vs2EEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))

		builder.WriteString(fmt.Sprintf("%s v%d, v%d, v%d%s\n",
			i.Name, vd, vss[1], vss[0], v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, vdEMUL1, vdEEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}

	return res
}
