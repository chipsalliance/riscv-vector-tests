package generator

import (
	"fmt"
	"log"
	"math"
	"strings"
)

func getSew(vtype int64) int {
	sewMask := 0x38
	sew := vtype & int64(sewMask)
	if (sew >> 3) > 3 {
		log.Fatalln("illegal vsew")
	}
	return 8 << (sew >> 3)
}

func getLmul(vtype int64) LMUL {
	lmulmask := 0x7
	t := vtype & int64(lmulmask)
	lmul := LMUL(int(1) << t)
	if t >= 5 {
		lmul = 1 / LMUL(int(1)<<(8-t))
	}
	return LMUL(lmul)
}

func getVlmax(vtype int64, vlen int) int {
	sew := getSew(vtype)
	vlmax := LMUL((vlen / sew)) * getLmul(vtype)
	return int(vlmax)
}

func (v *vtype) vtypeRaw(XLEN int, VLEN int, newVtypeRaw int64, curVtypeRaw int64, curVl int64, rd int64, rs1 int64) []int64 {
	res := make([]int64, 0, 3)
	// vstart
	res = append(res, 0)

	vlmax := int64(getVlmax(curVtypeRaw, VLEN))
	VtypeRaw := newVtypeRaw
	if VtypeRaw != curVtypeRaw {
		vlmax = int64(getVlmax(VtypeRaw, VLEN))
		newLmul := getLmul(VtypeRaw)
		newVsew := getSew(VtypeRaw)
		// We assume ELEN is consistent with XLEN.
		new_vill :=
			!(newLmul >= 0.125 && newLmul <= 8) ||
				float64(newVsew) > math.Min(float64(newLmul), float64(1))*float64(XLEN) ||
				bits64(VtypeRaw, XLEN-2, 8) != 0
		if new_vill {
			vlmax = 0
			VtypeRaw = 1 << (XLEN - 1)
		}
	}
	// vtype
	res = append(res, VtypeRaw)

	vl := curVl
	if vlmax == 0 {
		vl = 0
	} else if rd == 0 && rs1 == 0 {
		if curVl > int64(vlmax) {
			vl = curVl
		} else {
			vl = curVl
		}
	} else if rd != 0 && rs1 == 0 {
		vl = vlmax
	} else if rs1 != 0 {
		if rs1 > vlmax {
			vl = vlmax
		} else {
			vl = rs1
		}
	}
	// vl
	res = append(res, vl)

	return res
}

func (i *Insn) genCodevsetvl() []string {
	res := make([]string, 0, 1)
	ncase := 3
	builder := strings.Builder{}
	cases := i.testCases(false, 8)

	curvtype := int64(0)
	curvl := int64(0)
	for _, cs := range cases {
		builder.WriteString(fmt.Sprintf("# ------case test begin---------\n"))

		builder.WriteString(fmt.Sprintf("li t0, %d\n", cs[0]))
		builder.WriteString(fmt.Sprintf("li t1, %d\n", cs[1]))
		builder.WriteString(fmt.Sprintf("li t2, %d\n", cs[2]))
		builder.WriteString(fmt.Sprintf("vsetvl t0, t1, t2\n"))

		var v vtype
		t := v.vtypeRaw(int(i.Option.XLEN), int(i.Option.VLEN), int64(cs[2].(uint8)), curvtype, curvl, int64(cs[0].(uint8)), int64(cs[1].(uint8)))

		builder.WriteString(fmt.Sprintf("csrr a4, vstart\n"))
		builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[0]))
		ncase = ncase + 1

		builder.WriteString(fmt.Sprintf("csrr a4, vtype\n"))
		if t[1] == 1<<(i.Option.XLEN-1) {
			builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, uint64(1<<(i.Option.XLEN-1))))
		} else {
			builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[1]))
		}
		ncase = ncase + 1

		builder.WriteString(fmt.Sprintf("csrr a4, vl\n"))
		builder.WriteString(fmt.Sprintf("TEST_CASE(%d, a4, %d)\n", ncase, t[2]))
		ncase = ncase + 1

		curvtype = t[1]
		curvl = t[2]

		builder.WriteString(fmt.Sprintf("# ------case test end---------\n"))
	}
	res = append(res, builder.String())

	return res
}
