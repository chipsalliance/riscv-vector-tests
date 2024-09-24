package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVdVs1Vs2Vm(pos int) []string {
	float := strings.HasPrefix(i.Name, "vf")
	vdWidening := strings.HasPrefix(i.Name, "vw") || strings.HasPrefix(i.Name, "vfw")
	vdSize := iff(vdWidening, 2, 1)
	vs1Size := 1

	sews := iff(float, i.floatSEWs(), allSEWs)
	sews = iff(vdWidening, sews[:len(sews)-1], sews)
	combinations := i.combinations(
		iff(vdWidening, wideningMULs, allLMULs),
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
		vs1EMUL1 := LMUL(math.Max(float64(int(c.LMUL)*vs1Size), 1))
		vdEEW := c.SEW * SEW(vdSize)
		vs1EEW := c.SEW * SEW(vs1Size)
		if vdEEW > SEW(i.Option.XLEN) || vs1EEW > SEW(i.Option.XLEN) {
			res = append(res, "")
			continue
		}

		vd := int(vdEMUL1)
		vss := []int{
			vd * 2,
			vd*2 + int(vs1EMUL1),
		}

		for r := 0; r < i.Option.Repeat; r += 1 {
			builder.WriteString(i.gWriteTestData(float, !i.NoTestfloat3, r != 0, vdEMUL1, vdEEW, 0, 3))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, vdEMUL1, vdEEW))

			builder.WriteString(i.gWriteTestData(float, !i.NoTestfloat3, r != 0, vs1EMUL1, vs1EEW, 1, 3))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vss[1], vs1EMUL1, vs1EEW))

			builder.WriteString(i.gWriteTestData(float, !i.NoTestfloat3, r != 0, c.LMUL1, c.SEW, 2, 3))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vss[0], c.LMUL1, c.SEW))

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
