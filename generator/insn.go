package generator

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
)

type insnFormat string

const (
	insnFormatVdRs1mVm    insnFormat = "vd,(rs1),vm"
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
}

type insn struct {
	Name                string     `toml:"name"`
	Format              insnFormat `toml:"format"`
	OperandExchangeable bool       `toml:"operand_exchangeable"`
	Tests               tests      `toml:"tests"`
	Option              Option     `toml:"-"`
}

func ReadInsnFromToml(contents []byte, option Option) (*insn, error) {
	i := insn{Option: option}

	if err := i.Check(); err != nil {
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

func (i *insn) Check() error {
	if !i.Option.VLEN.Valid() {
		return fmt.Errorf("wrong VLEN: %d", i.Option.VLEN)
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

  TEST_PASSFAIL
RVTEST_CODE_END
`, i.genTestCases()))...)
	return buf
}

func (i *insn) genData(buf []byte) []byte {
	buf = append(buf, []byte(fmt.Sprintf(`
.data
RVTEST_DATA_BEGIN

# Generates a 2KiB space for test data.
testdata:
	.zero 2048

RVTEST_DATA_END
`))...)
	return buf
}

func (i *insn) genTestCases() string {
	switch i.Format {
	case insnFormatVdVs2Vs1Vm:
		return i.genCodeVdVs2Vs1Vm()
	default:
		log.Fatalln("unreachable")
		return ""
	}
}
