package generator

import (
	"reflect"
	"testing"
)

func TestGenRandomData(t *testing.T) {
	for i := int64(1); i <= 10; i++ {
		if !reflect.DeepEqual(genRandomData(i*3), genRandomData(i*3)) {
			t.Fatal()
		}
	}
}
