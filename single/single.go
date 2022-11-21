package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ksco/riscv-vector-tests/generator"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func fatalIf(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\033[0;1;31mfatal:\033[0m %s\n", err.Error())
	os.Exit(1)
}

var vlenF = flag.Int("VLEN", 256, "")
var xlenF = flag.Int("XLEN", 64, "")
var outputDirF = flag.String("output", "", "output directory.")
var configFileF = flag.String("configfile", "", "config file path.")

func main() {
	flag.Parse()

	if outputDirF == nil || *outputDirF == "" {
		fatalIf(errors.New("-output is required"))
	}
	if configFileF == nil || *configFileF == "" {
		fatalIf(errors.New("-configfile is required"))
	}

	option := generator.Option{
		VLEN: generator.VLEN(*vlenF),
		XLEN: generator.XLEN(*xlenF),
	}

	println("Generating...")

	fp := *configFileF
	contents, err := os.ReadFile(fp)
	fatalIf(err)

	insn, err := generator.ReadInsnFromToml(contents, option)
	fatalIf(err)

	for idx, testContent := range insn.Generate(false) {
		asmFilename := strings.TrimSuffix(filepath.Base(fp), ".toml") + "-" + strconv.Itoa(idx)
		writeTo(
			*outputDirF,
			asmFilename+".S",
			testContent)
	}

	println("\033[32mOK\033[0m")
}

func writeTo(path string, name string, contents string) {
	err := os.MkdirAll(path, 0777)
	fatalIf(err)
	err = os.WriteFile(filepath.Join(path, name), []byte(contents), 0644)
	fatalIf(err)
}
