package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ksco/riscv-vector-tests/generator"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
var xlenF = flag.Int("XLEN", 64, "we do not support specifying ELEN yet, ELEN is consistent with XLEN.")
var splitF = flag.Int("split", 10000, "split per lines.")
var integerF = flag.Bool("integer", false, "only generate integer tests.")
var patternF = flag.String("pattern", ".*", "regex to filter out tests.")
var stage1OutputDirF = flag.String("stage1output", "", "stage1 output directory.")
var configsDirF = flag.String("configs", "configs/", "config files directory.")

func main() {
	flag.Parse()

	pattern, err := regexp.Compile(*patternF)
	fatalIf(err)

	if stage1OutputDirF == nil || *stage1OutputDirF == "" {
		fatalIf(errors.New("-stage1output is required"))
	}

	option := generator.Option{
		VLEN: generator.VLEN(*vlenF),
		XLEN: generator.XLEN(*xlenF),
	}

	files, err := os.ReadDir(*configsDirF)
	fatalIf(err)

	println("Generating...")

	makefrag := make([]string, 0)
	lk := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, file := range files {
		if *integerF && strings.HasPrefix(file.Name(), "vf") {
			continue
		}

		if !pattern.MatchString(strings.TrimSuffix(file.Name(), ".toml")) {
			continue
		}

		wg.Add(1)
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

			for idx, testContent := range insn.Generate(*splitF) {
				asmFilename := strings.TrimSuffix(name, ".toml") + "-" + strconv.Itoa(idx)
				writeTo(*stage1OutputDirF, asmFilename+".S", testContent)
				lk.Lock()
				makefrag = append(makefrag, fmt.Sprintf("  %s \\\n", asmFilename))
				lk.Unlock()
			}

			wg.Done()
		}(file)
	}
	wg.Wait()

	sort.Slice(makefrag, func(i, j int) bool {
		return makefrag[i] < makefrag[j]
	})
	writeTo("./", "Makefrag", "tests = \\\n"+strings.Join(makefrag, ""))

	println("\033[32mOK\033[0m")
}

func writeTo(path string, name string, contents string) {
	err := os.MkdirAll(path, 0777)
	fatalIf(err)
	err = os.WriteFile(filepath.Join(path, name), []byte(contents), 0644)
	fatalIf(err)
}
