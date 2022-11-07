package generator

import (
	"encoding/binary"
	"fmt"
)

func (i *insn) gTestDataAddr() string {
	return fmt.Sprintf("la a0, testdata\n")
}

func (i *insn) gWriteRandomData(lmul LMUL) string {
	nBytes := i.vlenb() * int(lmul)
	rdata := genRandomData(int64(nBytes))

	res := "# Write random data into test data area.\n"
	for a := 0; a < nBytes/8; a++ {
		elem := binary.LittleEndian.Uint64(rdata)
		rdata = rdata[8:]

		res += fmt.Sprintf("li a1, 0x%x\n", elem)
		res += fmt.Sprintf("sd a1, %d(a0)\n", a*8)
	}
	return res + "\n"
}

func (i *insn) gWriteTestData(lmul LMUL, sew SEW, idx int) string {
	nBytes := i.vlenb() * int(lmul)
	res := fmt.Sprintf("# Write test data into test data area.\n")
	cases := i.testCases(sew)
	for a := 0; a < (nBytes / (int(sew) / 8)); a++ {
		b := a % len(cases)
		switch sew {
		case 8:
			res += fmt.Sprintf("li a1, 0x%x\n", convNum[uint8](cases[b][idx]))
			res += fmt.Sprintf("sb a1, %d(a0)\n", a*(int(sew)/8))
		case 16:
			res += fmt.Sprintf("li a1, 0x%x\n", convNum[uint16](cases[b][idx]))
			res += fmt.Sprintf("sh a1, %d(a0)\n", a*(int(sew)/8))
		case 32:
			res += fmt.Sprintf("li a1, 0x%x\n", convNum[uint32](cases[b][idx]))
			res += fmt.Sprintf("sw a1, %d(a0)\n", a*(int(sew)/8))
		case 64:
			res += fmt.Sprintf("li a1, 0x%x\n", convNum[uint64](cases[b][idx]))
			res += fmt.Sprintf("sd a1, %d(a0)\n", a*(int(sew)/8))
		}
	}
	return res + "\n"
}

func (i *insn) gLoadDataIntoRegisterGroup(
	group int, lmul LMUL, sew SEW) string {
	res := fmt.Sprintf("\n# Load data into v%d register group.\n", group)
	res += "li t0, -1\n"
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n", sew.String(), lmul.String())
	res += fmt.Sprintf("vle%d.v v%d, (a0)\n\n", sew, group)
	return res
}

func (i *insn) gStoreRegisterGroupIntoData(
	group int, lmul LMUL, sew SEW) string {
	res := fmt.Sprintf("\n# Store v%d register group into test data area.\n", group)
	res += "li t0, -1\n"
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
		sew.String(), lmul.String())
	res += fmt.Sprintf("vse%d.v v%d, (a0)\n\n", sew, group)
	return res
}

func (i *insn) gMagicInsn(group int) string {
	return fmt.Sprintf("addi x0, x%d, %d\n\n", 1*int(group), 2*int(group))
}

func (i *insn) gVsetvli(vl int, sew SEW, lmul LMUL) string {
	res := fmt.Sprintf("li t0, %d\n", vl)
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
		sew.String(), lmul.String())
	return res
}
