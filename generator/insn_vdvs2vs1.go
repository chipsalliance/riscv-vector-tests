package generator

import (
	"fmt"
	"strings"
)

func (i *Insn) genCodeVdVs2Vs1(pos int) []string {
    vsm3Insn := strings.HasPrefix(i.Name, "vsm3")
	sew32OnlyInsn := strings.HasPrefix(i.Name, "vg") || strings.HasPrefix(i.Name, "vsha")
	sews := iff(sew32OnlyInsn || vsm3Insn, []SEW{32}, allSEWs)
	combinations := i.combinations(
		iff(vsm3Insn, []LMUL {1, 2, 4, 8}, allLMULs),
		sews,
		[]bool{false},
		i.vxrms(),
	)
	res := make([]string, 0, len(combinations))
	for _, c := range combinations[pos:] {
		if sew32OnlyInsn && c.Vl % 4 != 0 {
			c.Vl = (c.Vl + 3) &^ 3
		}
		if vsm3Insn {
			c.Vl = (c.Vl + 7) &^ 7 
		}
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		vd := int(c.LMUL1)
		vss := []int{2 * int(c.LMUL1), 3 * int(c.LMUL1)}
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8)))

		for idx, vs := range vss {
			builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, idx))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs, c.LMUL1, c.SEW))
		}

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, v%d, v%d\n",
			i.Name, vd, vss[1], vss[0]))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd, c.LMUL1))

		res = append(res, builder.String())
	}

	return res
}
