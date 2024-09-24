##-----------------------------
##RISC-V Vector Tests Generator
##-----------------------------
##
##Usage: make all -j$(nproc) --environment-overrides [OPTIONS]
##
##Example: to generate isa=rv32gcv_zvl128b_zve32f mode=machine tests, use:
## make all --environment-overrides VLEN=128 XLEN=32 MODE=machine -j$(nproc)
##
##Subcommands:
  help: ## Show this help message.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

  all:  ## Generate tests.
  all: compile-stage2

##
##Options:

  MODE = machine##
        ##Can be [machine], [virtual] or [user]; for more info, see https://github.com/riscv/riscv-test-env
        ##
  VLEN = 256##
        ##Can vary from [64] to [4096] (upper boundary is limited by Spike)
        ##
  XLEN = 64##
        ##Can be [32] or [64]; note that we do not support specifying ELEN yet, ELEN is always consistent with XLEN
        ##
  SPLIT = 10000##
        ##Split on assembly file lines (excluding test data)
        ##
  INTEGER = 0##
        ##Set to [1] if you don't want float tests (i.e. for Zve32x or Zve64x)
        ##
  FLOAT16 = 1##
        ##Set to [0] if you don't want float16 (Zvfh) tests
        ##
  PATTERN = '.*'##
        ##Set to a valid regex to generate the tests of your interests (e.g. PATTERN='^v[ls].+\.v$' to generate only load/store tests)
        ##
  TESTFLOAT3LEVEL = 2##
        ##Testing level for testfloat3 generated cases, can be one of [1] or [2].
        ##
  REPEAT = 1##
        ##Set to greater value to repeat the same V instruction n times for a better coverage (only valid for float instructions).
        ##
  TEST_MODE = self##
        ##Change to [cosim] if you want to generate faster tests without self-verification (to be used with co-simulators).
        ##
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
MARCH = rv${XLEN}gcv_zvbb_zvbc_zfh_zvfh_zvkg_zvkned_zvknha_zvksed_zvksh
MABI = lp64d

ifeq ($(XLEN), 32)
MABI = ilp32f
VARCH = zvl${VLEN}b_zve32f_zfh_zfhmin_zvfh
else
VARCH = zvl${VLEN}b_zve64d_zfh_zfhmin_zvfh
endif

RISCV_PREFIX = riscv64-unknown-elf-
RISCV_GCC = $(RISCV_PREFIX)gcc
RISCV_GCC_OPTS = -static -mcmodel=medany -fvisibility=hidden -nostdlib -nostartfiles -DENTROPY=0xdeadbeef -DLFSR_BITS=9 -fno-tree-loop-distribute-patterns

ifeq ($(TEST_MODE), self)
STAGE2_GCC_OPTS =
else ifeq ($(TEST_MODE), cosim)
STAGE2_GCC_OPTS = -DCOSIM_TEST_CASE
else
$(error "Only self and cosim are supported for TEST_MODE")
endif

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
	build/generator -VLEN ${VLEN} -XLEN ${XLEN} -split=${SPLIT} -integer=${INTEGER} -float16=${FLOAT16} -pattern='${PATTERN}' -testfloat3level='${TESTFLOAT3LEVEL}' -repeat='${REPEAT}' -stage1output ${OUTPUT_STAGE1} -configs ${CONFIGS}

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
	LD_LIBRARY_PATH=$(SPIKE_INSTALL)/lib ${PATCHER_SPIKE} --isa=${MARCH}_${VARCH} $(PK) ${OUTPUT_STAGE1_BIN}$(shell basename $@ .patch) > ${OUTPUT_STAGE2_PATCH}$@

generate-stage2: patching-stage2
	build/merger -stage1output ${OUTPUT_STAGE1} -stage2output ${OUTPUT_STAGE2} -stage2patch ${OUTPUT_STAGE2_PATCH}

compile-stage2: generate-stage2
	@mkdir -p ${OUTPUT_STAGE2_BIN}
	$(MAKE) $(tests_stage2)

tests_stage2 = $(addsuffix .stage2, $(tests))

$(tests_stage2):
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) $(STAGE2_GCC_OPTS) -I$(ENV) -Imacros/general -T$(ENV)/link.ld $(ENV_CSRCS) ${OUTPUT_STAGE2}$(shell basename $@ .stage2).S -o ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
	${SPIKE} --isa=${MARCH}_${VARCH} $(PK) ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)


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
