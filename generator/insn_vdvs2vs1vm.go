package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *insn) genCodeVdVs2Vs1Vm() []string {
	float := strings.HasPrefix(i.Name, "vf")
	vdWidening := strings.HasPrefix(i.Name, "vw")
	vdSize := iff(vdWidening, 2, 1)

	sews := iff(float, floatSEWs, allSEWs)
	sews = iff(vdWidening, sews[:len(sews)-1], sews)
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

		emul1 := LMUL(math.Max(float64(int(c.LMUL)*vdSize), 1))
		vd := int(emul1)
		vss := []int{
			vd * 2,
			vd*2 + int(c.LMUL1),
		}
		builder.WriteString(i.gWriteRandomData(emul1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, emul1, SEW(8)))

		for idx, vs := range vss {
			builder.WriteString(i.gWriteTestData(float, c.LMUL1, c.SEW, idx))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs, c.LMUL1, c.SEW))
		}

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))

		builder.WriteString(fmt.Sprintf("%s v%d, v%d, v%d%s\n",
			i.Name, vd, vss[1], vss[0], v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, emul1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}

	return res
}
