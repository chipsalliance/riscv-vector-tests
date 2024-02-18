package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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
var outputFileF = flag.String("outputfile", "", "output file name.")
var configFileF = flag.String("configfile", "", "config file path.")
var testfloat3LevelF = flag.Int("testfloat3level", 2, "testfloat3 testing level (1 or 2).")

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

	testfloat3.SetLevel(*testfloat3LevelF)

	option := generator.Option{
		VLEN: generator.VLEN(*vlenF),
		XLEN: generator.XLEN(*xlenF),
	}

	fp := *configFileF
	contents, err := os.ReadFile(fp)
	fatalIf(err)

	insn, err := generator.ReadInsnFromToml(contents, option)
	fatalIf(err)

	r := regexp.MustCompile(".word 0x.+")
	writeTo(*outputFileF, r.ReplaceAllString(insn.Generate(-1)[0], ""))
}

func writeTo(path string, contents string) {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	fatalIf(err)
	err = os.WriteFile(path, []byte(contents), 0644)
	fatalIf(err)
}
