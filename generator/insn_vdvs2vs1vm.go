package generator

import "fmt"

func (i *insn) genCodeVdVs2Vs1Vm() string {
	res := ""
	for _, lmul := range allLMULs {
		for _, sew := range allSEWs {
			res += fmt.Sprintf(
				"# Generating tests for LMUL: %s, SEW: %s\n\n",
				lmul.String(),
				sew.String())

			// TODO
		}
	}
	return res
}
