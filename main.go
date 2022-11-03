package main

import (
	"fmt"
	"github.com/ksco/riscv-vector-tests/generator"
)

func main() {
	i, err := generator.ReadInsnFromFile("tests/vadd.vv.toml")
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Printf("%v\n", i.Name)
}
