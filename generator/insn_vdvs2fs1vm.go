package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdVs2Fs1Vm(pos int) []string {
	vdWidening := strings.HasPrefix(i.Name, "vfw")
	vs2Widening := strings.HasSuffix(i.Name, ".wf")
	vdSize := iff(vdWidening, 2, 1)
	vs2Size := iff(vs2Widening, 2, 1)

	sews := iff(vdWidening || vs2Widening, floatSEWs[:len(floatSEWs)-1], floatSEWs)
	combinations := i.combinations(
		iff(vdWidening || vs2Widening, wideningMULs, allLMULs),
		sews,
		[]bool{false, true},
		i.vxrms(),
	)
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		builder.WriteString(i.gWriteRandomData(LMUL(1)))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

		vdEMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vdSize), 1))
		vs2EMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vs2Size), 1))
		vdEEW := c.SEW * SEW(vdSize)
		vs2EEW := c.SEW * SEW(vs2Size)
		if vdEEW > SEW(i.Option.XLEN) {
			res = append(res, "")
			continue
		}

		vd := int(vdEMUL1)
		vs2 := vd * 2

		for r := 0; r < i.Option.Repeat; r += 1 {
			builder.WriteString(i.gWriteRandomData(vdEMUL1))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, vdEMUL1, SEW(8)))

			builder.WriteString(i.gWriteTestData(true, !i.NoTestfloat3, r != 0, vs2EMUL1, vs2EEW, 0, 1))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, vs2EMUL1, vs2EEW))

			cases := i.testCases(true, c.SEW)
			for a := 0; a < len(cases); a++ {
				builder.WriteString("# -------------- TEST BEGIN --------------\n")
				switch c.SEW {
				case 32:
					builder.WriteString(fmt.Sprintf("li s0, 0x%x\n", convNum[uint32](cases[a][0])))
					builder.WriteString(fmt.Sprintf("fmv.w.x f0, s0\n"))
				case 64:
					builder.WriteString(fmt.Sprintf("li s0, 0x%x\n", convNum[uint64](cases[a][0])))
					builder.WriteString(fmt.Sprintf("fmv.d.x f0, s0\n"))
				}
				builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
				builder.WriteString(fmt.Sprintf("%s v%d, v%d, f0%s\n",
					i.Name, vd, vs2, v0t(c.Mask)))
				builder.WriteString("# -------------- TEST END   --------------\n")

				builder.WriteString(i.gResultDataAddr())
				builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, vdEMUL1, vdEEW))
				builder.WriteString(i.gMagicInsn(vd))
			}
		}

		res = append(res, builder.String())
	}
	return res
}
