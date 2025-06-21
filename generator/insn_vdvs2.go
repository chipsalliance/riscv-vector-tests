package generator

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (i *Insn) genCodeVdVs2(pos int) []string {
	sew32OnlyInsn := strings.HasPrefix(i.Name, "vg") || strings.HasPrefix(i.Name, "vsm4") || strings.HasPrefix(i.Name, "vaes")
	sews := iff(sew32OnlyInsn, []SEW{32}, allSEWs)

	if sew32OnlyInsn {
		i.EGW = 128
	}

	var nr int
	var err error

	if match := regexp.MustCompile(`vmv(\d+)r.v`).FindStringSubmatch(i.Name); len(match) > 1 {
		nr, err = strconv.Atoi(match[1])
		if err != nil {
			log.Fatalf("Error parsing register number: %v", err)
		}
	}

	var combinations []combination
	if sew32OnlyInsn {
		combinations = i.combinations(allLMULs, sews, []bool{false}, i.rms())
	} else {
		combinations = i.combinations([]LMUL{LMUL(nr)}, sews, []bool{false}, i.rms())
	}
	res := make([]string, 0, len(combinations))

	for _, c := range combinations[pos:] {
		if sew32OnlyInsn && c.Vl%4 != 0 {
			c.Vl = (c.Vl + 3) / 4 * 4
		}

		builder := strings.Builder{}
		builder.WriteString(c.initialize())

		var vd, vs2 int
		if sew32OnlyInsn {
			vd = int(c.LMUL1)
			vs2 = 3 * int(c.LMUL1)
		} else {
			vd, vs2, _ = getVRegs(c.LMUL, true, i.Name)
		}

		builder.WriteString(i.gWriteRandomData(c.LMUL * 2))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vd, c.LMUL, c.SEW))
		builder.WriteString(fmt.Sprintf("li t1, %d\n", int(c.LMUL)*i.vlenb()))
		builder.WriteString(fmt.Sprintf("add a0, a0, t1\n"))
		builder.WriteString(i.gLoadDataIntoRegisterGroup(vs2, c.LMUL, c.SEW))

		builder.WriteString("# -------------- TEST BEGIN --------------\n")
		builder.WriteString(i.gVsetvli(c.Vl, c.SEW, c.LMUL))
		builder.WriteString(fmt.Sprintf("%s v%d, v%d\n", i.Name, vd, vs2))
		builder.WriteString("# -------------- TEST END   --------------\n")

		builder.WriteString(i.gResultDataAddr())
		builder.WriteString(i.gStoreRegisterGroupIntoResultData(vd, c.LMUL, c.SEW))
		builder.WriteString(i.gMagicInsn(vd, c.LMUL))

		res = append(res, builder.String())
	}

	return res
}
