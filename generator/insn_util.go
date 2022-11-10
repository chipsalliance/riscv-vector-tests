package generator

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type SEW int

var allSEWs = []SEW{8, 16, 32, 64}
var floatSEWs = []SEW{32, 64}
var validSEWs = map[SEW]struct{}{
	allSEWs[0]: {},
	allSEWs[1]: {},
	allSEWs[2]: {},
	allSEWs[3]: {},
}

func (s SEW) String() string {
	if _, ok := validSEWs[s]; !ok {
		log.Fatalln("unreachable")
	}

	return fmt.Sprintf("e%d", s)
}

type LMUL float32

var allLMULs = []LMUL{LMUL(1) / 8, LMUL(1) / 4, LMUL(1) / 2, 1, 2, 4, 8}
var wideningMULs = []LMUL{LMUL(1) / 8, LMUL(1) / 4, LMUL(1) / 2, 1, 2, 4}
var validLMULs = map[LMUL]struct{}{
	allLMULs[0]: {},
	allLMULs[1]: {},
	allLMULs[2]: {},
	allLMULs[3]: {},
	allLMULs[4]: {},
	allLMULs[5]: {},
	allLMULs[6]: {},
}

func (l LMUL) String() string {
	if _, ok := validLMULs[l]; !ok {
		log.Fatalln("unreachable")
	}

	if l < 1 {
		return fmt.Sprintf("mf%d", int(1/l))
	}
	return fmt.Sprintf("m%d", int(l))
}

type VLEN int

func (v VLEN) Valid() bool {
	return 128 <= v && v <= 4096 && v&(v-1) == 0
}

type ELEN int

func (e ELEN) Valid(v VLEN) bool {
	return e >= 64 && e <= ELEN(v) && e&(e-1) == 0
}

func v0t(mask bool) string {
	if mask {
		return ", v0.t"
	}
	return ""
}

func trimBoth(name, prefix, suffix string) string {
	return strings.TrimSuffix(strings.TrimPrefix(name, prefix), suffix)
}

func getEEW(name string) SEW {
	s := regexp.MustCompile(`v.+?(\d+)\.v`)
	eew, err := strconv.Atoi(s.FindStringSubmatch(name)[1])
	if err != nil {
		log.Fatalln("unreachable")
	}
	return SEW(eew)
}

func iff[T any](condition bool, t T, f T) T {
	if condition {
		return t
	}
	return f
}
