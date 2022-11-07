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

	makefrag := "tests = \\\n"
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

			for idx, testContent := range insn.Generate() {
				asmFilename := strings.TrimSuffix(name, ".toml") + "-" + strconv.Itoa(idx)
				writeTo(
					*stage1OutputDirF,
					asmFilename+".S",
					testContent)

				lk.Lock()
				makefrag += fmt.Sprintf("  %s \\\n", asmFilename)
				lk.Unlock()
			}

			wg.Done()
		}(file)
	}
	wg.Wait()

	writeTo(".", "Makefrag", makefrag)

	println("\033[32mOK\033[0m")
}

func writeTo(path string, name string, contents string) {
	err := os.MkdirAll(path, 0777)
	fatalIf(err)
	err = os.WriteFile(filepath.Join(path, name), []byte(contents), 0644)
	fatalIf(err)
}
