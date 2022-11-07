package main

import (
	"flag"
	"fmt"
	"github.com/ksco/riscv-vector-tests/generator"
	"os"
	"path/filepath"
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
var outputDirF = flag.String("output", "out/", "output directory.")
var configsDirF = flag.String("configs", "configs/", "config files directory.")

var stage1OutputDir = filepath.Join(*outputDirF, "tests/stage1")

func main() {
	flag.Parse()

	option := generator.Option{
		VLEN: generator.VLEN(*vlenF),
	}

	err := os.RemoveAll(*outputDirF)
	fatalIf(err)

	files, err := os.ReadDir(*configsDirF)
	fatalIf(err)

	println("Generating...")
	for _, file := range files {
		name := file.Name()
		fp := filepath.Join(*configsDirF, name)
		if file.IsDir() ||
			!strings.HasPrefix(name, "v") ||
			!strings.HasSuffix(name, ".toml") {
			fmt.Printf("\033[0;1;31mskipping:\033[0m %s, unrecognized filename\n", fp)
			continue
		}

		contents, err := os.ReadFile(fp)
		fatalIf(err)

		insn, err := generator.ReadInsnFromToml(contents, option)
		fatalIf(err)

		writeTo(stage1OutputDir,
			strings.TrimSuffix(name, ".toml")+".S", insn.Generate())
	}

	println("\033[32mOK\033[0m")
}

func writeTo(path string, name string, contents []byte) {
	err := os.MkdirAll(path, 0777)
	fatalIf(err)
	err = os.WriteFile(filepath.Join(path, name), contents, 0644)
	fatalIf(err)
}
