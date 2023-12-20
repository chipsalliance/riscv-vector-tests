MODE = machine
VLEN = 256
XLEN = 64
SPLIT = 10000
INTEGER = 0
PATTERN = '.*'
SPIKE_INSTALL = $(RISCV)
OUTPUT = out/v$(VLEN)x$(XLEN)$(MODE)
OUTPUT_STAGE1 = $(OUTPUT)/tests/stage1/
OUTPUT_STAGE2 = $(OUTPUT)/tests/stage2/
OUTPUT_STAGE2_PATCH = $(OUTPUT)/patches/stage2/
OUTPUT_STAGE1_BIN = $(OUTPUT)/bin/stage1/
OUTPUT_STAGE2_BIN = $(OUTPUT)/bin/stage2/
CONFIGS = configs/

SPIKE = spike
PATCHER_SPIKE = build/pspike
MARCH = rv${XLEN}gcv
MABI = lp64d

ifeq ($(XLEN), 32)
MABI = ilp32f
endif

RISCV_PREFIX = riscv64-unknown-elf-
RISCV_GCC = $(RISCV_PREFIX)gcc
RISCV_GCC_OPTS = -static -mcmodel=medany -fvisibility=hidden -nostdlib -nostartfiles -DENTROPY=0xdeadbeef -DLFSR_BITS=9 -fno-tree-loop-distribute-patterns

PK =

ifeq ($(MODE),machine)
ENV = env/riscv-test-env/p
endif
ifeq ($(MODE),user)
ENV = env/ps
PK = $(shell which pk)
ifeq (, $(PK))
$(error "No pk found, please install it to your path.")
endif
endif
ifeq ($(MODE),virtual)
ENV = env/riscv-test-env/v
ENV_CSRCS = env/riscv-test-env/v/vm.c env/riscv-test-env/v/string.c env/riscv-test-env/v/entry.S
endif

all: compile-stage2

build: build-generator build-merger

build-generator:
	go build -o build/generator

build-merger:
	go build -o build/merger merger/merger.go

build-patcher-spike: pspike/pspike.cc
	rm -rf build
	mkdir -p build
	g++ -std=c++17 -I$(SPIKE_INSTALL)/include -L$(SPIKE_INSTALL)/lib $< -lriscv -lfesvr -o $(PATCHER_SPIKE)

unittest:
	go test ./...

generate-stage1: clean-out build
	@mkdir -p ${OUTPUT_STAGE1}
	build/generator -VLEN ${VLEN} -XLEN ${XLEN} -split=${SPLIT} -integer=${INTEGER} -pattern='${PATTERN}' -stage1output ${OUTPUT_STAGE1} -configs ${CONFIGS}

include Makefrag

compile-stage1: generate-stage1
	@mkdir -p ${OUTPUT_STAGE1_BIN}
	$(MAKE) $(tests)

$(tests): %: ${OUTPUT_STAGE1}%.S
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) -I$(ENV) -Imacros/general -T$(ENV)/link.ld $(ENV_CSRCS) $< -o ${OUTPUT_STAGE1_BIN}$@

tests_patch = $(addsuffix .patch, $(tests))

patching-stage2: build-patcher-spike compile-stage1
	@mkdir -p ${OUTPUT_STAGE2_PATCH}
	$(MAKE) $(tests_patch)

$(tests_patch):
	LD_LIBRARY_PATH=$(SPIKE_INSTALL)/lib ${PATCHER_SPIKE} --isa=${MARCH} --varch=vlen:${VLEN},elen:${XLEN} $(PK) ${OUTPUT_STAGE1_BIN}$(shell basename $@ .patch) > ${OUTPUT_STAGE2_PATCH}$@

generate-stage2: patching-stage2
	build/merger -stage1output ${OUTPUT_STAGE1} -stage2output ${OUTPUT_STAGE2} -stage2patch ${OUTPUT_STAGE2_PATCH}

compile-stage2: generate-stage2
	@mkdir -p ${OUTPUT_STAGE2_BIN}
	$(MAKE) $(tests_stage2)

tests_stage2 = $(addsuffix .stage2, $(tests))

$(tests_stage2):
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) -I$(ENV) -Imacros/general -T$(ENV)/link.ld $(ENV_CSRCS) ${OUTPUT_STAGE2}$(shell basename $@ .stage2).S -o ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
	${SPIKE} --isa=${MARCH} --varch=vlen:${VLEN},elen:${XLEN} $(PK) ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)


clean-out:
	rm -rf $(OUTPUT)

clean: clean-out
	go clean
	rm -rf build/

.PHONY: all \
		build build-generator unittest \
		generate-stage1 compile-stage1 $(tests) \
		$(tests_patch) generate-stage2 compile-stage2 $(tests_stage2) \
		clean-out clean
