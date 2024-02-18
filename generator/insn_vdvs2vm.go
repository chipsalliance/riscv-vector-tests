package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdVs2Vm(pos int) []string {
	vdWidening := strings.HasPrefix(i.Name, "vfw")
	vdNarrowing := strings.HasPrefix(i.Name, "vfn")
	vdSize := iff(vdWidening, 2, 1)
	vs2Size := iff(vdNarrowing, 2, 1)

	lmuls := iff(vdWidening || vdNarrowing, wideningMULs, allLMULs)
	sews := iff(vdWidening || vdNarrowing, floatSEWs[:len(floatSEWs)-1], floatSEWs)
	combinations := i.combinations(lmuls, sews, []bool{false, true}, i.vxrms())

	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, c.LMUL1, SEW(8)))

		vdEMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vdSize), 1))
		vdEEW := c.SEW * SEW(vdSize)
		vs2EMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vs2Size), 1))
		vs2EEW := c.SEW * SEW(vs2Size)
		if vdEEW > SEW(i.Option.XLEN) || vs2EEW > SEW(i.Option.XLEN) {
			res = append(res, "")
			continue
		}

		vd := int(vdEMUL1)
		vs2 := vd * 2

		builder.WriteString(i.gWriteRandomData(vdEMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, vdEMUL1, SEW(8)))

		builder.WriteString(i.gWriteTestData(false, false, vs2EMUL1, vs2EEW, 0, 1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, vs2EMUL1, vs2EEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, v%d%s\n",
			i.Name, vd, vs2, v0t(c.Mask)))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, vdEMUL1, vdEEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}

	return res
}
