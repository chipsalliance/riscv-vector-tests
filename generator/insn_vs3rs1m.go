package generator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (i *Insn) genCodeVs3Rs1m(pos int) []string {
	s := regexp.MustCompile(`vs(\d)r\.v`)
	subs := s.FindStringSubmatch(i.Name)

	lmuls := []LMUL{1}
	if len(subs) == 2 {
		nfields, _ := strconv.Atoi(subs[1])
		lmuls = []LMUL{LMUL(nfields)}
	}

	combinations := i.combinations(lmuls, []SEW{8}, []bool{false}, i.rms())
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		vs3, _, _ := getVRegs(c.LMUL1, true, i.Name)
		builder.WriteString(i.gWriteIntegerTestData(c.LMUL1, c.SEW, 0))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, (a0)\n", i.Name, vs3))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs3, c.LMUL1, c.SEW))
		builder.WriteString(i.gMagicInsn(vs3, c.LMUL1))

		res = append(res, builder.String())
	}
	return res
}
