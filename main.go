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

func parse_extension(march string) []string {
	var valid_exts = []string{"zvbb", "zvbc", "zfh", "zvfh", "zvkg", "zvkned", "zvknha", "zvksed", "zvksh"}
	var exts = []string{}
	exts = append(exts, "v") // standard RVV
	for _, s := range valid_exts {
		for _, e := range strings.Split(march, "_") {
			if e == s {
				exts = append(exts, e)
				break
			}
		}
	}
	return exts
}

type FileTuple struct {
	Entry    os.DirEntry
	FullPath string
}

func walkDirectories(root string) (map[string][]FileTuple, error) {
	filesByDir := make(map[string][]FileTuple)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			parentDir := filepath.Base(filepath.Dir(path))
			filesByDir[parentDir] = append(filesByDir[parentDir], FileTuple{Entry: d, FullPath: path})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return filesByDir, nil
}

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
var float16F = false
var patternF = flag.String("pattern", ".*", "regex to filter out tests.")
var stage1OutputDirF = flag.String("stage1output", "", "stage1 output directory.")
var configsDirF = flag.String("configs", "configs/", "config files directory.")
var testfloat3LevelF = flag.Int("testfloat3level", 2, "testfloat3 testing level (1 or 2).")
var repeatF = flag.Int("repeat", 1, "repeat same V instruction n times for a better coverage (only valid for float instructions).")
var march = flag.String("march", "gcv_zvbb_zvbc_zfh_zvfh_zvkg_zvkned_zvknha_zvksed_zvksh", "march")

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

	extensions := parse_extension(*march)
	extFiles, err := walkDirectories(*configsDirF)
	fatalIf(err)

	fileTuples := make([]FileTuple, 0)
	for _, e := range extensions {
		fs, ok := extFiles[e]
		if ok {
			fileTuples = append(fileTuples, fs...)
			if e == "zvfh" {
				float16F = true
			}
			println("Test extension: ", e)
		}
	}

	println("Generating...")

	makefrag := make([]string, 0)
	lk := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, fileTuple := range fileTuples {
		fp := fileTuple.FullPath
		file := fileTuple.Entry
		if *integerF && (strings.HasPrefix(file.Name(), "vf") || strings.HasPrefix(file.Name(), "vmf")) && !strings.HasPrefix(file.Name(), "vfirst") {
			continue
		}

		if !pattern.MatchString(strings.TrimSuffix(file.Name(), ".toml")) {
			continue
		}

		wg.Add(1)
		go func(file os.DirEntry) {
			name := file.Name()
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
				VLEN:    generator.VLEN(*vlenF),
				XLEN:    generator.XLEN(*xlenF),
				Repeat:  *repeatF,
				Float16: float16F,
			}
			if (!strings.HasPrefix(file.Name(), "vf") && !strings.HasPrefix(file.Name(), "vmf")) || strings.HasPrefix(file.Name(), "vfirst") {
				option.Repeat = 1
			} else {
				option.Fp = true
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
