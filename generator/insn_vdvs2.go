package generator

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (i *Insn) genCodeVdVs2() []string {
	s := regexp.MustCompile(`vmv(\d)r.v`)
	nr, err := strconv.Atoi(s.FindStringSubmatch(i.Name)[1])
	if err != nil {
		log.Fatal("unreachable")
	}

	combinations := i.combinations([]LMUL{LMUL(nr)}, allSEWs, []bool{false})
	res := make([]string, 0, len(combinations))

	for _, c := range combinations {
		builder := strings.Builder{}
		builder.WriteString(i.gTestDataAddr())
		builder.WriteString(c.comment())

		vd := int(c.LMUL)
		vs2 := 2 * int(c.LMUL)
		builder.WriteString(i.gWriteRandomData(c.LMUL * 2))

		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL, c.SEW))
		builder.WriteString(fmt.Sprintf("addi a0, a0, %d\n", int(c.LMUL)*i.vlenb()))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL, c.SEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, v%d\n",
			i.Name, vd, vs2))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gStoreRegisterGroupIntoData(vd, c.LMUL, c.SEW))
		builder.WriteString(i.gMagicInsn(vd))

		res = append(res, builder.String())

	}

	return res
}
