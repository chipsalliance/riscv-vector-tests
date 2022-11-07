package generator

import (
	"fmt"
)

func (i *insn) genCodeVdVs2Vs1Vm() string {
	res := i.genCodeTestDataAddr()
	res += i.genCodeWriteRandomData(LMUL(1))
	res += i.genCodeLoadDataIntoRegisterGroup(0, LMUL(1), SEW(8))

	for _, c := range i.combinations(allSEWs) {
		res += c.comment()

		vd := int(c.LMUL1)
		res += i.genCodeWriteRandomData(c.LMUL1)
		res += i.genCodeLoadDataIntoRegisterGroup(vd, c.LMUL1, SEW(8))

		for idx := 0; idx < 2; idx++ {
			res += i.genCodeWriteTestData(c.LMUL1, c.SEW, idx)
			res += i.genCodeLoadDataIntoRegisterGroup((idx+2)*int(c.LMUL1), c.LMUL1, c.SEW)
		}

		res += "# -------------- TEST BEGIN --------------\n"
		res += i.genCodeVsetvli(c.Vl, c.SEW, c.LMUL)

		res += fmt.Sprintf("%s v%d, v%d, v%d%s\n",
			i.Name, vd, 3*int(c.LMUL1), 2*int(c.LMUL1), v0t(c.Mask))
		res += "# -------------- TEST END   --------------\n"

		res += i.genCodeStoreRegisterGroupIntoDataArea(vd, c.LMUL1, c.SEW)
		res += i.genCodeMagicInsn(int(c.LMUL1))
	}
	return res
}
