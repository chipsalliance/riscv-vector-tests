package generator

import (
	"strconv"
	"strings"
)

type num interface {
	uint8 | uint16 | uint32 | uint64 | string
}

type testCase[T num] []T

type tests struct {
	Base  []testCase[uint8]  `toml:"base"`
	SEW8  []testCase[uint8]  `toml:"sew8"`
	SEW16 []testCase[uint16] `toml:"sew16"`
	SEW32 []testCase[uint32] `toml:"sew32"`

	// Go toml cannot parse uint64, we parse it ourself.
	SEW64_ []testCase[string] `toml:"sew64"`
	SEW64  []testCase[uint64] `toml:"-"`
}

func (t *tests) initialize() error {
	var err error
	for i, ss := range t.SEW64_ {
		t.SEW64 = append(t.SEW64, make([]uint64, len(ss)))
		for j, s := range ss {
			t.SEW64[i][j], err = strconv.ParseUint(
				strings.TrimPrefix(s, "0x"), 16, 64)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
