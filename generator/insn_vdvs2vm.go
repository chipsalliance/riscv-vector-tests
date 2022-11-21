package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdVs2Vm() []string {
	vdWidening := strings.HasPrefix(i.Name, "vfw")
	vdNarrowing := strings.HasPrefix(i.Name, "vfn")
	vdSize := iff(vdWidening, 2, 1)
	vs2Size := iff(vdNarrowing, 2, 1)

	lmuls := iff(vdWidening || vdNarrowing, wideningMULs, allLMULs)
	sews := iff(vdWidening || vdNarrowing, floatSEWs[:len(floatSEWs)-1], floatSEWs)
	combinations := i.combinations(lmuls, sews, []bool{false, true})

	res := make([]string, 0, len(combinations))
	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(i.gTestDataAddr())
		builder.WriteString(c.comment())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, c.LMUL1, SEW(8)))

		vdEMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vdSize), 1))
		vdEEW := c.SEW * SEW(vdSize)
		vs2EMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vs2Size), 1))
		vs2EEW := c.SEW * SEW(vs2Size)
		if vdEEW > SEW(i.Option.XLEN) || vs2EEW > SEW(i.Option.XLEN) {
			continue
		}

		vd := int(vdEMUL1)
		vs2 := vd * 2

		builder.WriteString(i.gWriteRandomData(vdEMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, vdEMUL1, SEW(8)))

		builder.WriteString(i.gWriteTestData(false, vs2EMUL1, vs2EEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, vs2EMUL1, vs2EEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, v%d%s\n",
			i.Name, vd, vs2, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, vdEMUL1, vdEEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}

	return res
}
