package generator

import (
	"fmt"
	"math"
	"strings"
)

func (i *Insn) genCodeVs3Rs1mVs2Vm(pos int) []string {
	nfields := getNfieldsRoundedUp(i.Name)
	combinations := i.combinations(
		nfieldsLMULs(nfields),
		[]SEW{getEEW(i.Name)},
		[]bool{false, true})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}

		var sew SEW = 8 // sew is data width, c.SEW is offset width
		for ; sew <= SEW(i.Option.XLEN); sew <<= 1 {
			if int(sew) > int(float64(i.Option.XLEN)*float64(c.LMUL)) {
				continue
			}
			emul := (float64(c.SEW) / float64(sew)) * float64(c.LMUL)
			if emul > 8 || emul < 1./8 {
				continue
			}
			emul = math.Max(emul, 1)
			builder.WriteString(c.comment())

			builder.WriteString(i.gWriteRandomData(LMUL(1)))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(0, LMUL(1), SEW(32)))

			lmul1 := LMUL(math.Max(float64(c.LMUL)*float64(nfields), 1))
			vs3 := int(lmul1)
			vs1 := 2 * int(math.Max(emul, float64(int(c.LMUL1)*nfields)))
			builder.WriteString(i.gWriteIntegerTestData(lmul1, sew, 0))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, lmul1, sew))
			builder.WriteString(i.gWriteIndexData(lmul1, LMUL(emul), c.Vl, sew, c.SEW))
			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs1, LMUL(emul), c.SEW))

			builder.WriteString(i.gResultDataAddr())
			builder.WriteString("# -------------- TEST BEGIN --------------\n")
			builder.WriteString(i.gVsetvli(c.Vl, sew, c.LMUL))
			builder.WriteString(fmt.Sprintf("%s v%d, (a0), v%d%s\n", i.Name, vs3, vs1, v0t(c.Mask)))
			builder.WriteString("# -------------- TEST END   --------------\n")

			builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, lmul1, sew))
			builder.WriteString(i.gMagicInsn(vs3))
		}
		res = append(res, builder.String())
	}
	return res
}
