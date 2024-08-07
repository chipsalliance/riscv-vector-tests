package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdVs2Vs1Vm(pos int) []string {
	float := strings.HasPrefix(i.Name, "vf") || strings.HasPrefix(i.Name, "vmf")
	sew64_insn := i.Name == "vclmul.vv" || i.Name == "vclmulh.vv"
	vdWidening := strings.HasPrefix(i.Name, "vw") || strings.HasPrefix(i.Name, "vfw")
	vs2Widening := strings.HasSuffix(i.Name, ".wv")
	vdSize := iff(vdWidening, 2, 1)
	vs2Size := iff(vs2Widening, 2, 1)

	sews := iff(float, floatSEWs, allSEWs)
	sews = iff(vdWidening || vs2Widening, sews[:len(sews)-1], sews)
    sews = iff(sew64_insn, []SEW{64}, sews)
	combinations := i.combinations(
		iff(vdWidening || vs2Widening, wideningMULs, iff(sew64_insn, []LMUL{1, 2, 4, 8}, allLMULs)),
		sews,
		[]bool{false, true},
		i.vxrms(),
	)
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		if strings.HasPrefix(i.Name, "vrgatherei16") && (16*c.LMUL/LMUL(c.SEW) > LMUL(8.0) || 16*c.LMUL/LMUL(c.SEW) < LMUL(0.125)) {
			res = append(res, "")
			continue
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
		vss := []int{
			vd * 2,
			vd*2 + int(vs2EMUL1),
		}

		if vdEMUL1 == vs2EMUL1 && !strings.HasPrefix(i.Name, "vrgatherei16") {
			vd1, vs1, vs2 := getVRegs(vdEMUL1, false, i.Name)
			vd = vd1
			vss = []int{vs1, vs2}
		}

		for r := 0; r < i.Option.Repeat; r += 1 {
			builder.WriteString(i.gWriteRandomData(vdEMUL1))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, vdEMUL1, SEW(8)))

			builder.WriteString(i.gWriteTestData(float, !i.NoTestfloat3, r != 0, c.LMUL1, c.SEW, 0, 2))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vss[0], c.LMUL1, c.SEW))

			builder.WriteString(i.gWriteTestData(float, !i.NoTestfloat3, r != 0, vs2EMUL1, vs2EEW, 1, 2))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vss[1], vs2EMUL1, vs2EEW))

			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))

			builder.WriteString(fmt.Sprintf("%s v%d, v%d, v%d%s\n",
				i.Name, vd, vss[1], vss[0], v0t(c.Mask)))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, vdEMUL1, vdEEW))
			builder.WriteString(i.gMagicInsn(vd, vdEMUL1))
		}

		res = append(res, builder.String())
	}

	return res
}
