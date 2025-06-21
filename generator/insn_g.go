package generator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"github.com/ksco/riscv-vector-tests/testfloat3"
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

func (i *Insn) gWriteIndexData(dataLmul1 LMUL, offsetLmul1 LMUL, n int, dataSew SEW, offsetSew SEW) string {
	if n <= 0 {
		n = 1
	}
	nBytes := i.vlenb() * int(offsetLmul1)
	builder := strings.Builder{}
	buf := &bytes.Buffer{}
	n = int(math.Min(float64(n), float64(int(dataLmul1)*i.vlenb()/(int(dataSew)/8))))
	s := genShuffledSlice(n)
	for a := 0; a < nBytes/(int(offsetSew)/8); a++ {
		switch offsetSew {
		case 8:
			_ = binary.Write(buf, binary.LittleEndian, uint8(s[a%n]*int(dataSew)/8))
		case 16:
			_ = binary.Write(buf, binary.LittleEndian, uint16(s[a%n]*int(dataSew)/8))
		case 32:
			_ = binary.Write(buf, binary.LittleEndian, uint32(s[a%n]*int(dataSew)/8))
		case 64:
			_ = binary.Write(buf, binary.LittleEndian, uint64(s[a%n]*int(dataSew)/8))
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
	return i.gWriteTestData(false, false, false, lmul, sew, idx, 0)
}

func (i *Insn) gWriteTestData(float bool, testfloat bool, cont bool, lmul LMUL, sew SEW, idx int, numops int) string {
	if !float && cont {
		panic("unreachable")
	}

	if float && testfloat && !cont {
		testfloat3.InitF16(numops)
		testfloat3.InitF32(numops)
		testfloat3.InitF64(numops)
	}
	nextf16 := func() uint16 {
		return testfloat3.GenF16(numops)[idx]
	}
	nextf32 := func() float32 {
		return testfloat3.GenF32(numops)[idx]
	}
	nextf64 := func() float64 {
		return testfloat3.GenF64(numops)[idx]
	}

	nBytes := i.vlenb() * int(lmul)
	cases := i.testCases(float, sew)
	builder := strings.Builder{}
	buf := &bytes.Buffer{}
	for a := 0; a < (nBytes / (int(sew) / 8)); a++ {
		switch sew {
		case 8:
			b := a % len(cases)
			_ = binary.Write(buf, binary.LittleEndian, convNum[uint8](cases[b][idx]))
		case 16:
			if (float && testfloat && a >= len(cases)) || cont {
				_ = binary.Write(buf, binary.LittleEndian, nextf16())
			} else {
				b := a % len(cases)
				_ = binary.Write(buf, binary.LittleEndian, convNum[uint16](cases[b][idx]))
			}
		case 32:
			// Manual test cases exhausted, use testfloat3 to generate new ones.
			if (float && testfloat && a >= len(cases)) || cont {
				_ = binary.Write(buf, binary.LittleEndian, math.Float32bits(nextf32()))
			} else {
				b := a % len(cases)
				_ = binary.Write(buf, binary.LittleEndian, convNum[uint32](cases[b][idx]))
			}
		case 64:
			// Manual test cases exhausted, use testfloat3 to generate new ones.
			if (float && testfloat && a >= len(cases)) || cont {
				_ = binary.Write(buf, binary.LittleEndian, math.Float64bits(nextf64()))
			} else {
				b := a % len(cases)
				_ = binary.Write(buf, binary.LittleEndian, convNum[uint64](cases[b][idx]))
			}
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
	builder.WriteString(fmt.Sprintf("vsetvli t1, x0, %s,%s,tu,mu\n", sew.String(), lmul.String()))
	builder.WriteString(fmt.Sprintf("vle%d.v v%d, (a0)\n\n", sew, group))
	return builder.String()
}

func (i *Insn) gStoreRegisterGroupIntoResultData(
	group int, lmul LMUL, sew SEW) string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\n# Store v%d register group into result data area.\n", group))
	builder.WriteString(fmt.Sprintf("vsetvli t1, x0, %s,%s,tu,mu\n",
		sew.String(), lmul.String()))
	builder.WriteString(fmt.Sprintf("vse%d.v v%d, (a0)\n\n", sew, group))
	return builder.String()
}

func (i *Insn) gMoveScalarToVector(scalar string, vector int, sew SEW) string {
	float := strings.HasPrefix(scalar, "f")
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\n# Move %s to the elem 0 of v%d\n", scalar, vector))

	builder.WriteString(fmt.Sprintf("vsetvli t1, x0, %s,m1,tu,mu\n", sew.String()))
	builder.WriteString(fmt.Sprintf("v%smv.s.%s v%d, %s\n",
		iff(float, "f", ""), iff(float, "f", "x"), vector, scalar))

	return builder.String()
}

func (i *Insn) gMagicInsn(group int, lmul1 LMUL) string {

	// opcode
	insn := 0b0001011

	// rs1 for vreg group
	insn += (group & 0b11111) << 15

	// rs2[0] for vxsat CSR
	if i.Vxsat {
		insn += 1 << 20
	}
	// rs2[4:1] for EMUL
	insn += int(lmul1) << 21

	return fmt.Sprintf(".word 0x%x\n", insn)
}

func (i *Insn) gVsetvli(vl int, sew SEW, lmul LMUL) string {
	res := fmt.Sprintf("li t0, %d\n", vl)
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,tu,mu\n",
		sew.String(), lmul.String())
	return res
}
