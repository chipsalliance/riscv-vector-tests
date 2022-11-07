package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ksco/riscv-vector-tests/generator"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func fatalIf(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\033[0;1;31mfatal:\033[0m %s\n", err.Error())
	os.Exit(1)
}

var vlenF = flag.Int("VLEN", 256, "")
var elenF = flag.Int("ELEN", 64, "")
var stage1OutputDirF = flag.String("stage1output", "", "stage1 output directory.")
var configsDirF = flag.String("configs", "configs/", "config files directory.")
var rewriteMakeFrag = flag.Bool("rewrite-makefrag", true, "rewrite makefrag file.")

func main() {
	flag.Parse()

	if stage1OutputDirF == nil || *stage1OutputDirF == "" {
		fatalIf(errors.New("-stage1output is required"))
	}

	option := generator.Option{
		VLEN: generator.VLEN(*vlenF),
		ELEN: generator.ELEN(*elenF),
	}

	files, err := os.ReadDir(*configsDirF)
	fatalIf(err)

	println("Generating...")

	if rewriteMakeFrag != nil && *rewriteMakeFrag {
		makefrag := "tests = \\\n"
		for _, file := range files {
			filename := strings.TrimSuffix(file.Name(), ".toml")
			makefrag += fmt.Sprintf("  %s \\\n", filename)
		}
		writeTo(".", "Makefrag", []byte(makefrag))
	}

	lk := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(files))
	for _, file := range files {
		go func(file os.DirEntry) {
			name := file.Name()
			fp := filepath.Join(*configsDirF, name)
			if file.IsDir() ||
				!strings.HasPrefix(name, "v") ||
				!strings.HasSuffix(name, ".toml") {
				lk.Lock()
				fmt.Printf("\033[0;1;31mskipping:\033[0m %s, unrecognized filename\n", fp)
				lk.Unlock()
				return
			}

			contents, err := os.ReadFile(fp)
			fatalIf(err)

			insn, err := generator.ReadInsnFromToml(contents, option)
			fatalIf(err)

			writeTo(*stage1OutputDirF,
				strings.TrimSuffix(name, ".toml")+".S", insn.Generate())
			wg.Done()
		}(file)
	}
	wg.Wait()

	println("\033[32mOK\033[0m")
}

func writeTo(path string, name string, contents []byte) {
	err := os.MkdirAll(path, 0777)
	fatalIf(err)
	err = os.WriteFile(filepath.Join(path, name), contents, 0644)
	fatalIf(err)
}
