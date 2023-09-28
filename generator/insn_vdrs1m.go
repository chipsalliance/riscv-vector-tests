package generator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (i *Insn) genCodeVdRs1m(pos int) []string {
	s := regexp.MustCompile(`vl(\d)re(\d+)\.v`)
	subs := s.FindStringSubmatch(i.Name)

	lmuls := []LMUL{1}
	sews := []SEW{8}
	if len(subs) == 3 {
		nfields, _ := strconv.Atoi(subs[1])
		eew, _ := strconv.Atoi(subs[2])
		lmuls = []LMUL{LMUL(nfields)}
		sews = []SEW{SEW(eew)}
	}
	combinations := i.combinations(lmuls, sews, []bool{false})

	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.comment())

		vd := int(c.LMUL1)
		builder.WriteString(i.gWriteRandomData(c.LMUL1))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 0))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)\n", i.Name, vd))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())
	}
	return res
}
