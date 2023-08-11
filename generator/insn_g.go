package generator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

func (i *Insn) gResultDataAddr() string {
	return fmt.Sprintf("la a0, resultdata\n")
}

func (i *Insn) gWriteRandomData(lmul LMUL) string {
	nBytes := i.vlenb() * int(lmul)
	off := i.TestData.Append(genRandomData(int64(nBytes)))

	builder := strings.Builder{}
	builder.WriteString("# Move a0 to test data area.\n")
	builder.WriteString("la a0, testdata\n")
	builder.WriteString(fmt.Sprintf("li a5, %d\n", off))
	builder.WriteString("add a0, a0, a5\n")
	return builder.String()
}

func (i *Insn) gWriteIndexData(lmul1 LMUL, n int, sew SEW) string {
	if n <= 0 {
		n = 1
	}
	nBytes := i.vlenb() * int(lmul1)
	builder := strings.Builder{}
	buf := &bytes.Buffer{}
	s := genShuffledSlice(n)
	for a := 0; a < nBytes/(int(sew)/8); a++ {
		switch sew {
		case 8:
			_ = binary.Write(buf, binary.LittleEndian, uint8(s[a%n]))
		case 16:
			_ = binary.Write(buf, binary.LittleEndian, uint16(s[a%n]*2))
		case 32:
			_ = binary.Write(buf, binary.LittleEndian, uint32(s[a%n]*4))
		case 64:
			_ = binary.Write(buf, binary.LittleEndian, uint64(s[a%n]*8))
		}
	}
	off := i.TestData.Append(buf.Bytes())
	builder.WriteString("# Move a0 to test data area.\n")
	builder.WriteString("la a0, testdata\n")
	builder.WriteString(fmt.Sprintf("li a5, %d\n", off))
	builder.WriteString("add a0, a0, a5\n")

	return builder.String()
}

func (i *Insn) gWriteIntegerTestData(lmul LMUL, sew SEW, idx int) string {
	return i.gWriteTestData(false, lmul, sew, idx)
}

func (i *Insn) gWriteTestData(float bool, lmul LMUL, sew SEW, idx int) string {
	nBytes := i.vlenb() * int(lmul)
	cases := i.testCases(float, sew)
	builder := strings.Builder{}
	buf := &bytes.Buffer{}
	for a := 0; a < (nBytes / (int(sew) / 8)); a++ {
		b := a % len(cases)
		switch sew {
		case 8:
			_ = binary.Write(buf, binary.LittleEndian, convNum[uint8](cases[b][idx]))
		case 16:
			_ = binary.Write(buf, binary.LittleEndian, convNum[uint16](cases[b][idx]))
		case 32:
			_ = binary.Write(buf, binary.LittleEndian, convNum[uint32](cases[b][idx]))
		case 64:
			_ = binary.Write(buf, binary.LittleEndian, convNum[uint64](cases[b][idx]))
		}
	}
	off := i.TestData.Append(buf.Bytes())
	builder.WriteString("# Move a0 to test data area.\n")
	builder.WriteString("la a0, testdata\n")
	builder.WriteString(fmt.Sprintf("li a5, %d\n", off))
	builder.WriteString("add a0, a0, a5\n")

	return builder.String()
}

func (i *Insn) gLoadDataIntoRegisterGroup(
	group int, lmul LMUL, sew SEW) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("\n# Load data into v%d register group.\n", group))
	builder.WriteString("li t0, -1\n")
	builder.WriteString(fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n", sew.String(), lmul.String()))
	builder.WriteString(fmt.Sprintf("vle%d.v v%d, (a0)\n\n", sew, group))
	return builder.String()
}

func (i *Insn) gStoreRegisterGroupIntoResultData(
	group int, lmul LMUL, sew SEW) string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\n# Store v%d register group into result data area.\n", group))
	builder.WriteString("li t0, -1\n")
	builder.WriteString(fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
		sew.String(), lmul.String()))
	builder.WriteString(fmt.Sprintf("vse%d.v v%d, (a0)\n\n", sew, group))
	return builder.String()
}

func (i *Insn) gMoveScalarToVector(scalar string, vector int, sew SEW) string {
	float := strings.HasPrefix(scalar, "f")
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\n# Move %s to the elem 0 of v%d\n", scalar, vector))

	builder.WriteString("li t0, -1\n")
	builder.WriteString(fmt.Sprintf("vsetvli t1, t0, %s,m1,ta,ma\n", sew.String()))
	builder.WriteString(fmt.Sprintf("v%smv.s.%s v%d, %s\n",
		iff(float, "f", ""), iff(float, "f", "x"), vector, scalar))

	return builder.String()
}

func (i *Insn) gMagicInsn(group int) string {
	insn := 0b0001011 + (group & 0b11111) << 15
	return fmt.Sprintf(".word 0x%x\n", insn);
}

func (i *Insn) gVsetvli(vl int, sew SEW, lmul LMUL) string {
	res := fmt.Sprintf("li t0, %d\n", vl)
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
		sew.String(), lmul.String())
	return res
}
