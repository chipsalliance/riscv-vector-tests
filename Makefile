VLEN = 256
ELEN = 64
OUTPUT = out/
OUTPUT_STAGE1 = ${OUTPUT}tests/stage1/
OUTPUT_STAGE1_BIN = ${OUTPUT}bin/stage1/
CONFIGS = configs/

include Makefrag

RISCV_PREFIX = riscv64-unknown-elf-
RISCV_GCC = $(RISCV_PREFIX)gcc
RISCV_GCC_OPTS = -static -mcmodel=medany -fvisibility=hidden -nostdlib -nostartfiles

build:
	go build -o build/generator

unittest:
	go test ./...

generate-stage1: build
	@mkdir -p ${OUTPUT_STAGE1}
	build/generator -VLEN ${VLEN} -ELEN ${ELEN} -stage1output ${OUTPUT_STAGE1} -configs ${CONFIGS}

compile-stage1: generate-stage1
	@mkdir -p ${OUTPUT_STAGE1_BIN}
	$(MAKE) $(tests)

$(tests): %: ${OUTPUT_STAGE1}%.S
	$(RISCV_GCC) -march=rv64gv $(RISCV_GCC_OPTS) -Ienv/p -Imacros -Tenv/p/link.ld $< -o ${OUTPUT_STAGE1_BIN}$@

clean:
	go clean
	rm -rf out/
	rm -rf build/

.PHONY: build unittest generate-stage1 compile-stage1 $(tests) clean