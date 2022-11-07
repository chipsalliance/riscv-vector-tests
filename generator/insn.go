package generator

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
	"math"
)

type insnFormat string

const (
	insnFormatVdRs1mVm    insnFormat = "vd,(rs1),vm"
	insnFormatVs2Rs1mVm   insnFormat = "vs3,(rs1),vm"
	insnFormatVdRs1m      insnFormat = "vd,(rs1)"
	insnFormatVdRs1mRs2Vm insnFormat = "vd,(rs1),rs2,vm"
	insnFormatVdRs1mVs2Vm insnFormat = "vd,(rs1),vs2,vm"
	insnFormatVdVs2Vs1    insnFormat = "vd,vs2,vs1"
	insnFormatVdVs2Vs1V0  insnFormat = "vd,vs2,vs1,v0"
	insnFormatVdVs2Vs1Vm  insnFormat = "vd,vs2,vs1,vm"
	insnFormatVdVs2Rs1V0  insnFormat = "vd,vs2,rs1,v0"
	insnFormatVdVs2Fs1V0  insnFormat = "vd,vs2,fs1,v0"
	insnFormatVdVs2Rs1Vm  insnFormat = "vd,vs2,rs1,vm"
	insnFormatVdVs2Fs1Vm  insnFormat = "vd,vs2,fs1,vm"
	insnFormatVdVs2ImmV0  insnFormat = "vd,vs2,imm,v0"
	insnFormatVdVs2ImmVm  insnFormat = "vd,vs2,imm,vm"
	insnFormatVdVs2UimmVm insnFormat = "vd,vs2,uimm,vm"
	insnFormatVdVs1Vs2Vm  insnFormat = "vd,vs1,vs2,vm"
	insnFormatVdRs1Vs2Vm  insnFormat = "vd,rs1,vs2,vm"
	insnFormatVdFs1Vs2Vm  insnFormat = "vd,fs1,vs2,vm"
	insnFormatVdVs1       insnFormat = "vd,vs1"
	insnFormatVdRs1       insnFormat = "vd,rs1"
	insnFormatVdFs1       insnFormat = "vd,fs1"
	insnFormatVdImm       insnFormat = "vd,imm"
	insnFormatVdVs2       insnFormat = "vd,vs2"
	insnFormatVdVs2Vm     insnFormat = "vd,vs2,vm"
	insnFormatRdVs2Vm     insnFormat = "rd,vs2,vm"
	insnFormatRdVs2       insnFormat = "rd,vs2"
	insnFormatFdVs2       insnFormat = "fd,vs2"
	insnFormatVdVm        insnFormat = "vd,vm"
)

var formats = map[insnFormat]struct{}{
	insnFormatVdRs1mVm:    {},
	insnFormatVs2Rs1mVm:   {},
	insnFormatVdRs1m:      {},
	insnFormatVdRs1mRs2Vm: {},
	insnFormatVdRs1mVs2Vm: {},
	insnFormatVdVs2Vs1:    {},
	insnFormatVdVs2Vs1V0:  {},
	insnFormatVdVs2Vs1Vm:  {},
	insnFormatVdVs2Rs1V0:  {},
	insnFormatVdVs2Fs1V0:  {},
	insnFormatVdVs2Rs1Vm:  {},
	insnFormatVdVs2Fs1Vm:  {},
	insnFormatVdVs2ImmV0:  {},
	insnFormatVdVs2ImmVm:  {},
	insnFormatVdVs2UimmVm: {},
	insnFormatVdVs1Vs2Vm:  {},
	insnFormatVdRs1Vs2Vm:  {},
	insnFormatVdFs1Vs2Vm:  {},
	insnFormatVdVs1:       {},
	insnFormatVdRs1:       {},
	insnFormatVdFs1:       {},
	insnFormatVdImm:       {},
	insnFormatVdVs2:       {},
	insnFormatVdVs2Vm:     {},
	insnFormatRdVs2Vm:     {},
	insnFormatRdVs2:       {},
	insnFormatFdVs2:       {},
	insnFormatVdVm:        {},
}

type Option struct {
	VLEN VLEN
	ELEN ELEN
}

type insn struct {
	Name   string     `toml:"name"`
	Format insnFormat `toml:"format"`
	Tests  tests      `toml:"tests"`
	Option Option     `toml:"-"`
}

func ReadInsnFromToml(contents []byte, option Option) (*insn, error) {
	i := insn{Option: option}

	if err := i.check(); err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(contents, &i); err != nil {
		return nil, err
	}

	if err := i.Tests.initialize(); err != nil {
		return nil, err
	}

	if _, ok := formats[i.Format]; !ok {
		return nil, errors.New("invalid test format")
	}

	return &i, nil
}

func (i *insn) check() error {
	if !i.Option.VLEN.Valid() {
		return fmt.Errorf("wrong VLEN: %d", i.Option.VLEN)
	}

	if !i.Option.ELEN.Valid(i.Option.VLEN) {
		return fmt.Errorf("wrong ELEN: %d", i.Option.ELEN)
	}
	return nil
}

func (i *insn) Generate() []byte {
	buf := make([]byte, 0)
	buf = i.genHeader(buf)
	buf = i.genCode(buf)
	buf = i.genData(buf)
	return buf
}

func (i *insn) genHeader(buf []byte) []byte {
	buf = append(buf, []byte(fmt.Sprintf(`
# This file is automatically generated. Do not edit.
# Instruction: %s

#include "riscv_test.h"
#include "test_macros.h"

RVTEST_RV64UV
`, i.Name))...)
	return buf
}

func (i *insn) genCode(buf []byte) []byte {
	buf = append(buf, []byte(fmt.Sprintf(`
RVTEST_CODE_BEGIN

%s

  TEST_CASE(2, x0, 0x0)
  TEST_PASSFAIL

RVTEST_CODE_END
`, i.genTestCases()))...)
	return buf
}

func (i *insn) genData(buf []byte) []byte {
	buf = append(buf, []byte(fmt.Sprintf(`
  .data
RVTEST_DATA_BEGIN

# Reserve space for test data.
testdata:
  .zero %d

RVTEST_DATA_END
`, i.vlenb()*(8 /* max LMUL */)))...)
	return buf
}

func (i *insn) genTestCases() string {
	switch i.Format {
	case insnFormatVdVs2Vs1Vm:
		return i.genCodeVdVs2Vs1Vm()
	case insnFormatVdRs1mVm:
		return i.genCodeVdRs1mVm()
	default:
		log.Fatalln("unreachable")
		return ""
	}
}

func (i *insn) vlenb() int {
	return int(i.Option.VLEN) / 8
}

func (i *insn) testCases(sew SEW) [][]any {
	res := make([][]any, 0)
	for _, c := range i.Tests.Base {
		l := make([]any, len(c))
		for b, op := range c {
			l[b] = op
		}
		res = append(res, l)
	}

	switch sew {
	case 8:
		for _, c := range i.Tests.SEW8 {
			l := make([]any, len(c))
			for b, op := range c {
				l[b] = op
			}
			res = append(res, l)
		}
	case 16:
		for _, c := range i.Tests.SEW16 {
			l := make([]any, len(c))
			for b, op := range c {
				l[b] = op
			}
			res = append(res, l)
		}
	case 32:
		for _, c := range i.Tests.SEW32 {
			l := make([]any, len(c))
			for b, op := range c {
				l[b] = op
			}
			res = append(res, l)
		}
	case 64:
		for _, c := range i.Tests.SEW64 {
			l := make([]any, len(c))
			for b, op := range c {
				l[b] = op
			}
			res = append(res, l)
		}
	}

	return res
}

type combination struct {
	SEW   SEW
	LMUL  LMUL
	LMUL1 LMUL
	Vl    int
	Mask  bool
}

func (c *combination) comment() string {
	return fmt.Sprintf(
		"\n\n# Generating tests for LMUL: %s, SEW or EEW: %s, Mask: %v\n\n",
		c.LMUL.String(),
		c.SEW.String(),
		c.Mask)
}

func (i *insn) combinations(sews []SEW) []combination {
	res := make([]combination, 0)
	for _, lmul := range allLMULs {
		for _, sew := range sews {
			if float64(lmul) < float64(sew)/float64(i.Option.ELEN) {
				continue
			}
			lmul1 := LMUL(math.Max(float64(lmul), 1))
			for _, mask := range []bool{false, true} {
				vlmax1 := int((float64(i.Option.VLEN) / float64(sew)) * float64(lmul1))
				for _, vl := range []int{0, vlmax1 / 2, vlmax1, vlmax1 + 1} {
					res = append(res, combination{
						SEW:   sew,
						LMUL:  lmul,
						LMUL1: lmul1,
						Vl:    vl,
						Mask:  mask,
					})
				}
			}
		}
	}

	return res
}

func (i *insn) genCodeTestDataAddr() string {
	return "la a0, testdata\n\n"
}

func (i *insn) genCodeWriteRandomData(lmul LMUL) string {
	nBytes := i.vlenb() * int(lmul)
	rdata := genRandomData(int64(nBytes))

	res := "# Write random data into test data area.\n"
	for a := 0; a < nBytes/8; a++ {
		elem := binary.LittleEndian.Uint64(rdata)
		rdata = rdata[8:]

		res += fmt.Sprintf("li a1, 0x%x\n", elem)
		res += fmt.Sprintf("sd a1, %d(a0)\n", a*8)
	}
	return res
}

func (i *insn) genCodeWriteTestData(lmul LMUL, sew SEW, idx int) string {
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
	return res
}

func (i *insn) genCodeLoadDataIntoRegisterGroup(
	group int, lmul LMUL, sew SEW) string {
	res := fmt.Sprintf("\n# Load data into v%d register group.\n", group)
	res += "li t0, -1\n"
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n", sew.String(), lmul.String())
	res += fmt.Sprintf("vle%d.v v%d, (a0)\n\n", sew, group)
	return res
}

func (i *insn) genCodeStoreRegisterGroupIntoDataArea(
	group int, lmul LMUL, sew SEW) string {
	res := fmt.Sprintf("\n# Store v%d register group into test data area.\n", group)
	res += "li t0, -1\n"
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
		sew.String(), lmul.String())
	res += fmt.Sprintf("vse%d.v v%d, (a0)\n\n", sew, group)
	return res
}

func (i *insn) genCodeMagicInsn(group int) string {
	return fmt.Sprintf("addi x0, x%d, %d", 1*int(group), 2*int(group))
}

func (i *insn) genCodeVsetvli(vl int, sew SEW, lmul LMUL) string {
	res := fmt.Sprintf("li t0, %d\n", vl)
	res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
		sew.String(), lmul.String())
	return res
}