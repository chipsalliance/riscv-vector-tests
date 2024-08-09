package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ksco/riscv-vector-tests/generator"
	"github.com/ksco/riscv-vector-tests/testfloat3"
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
var testfloat3LevelF = flag.Int("testfloat3level", 2, "testfloat3 testing level (1 or 2).")
var repeatF = flag.Int("repeat", 1, "repeat same V instruction n times for a better coverage (only valid for float instructions).")

func main() {
	flag.Parse()

	pattern, err := regexp.Compile(*patternF)
	fatalIf(err)

	if stage1OutputDirF == nil || *stage1OutputDirF == "" {
		fatalIf(errors.New("-stage1output is required"))
	}

	if !(*testfloat3LevelF == 1 || *testfloat3LevelF == 2) {
		fatalIf(errors.New("-testfloat3level must be 1 or 2"))
	}

	if *repeatF <= 0 {
		fatalIf(errors.New("-repeat must greater than 0"))
	}

	testfloat3.SetLevel(*testfloat3LevelF)

	files, err := os.ReadDir(*configsDirF)
	fatalIf(err)

	println("Generating...")

	makefrag := make([]string, 0)
	lk := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, file := range files {
		if *integerF && (strings.HasPrefix(file.Name(), "vf") || strings.HasPrefix(file.Name(), "vmf")) && !strings.HasPrefix(file.Name(), "vfirst") {
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

			name = strings.TrimSuffix(name, ".toml")
			name = strings.Replace(name, ".", "_", -1)

			contents, err := os.ReadFile(fp)
			fatalIf(err)

			option := generator.Option{
				VLEN:   generator.VLEN(*vlenF),
				XLEN:   generator.XLEN(*xlenF),
				Repeat: *repeatF,
			}
			if (!strings.HasPrefix(file.Name(), "vf") && !strings.HasPrefix(file.Name(), "vmf")) || strings.HasPrefix(file.Name(), "vfirst") {
				option.Repeat = 1
			}
			insn, err := generator.ReadInsnFromToml(contents, option)
			fatalIf(err)
			if insn.Name != strings.Replace(file.Name(), ".toml", "", -1) {
				fatalIf(errors.New("filename and instruction name unmatched"))
			}

			for idx, testContent := range insn.Generate(*splitF) {
				asmFilename := name + "-" + strconv.Itoa(idx)
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
