package generator

import (
	"encoding/binary"
	"fmt"
)

func (i *insn) genCodeVdVs2Vs1Vm() string {
	res := ""

	res += "la a0, testdata\n\n"

	res += "# Write random data into test data area.\n"
	rdata := genRandomData(int64(i.vlenb()))
	for a := 0; a < len(rdata)/8; a++ {
		elem := binary.LittleEndian.Uint64(rdata)
		rdata = rdata[8:]

		res += fmt.Sprintf("li a1, 0x%x\n", elem)
		res += fmt.Sprintf("sd a1, %d(a0)\n", a*8)
	}

	res += "\n# Load random data into v0 mask register.\n"
	res += "li t0, -1\n"
	res += "vsetvli t1, t0, e8,m1,ta,ma\n"
	res += "vle8.v v0, (a0)"

	for _, c := range i.combinations() {
		res += fmt.Sprintf(
			"\n\n# Generating tests for LMUL: %s, SEW: %s, Mask: %v\n\n",
			c.LMUL.String(),
			c.SEW.String(),
			c.Mask)

		nBytes := i.vlenb() * int(c.LMUL1)
		rdata = genRandomData(int64(nBytes))

		res += "# Write random data into test data area.\n"
		for a := 0; a < nBytes/8; a++ {
			elem := binary.LittleEndian.Uint64(rdata)
			rdata = rdata[8:]

			res += fmt.Sprintf("li a1, 0x%x\n", elem)
			res += fmt.Sprintf("sd a1, %d(a0)\n", a*8)
		}

		vd := fmt.Sprintf("v%d", 1*int(c.LMUL1))
		vss := []string{
			fmt.Sprintf("v%d", 2*int(c.LMUL1)),
			fmt.Sprintf("v%d", 3*int(c.LMUL1)),
		}

		res += "\n# Load random data into vd register group.\n"
		res += "li t0, -1\n"
		res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
			c.SEW.String(), c.LMUL1.String())
		res += fmt.Sprintf("vle%d.v %s, (a0)\n\n", c.SEW, vd)

		for _, op := range []int{1, 2} {
			idx := op - 1
			res += fmt.Sprintf("# Write op%d data into test data area.\n", op)
			cases := i.testCases(c.SEW)
			for a := 0; a < (nBytes / (int(c.SEW) / 8)); a++ {
				b := a % len(cases)
				switch c.SEW {
				case 8:
					res += fmt.Sprintf("li, a1, 0x%x\n", convNum[uint8](cases[b][idx]))
					res += fmt.Sprintf("sb, a1, %d(a0)\n", a*(int(c.SEW)/8))
				case 16:
					res += fmt.Sprintf("li, a1, 0x%x\n", convNum[uint16](cases[b][idx]))
					res += fmt.Sprintf("sh, a1, %d(a0)\n", a*(int(c.SEW)/8))
				case 32:
					res += fmt.Sprintf("li, a1, 0x%x\n", convNum[uint32](cases[b][idx]))
					res += fmt.Sprintf("sw, a1, %d(a0)\n", a*(int(c.SEW)/8))
				case 64:
					res += fmt.Sprintf("li a1, 0x%x\n", convNum[uint64](cases[b][idx]))
					res += fmt.Sprintf("sd, a1, %d(a0)\n", a*(int(c.SEW)/8))
				}
			}

			res += fmt.Sprintf("\n# Load test data into vs%d register group.\n", op)
			res += "li t0, -1\n"
			res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
				c.SEW.String(), c.LMUL1.String())
			res += fmt.Sprintf("vle%d.v %s, (a0)\n\n", c.SEW, vss[idx])
		}

		res += "# -------------- TEST BEGIN --------------\n"
		res += fmt.Sprintf("li t0, %d\n", c.Vl)
		res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
			c.SEW.String(), c.LMUL.String())
		res += fmt.Sprintf("%s %s, %s, %s%s\n", i.Name, vd, vss[1], vss[0], v0t(c.Mask))
		res += "# -------------- TEST END   --------------\n"

		res += "\n# Store vd register group into test data area.\n"
		res += "li t0, -1\n"
		res += fmt.Sprintf("vsetvli t1, t0, %s,%s,ta,ma\n",
			c.SEW.String(), c.LMUL1.String())
		res += fmt.Sprintf("vse%d.v %s, (a0)\n\n", c.SEW, vd)

		// Magic insn
		res += fmt.Sprintf("addi, x0, x%d, %d", 1*int(c.LMUL1), 2*int(c.LMUL1))
	}
	return res
}
