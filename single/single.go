package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ksco/riscv-vector-tests/generator"
	"github.com/ksco/riscv-vector-tests/testfloat3"
)

func fatalIf(err error) {
	if err == nil {
		return
	}
	fmt.Printf("fatal: %s\n", err.Error())
	os.Exit(1)
}

var vlenF = flag.Int("VLEN", 256, "")
var xlenF = flag.Int("XLEN", 64, "")
var splitF = flag.Int("split", -1, "split per lines.")
var float16F = flag.Bool("float16", true, "")
var outputFileF = flag.String("outputfile", "", "output file name.")
var configFileF = flag.String("configfile", "", "config file path.")
var testfloat3LevelF = flag.Int("testfloat3level", 2, "testfloat3 testing level (1 or 2).")
var repeatF = flag.Int("repeat", 1, "repeat same V instruction n times for a better coverage (only valid for float instructions).")

func main() {
	flag.Parse()

	if outputFileF == nil || *outputFileF == "" {
		fatalIf(errors.New("-outputfile is required"))
	}
	if configFileF == nil || *configFileF == "" {
		fatalIf(errors.New("-configfile is required"))
	}

	if !(*testfloat3LevelF == 1 || *testfloat3LevelF == 2) {
		fatalIf(errors.New("-testfloat3level must be 1 or 2"))
	}

	if *repeatF <= 0 {
		fatalIf(errors.New("-repeat must greater than 0"))
	}

	testfloat3.SetLevel(*testfloat3LevelF)

	option := generator.Option{
		VLEN:    generator.VLEN(*vlenF),
		XLEN:    generator.XLEN(*xlenF),
		Repeat:  *repeatF,
		Float16: *float16F,
	}

	fp := *configFileF
	contents, err := os.ReadFile(fp)
	fatalIf(err)

	if (!strings.HasPrefix(filepath.Base(fp), "vf") && !strings.HasPrefix(filepath.Base(fp), "vmf")) || strings.HasPrefix(filepath.Base(fp), "vfirst") {
		option.Repeat = 1
	} else {
		option.Fp = true
	}

	insn, err := generator.ReadInsnFromToml(contents, option)
	fatalIf(err)

	r := regexp.MustCompile(".word 0x.+")
	writeTo(*outputFileF, r.ReplaceAllString(insn.Generate(*splitF)[0], ""))
}

func writeTo(path string, contents string) {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	fatalIf(err)
	err = os.WriteFile(path, []byte(contents), 0644)
	fatalIf(err)
}
