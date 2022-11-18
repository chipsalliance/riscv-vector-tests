package generator

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type insnFormat string

type Option struct {
	VLEN VLEN
	ELEN ELEN
}

const minStride = -1 // Must be negative
const maxStride = 3  // Must be greater than 1
const strides = maxStride - minStride

type Insn struct {
	Name   string     `toml:"name"`
	Format insnFormat `toml:"format"`
	Tests  tests      `toml:"tests"`
	Option Option     `toml:"-"`
}

const (
	insnFormatVdRs1mVm     insnFormat = "vd,(rs1),vm"      // Added
	insnFormatVs3Rs1mVm    insnFormat = "vs3,(rs1),vm"     // Added
	insnFormatVdRs1m       insnFormat = "vd,(rs1)"         // Added
	insnFormatVs3Rs1m      insnFormat = "vs3,(rs1)"        // Added
	insnFormatVdRs1mRs2Vm  insnFormat = "vd,(rs1),rs2,vm"  // Added
	insnFormatVs3Rs1mRs2Vm insnFormat = "vs3,(rs1),rs2,vm" // Added
	insnFormatVdRs1mVs2Vm  insnFormat = "vd,(rs1),vs2,vm"
	insnFormatVs3Rs1mVs2Vm insnFormat = "vs3,(rs1),vs2,vm"
	insnFormatVdVs2Vs1     insnFormat = "vd,vs2,vs1"    // Added
	insnFormatVdVs2Vs1V0   insnFormat = "vd,vs2,vs1,v0" // Added
	insnFormatVdVs2Vs1Vm   insnFormat = "vd,vs2,vs1,vm" // Added
	insnFormatVdVs2Rs1V0   insnFormat = "vd,vs2,rs1,v0" // Added
	insnFormatVdVs2Fs1V0   insnFormat = "vd,vs2,fs1,v0"
	insnFormatVdVs2Rs1Vm   insnFormat = "vd,vs2,rs1,vm"  // Added
	insnFormatVdVs2Fs1Vm   insnFormat = "vd,vs2,fs1,vm"  // Added
	insnFormatVdVs2ImmV0   insnFormat = "vd,vs2,imm,v0"  // Added
	insnFormatVdVs2ImmVm   insnFormat = "vd,vs2,imm,vm"  // Added
	insnFormatVdVs2UimmVm  insnFormat = "vd,vs2,uimm,vm" // Added
	insnFormatVdVs1Vs2Vm   insnFormat = "vd,vs1,vs2,vm"
	insnFormatVdRs1Vs2Vm   insnFormat = "vd,rs1,vs2,vm"
	insnFormatVdFs1Vs2Vm   insnFormat = "vd,fs1,vs2,vm" // Added
	insnFormatVdVs1        insnFormat = "vd,vs1"        // Added
	insnFormatVdRs1        insnFormat = "vd,rs1"        // Added
	insnFormatVdFs1        insnFormat = "vd,fs1"        // Added
	insnFormatVdImm        insnFormat = "vd,imm"        // Added
	insnFormatVdVs2        insnFormat = "vd,vs2"        // Added
	insnFormatVdVs2Vm      insnFormat = "vd,vs2,vm"     // Added
	insnFormatVdVs2VmP2    insnFormat = "vd,vs2,vm/2"   // Added
	insnFormatVdVs2VmP3    insnFormat = "vd,vs2,vm/3"   // Added
	insnFormatRdVs2Vm      insnFormat = "rd,vs2,vm"     // Added
	insnFormatRdVs2        insnFormat = "rd,vs2"        // Added
	insnFormatFdVs2        insnFormat = "fd,vs2"        // Added
	insnFormatVdVm         insnFormat = "vd,vm"         // Added
)

var formats = map[insnFormat]struct{}{
	insnFormatVdRs1mVm:     {},
	insnFormatVs3Rs1mVm:    {},
	insnFormatVdRs1m:       {},
	insnFormatVs3Rs1m:      {},
	insnFormatVdRs1mRs2Vm:  {},
	insnFormatVs3Rs1mRs2Vm: {},
	insnFormatVdRs1mVs2Vm:  {},
	insnFormatVs3Rs1mVs2Vm: {},
	insnFormatVdVs2Vs1:     {},
	insnFormatVdVs2Vs1V0:   {},
	insnFormatVdVs2Vs1Vm:   {},
	insnFormatVdVs2Rs1V0:   {},
	insnFormatVdVs2Fs1V0:   {},
	insnFormatVdVs2Rs1Vm:   {},
	insnFormatVdVs2Fs1Vm:   {},
	insnFormatVdVs2ImmV0:   {},
	insnFormatVdVs2ImmVm:   {},
	insnFormatVdVs2UimmVm:  {},
	insnFormatVdVs1Vs2Vm:   {},
	insnFormatVdRs1Vs2Vm:   {},
	insnFormatVdFs1Vs2Vm:   {},
	insnFormatVdVs1:        {},
	insnFormatVdRs1:        {},
	insnFormatVdFs1:        {},
	insnFormatVdImm:        {},
	insnFormatVdVs2:        {},
	insnFormatVdVs2Vm:      {},
	insnFormatVdVs2VmP2:    {},
	insnFormatVdVs2VmP3:    {},
	insnFormatRdVs2Vm:      {},
	insnFormatRdVs2:        {},
	insnFormatFdVs2:        {},
	insnFormatVdVm:         {},
}

func (i *Insn) genTestCases() []string {
	switch i.Format {
	case insnFormatVdRs1mVm:
		return i.genCodeVdRs1mVm()
	case insnFormatVs3Rs1mVm:
		return i.genCodeVs3Rs1mVm()
	case insnFormatVdRs1m:
		return i.genCodeVdRs1m()
	case insnFormatVs3Rs1m:
		return i.genCodeVs3Rs1m()
	case insnFormatVdRs1mRs2Vm:
		return i.genCodeVdRs1mRs2Vm()
	case insnFormatVs3Rs1mRs2Vm:
		return i.genCodeVs3Rs1mRs2Vm()
	case insnFormatVdVs2Vs1:
		return i.genCodeVdVs2Vs1()
	case insnFormatVdVs2Vs1V0:
		return i.genCodeVdVs2Vs1V0()
	case insnFormatVdVs2Rs1V0:
		return i.genCodeVdVs2Rs1V0()
	case insnFormatVdVs2Vs1Vm:
		return i.genCodeVdVs2Vs1Vm()
	case insnFormatVdVs2Rs1Vm:
		return i.genCodeVdVs2Rs1Vm()
	case insnFormatVdVs2Fs1Vm:
		return i.genCodeVdVs2Fs1Vm()
	case insnFormatVdVs2ImmV0:
		return i.genCodeVdVs2ImmV0()
	case insnFormatVdVs2ImmVm:
		return i.genCodeVdVs2ImmVm()
	case insnFormatVdVs2UimmVm:
		return i.genCodeVdVs2UimmVm()
	case insnFormatVdFs1Vs2Vm:
		return i.genCodeVdFs1Vs2Vm()
	case insnFormatVdVs1:
		return i.genCodeVdVs1()
	case insnFormatVdRs1:
		return i.genCodeVdRs1()
	case insnFormatVdFs1:
		return i.genCodeVdFs1()
	case insnFormatVdImm:
		return i.genCodeVdImm()
	case insnFormatVdVs2:
		return i.genCodeVdVs2()
	case insnFormatVdVs2Vm:
		return i.genCodeVdVs2Vm()
	case insnFormatVdVs2VmP2:
		return i.genCodeVdVs2VmP2()
	case insnFormatVdVs2VmP3:
		return i.genCodeVdVs2VmP3()
	case insnFormatRdVs2Vm:
		return i.genCodeRdVs2Vm()
	case insnFormatRdVs2:
		return i.genCodeRdVs2()
	case insnFormatFdVs2:
		return i.genCodeFdVs2()
	case insnFormatVdVm:
		return i.genCodeVdVm()
	default:
		log.Fatalln("unreachable")
		return nil
	}
}

func ReadInsnFromToml(contents []byte, option Option) (*Insn, error) {
	i := Insn{Option: option}

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

func (i *Insn) check() error {
	if !i.Option.VLEN.Valid() {
		return fmt.Errorf("wrong VLEN: %d", i.Option.VLEN)
	}

	if !i.Option.ELEN.Valid(i.Option.VLEN) {
		return fmt.Errorf("wrong ELEN: %d", i.Option.ELEN)
	}
	return nil
}

func (i *Insn) Generate() []string {
	res := make([]string, 0)
	for _, code := range i.genCode() {
		builder := strings.Builder{}
		builder.WriteString(i.genHeader())
		builder.WriteString(code)
		builder.WriteString(i.genData())
		res = append(res, builder.String())
	}
	return res
}

func (i *Insn) genHeader() string {
	return fmt.Sprintf(`
# This file is automatically generated. Do not edit.
# Instruction: %s

#include "riscv_test.h"
#include "test_macros.h"

RVTEST_RV64UV
`, i.Name)
}

func (i *Insn) genCode() []string {
	res := make([]string, 0)

	builder := strings.Builder{}
	cs := i.genTestCases()
	for idx, c := range cs {
		builder.WriteString(c)
		if strings.Count(builder.String(), "\n") > 10000 || idx == len(cs)-1 {
			buf := fmt.Sprintf(`
RVTEST_CODE_BEGIN

%s

  TEST_CASE(2, x0, 0x0)
  TEST_PASSFAIL

RVTEST_CODE_END
`, builder.String())
			res = append(res, buf)
			builder.Reset()
		}
	}

	return res
}

func (i *Insn) genData() string {
	dataSize := i.vlenb() * (8 /* max LMUL */)
	// Stride insns
	if strings.HasPrefix(i.Name, "vlse") ||
		strings.HasPrefix(i.Name, "vsse") {
		dataSize *= strides
	} else if strings.HasPrefix(i.Name, "vw") ||
		strings.HasPrefix(i.Name, "vfw") {
		dataSize *= 2
	}

	return fmt.Sprintf(`
  .data
RVTEST_DATA_BEGIN

# Reserve space for test data.
testdata:
  .zero %d

RVTEST_DATA_END
`, dataSize)
}

func (i *Insn) vlenb() int {
	return int(i.Option.VLEN) / 8
}

func (i *Insn) integerTestCases(sew SEW) [][]any {
	return i.testCases(false, sew)
}

func (i *Insn) testCases(float bool, sew SEW) [][]any {
	res := make([][]any, 0)
	if !float {
		for _, c := range i.Tests.Base {
			l := make([]any, len(c))
			for b, op := range c {
				l[b] = op
			}
			res = append(res, l)
		}
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
		if float {
			for _, c := range i.Tests.FSEW32 {
				l := make([]any, len(c))
				for b, op := range c {
					l[b] = math.Float32bits(op)
				}
				res = append(res, l)
			}
			break
		}
		for _, c := range i.Tests.SEW32 {
			l := make([]any, len(c))
			for b, op := range c {
				l[b] = op
			}
			res = append(res, l)
		}
	case 64:
		if float {
			for _, c := range i.Tests.FSEW64 {
				l := make([]any, len(c))
				for b, op := range c {
					l[b] = math.Float64bits(op)
				}
				res = append(res, l)
			}
			break
		}
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
		"\n\n# Generating tests for VL: %d, LMUL: %s, SEW: %s, Mask: %v\n\n",
		c.Vl,
		c.LMUL.String(),
		c.SEW.String(),
		c.Mask)
}

func (i *Insn) combinations(lmuls []LMUL, sews []SEW, masks []bool) []combination {
	res := make([]combination, 0)
	for _, lmul := range lmuls {
		for _, sew := range sews {
			if float64(lmul) < float64(sew)/float64(i.Option.ELEN) {
				continue
			}
			lmul1 := LMUL(math.Max(float64(lmul), 1))
			for _, mask := range masks {
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
