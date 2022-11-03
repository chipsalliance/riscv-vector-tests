package generator

import (
	"errors"
	"github.com/pelletier/go-toml/v2"
	"os"
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

type insn struct {
	Name                string     `toml:"name"`
	Format              insnFormat `toml:"format"`
	OperandExchangeable bool       `toml:"operand_exchangeable"`
	Tests               tests      `toml:"tests"`
}

func ReadInsnFromFile(filepath string) (*insn, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	i := insn{}
	if err = toml.Unmarshal(contents, &i); err != nil {
		return nil, err
	}

	if err = i.Tests.initialize(); err != nil {
		return nil, err
	}

	if _, ok := formats[i.Format]; !ok {
		return nil, errors.New("invalid test format")
	}

	return &i, nil
}
