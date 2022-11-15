package generator

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func (i *insn) gTestDataAddr() string {
	return fmt.Sprintf("la a0, testdata\n")
}

func (i *insn) gWriteRandomData(lmul LMUL) string {
	nBytes := i.vlenb() * int(lmul)
	rdata := genRandomData(int64(nBytes))

	builder := strings.Builder{}
	builder.WriteString("# Write random data into test data area.\n")
	builder.WriteString("mv a3, a0\n")
	for a := 0; a < nBytes/8; a++ {
		elem := binary.LittleEndian.Uint64(rdata)
		rdata = rdata[8:]

		builder.WriteString(fmt.Sprintf("li a1, 0x%x\n", elem))
		builder.WriteString(fmt.Sprintf("sd a1, 0(a3)\n"))
		builder.WriteString("addi a3, a3, 8\n")
	}
	builder.WriteString("\n")

	return builder.String()
}

func (i *insn) gWriteIntegerTestData(lmul LMUL, sew SEW, idx int) string {
	return i.gWriteTestData(false, lmul, sew, idx)
}

func (i *insn) gWriteTestData(float bool, lmul LMUL, sew SEW, idx int) string {
	nBytes := i.vlenb() * int(lmul)
	cases := i.testCases(float, sew)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("# Write test data into test data area.\n"))
	builder.WriteString("mv a3, a0\n")
	for a := 0; a < (nBytes / (int(sew) / 8)); a++ {
		b := a % len(cases)
		switch sew {
		case 8:
			builder.WriteString(fmt.Sprintf("li a1, 0x%x\n", convNum[uint8](cases[b][idx])))
			builder.WriteString(fmt.Sprintf("sb a1, 0(a3)\n"))
			builder.WriteString(fmt.Sprintf("addi a3, a3, %d\n", int(sew)/8))
		case 16:
			builder.WriteString(fmt.Sprintf("li a1, 0x%x\n", convNum[uint16](cases[b][idx])))
			builder.WriteString(fmt.Sprintf("sh a1, 0(a3)\n"))
			builder.WriteString(fmt.Sprintf("addi a3, a3, %d\n", int(sew)/8))
		case 32:
			builder.WriteString(fmt.Sprintf("li a1, 0x%x\n", convNum[uint32](cases[b][idx])))
			builder.WriteString(fmt.Sprintf("sw a1, 0(a3)\n"))
			builder.WriteString(fmt.Sprintf("addi a3, a3, %d\n", int(sew)/8))
		case 64:
			builder.WriteString(fmt.Sprintf("li a1, 0x%x\n", convNum[uint64](cases[b][idx])))
			builder.WriteString(fmt.Sprintf("sd a1, 0(a3)\n"))
			builder.WriteString(fmt.Sprintf("addi a3, a3, %d\n", int(sew)/8))
		}
	}

	builder.WriteString("\n")
	return builder.String()
}

func (i *insn) gLoadDataIntoRegisterGroup(
	group int, lmul LMUL, sew SEW) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("\n# Load data into v%d register group.\n", group))
	builder.WriteString("li t0, -1\n")
	builder.WriteString(fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n", sew.String(), lmul.String()))
	builder.WriteString(fmt.Sprintf("vle%d.v v%d, (a0)\n\n", sew, group))
	return builder.String()
}

func (i *insn) gStoreRegisterGroupIntoData(
	group int, lmul LMUL, sew SEW) string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\n# Store v%d register group into test data area.\n", group))
	builder.WriteString("li t0, -1\n")
	builder.WriteString(fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
		sew.String(), lmul.String()))
	builder.WriteString(fmt.Sprintf("vse%d.v v%d, (a0)\n\n", sew, group))
	return builder.String()
}

func (i *insn) gMoveScalarToVector(scalar string, vector int, sew SEW) string {
	float := strings.HasPrefix(scalar, "f")
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\n# Move %s to the elem 0 of v%d\n", scalar, vector))

	builder.WriteString("li t0, -1\n")
	builder.WriteString(fmt.Sprintf("vsetvli t1, t0, %s,m1,ta,ma\n", sew.String()))
	builder.WriteString(fmt.Sprintf("v%smv.s.%s v%d, %s\n",
		iff(float, "f", ""), iff(float, "f", "x"), vector, scalar))

	return builder.String()
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
