package generator

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *insn) genCodeVdRs1mVm() string {
	getEEW := func(name string) SEW {
		eew, _ := strconv.Atoi(
			strings.TrimSuffix(strings.TrimPrefix(i.Name, "vle"), ".v"))
		return SEW(eew)
	}
	res := i.genCodeTestDataAddr()
	res += i.genCodeWriteRandomData(LMUL(1))
	res += i.genCodeLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8))

	for _, c := range i.combinations([]SEW{getEEW(i.Name)}) {
		res += c.comment()

		vd := 1 * int(c.LMUL1)

		res += i.genCodeWriteRandomData(c.LMUL1)
		res += i.genCodeLoadDataIntoRegisterGroup(vd, c.LMUL1, c.SEW)
		res += i.genCodeWriteTestData(c.LMUL1, c.SEW, 0)

		res += "# -------------- TEST BEGIN --------------\n"
		res += i.genCodeVsetvli(c.Vl, c.SEW, c.LMUL)
		res += fmt.Sprintf("%s v%d, (a0)%s\n", i.Name, vd, v0t(c.Mask))
		res += "# -------------- TEST END   --------------\n"

		res += i.genCodeStoreRegisterGroupIntoDataArea(vd, c.LMUL1, c.SEW)
		res += i.genCodeMagicInsn(int(c.LMUL1))
	}
	return res
}
