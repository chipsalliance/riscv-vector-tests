package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ksco/riscv-vector-tests/generator"
	"os"
	"path/filepath"
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

func main() {
	flag.Parse()

	if outputFileF == nil || *outputFileF == "" {
		fatalIf(errors.New("-outputfile is required"))
	}
	if configFileF == nil || *configFileF == "" {
		fatalIf(errors.New("-configfile is required"))
	}

	option := generator.Option{
		VLEN: generator.VLEN(*vlenF),
		XLEN: generator.XLEN(*xlenF),
	}

	fp := *configFileF
	contents, err := os.ReadFile(fp)
	fatalIf(err)

	insn, err := generator.ReadInsnFromToml(contents, option)
	fatalIf(err)

	writeTo(*outputFileF, insn.Generate(-1)[0])
}

func writeTo(path string, contents string) {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	fatalIf(err)
	err = os.WriteFile(path, []byte(contents), 0644)
	fatalIf(err)
}
