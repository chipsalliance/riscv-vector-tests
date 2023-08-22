# machine, user, or sequencer-vector
MODE = machine
VLEN = 256
XLEN = 64
INTEGER = 0
SPIKE_INSTALL = $(RISCV)
OUTPUT = out/v$(VLEN)x$(XLEN)$(MODE)
OUTPUT_STAGE1 = $(OUTPUT)/tests/stage1/
OUTPUT_STAGE2 = $(OUTPUT)/tests/stage2/
OUTPUT_STAGE2_PATCH = $(OUTPUT)/patches/stage2/
OUTPUT_STAGE1_ASM = $(OUTPUT)/asm/stage1/
OUTPUT_STAGE2_ASM = $(OUTPUT)/asm/stage2/
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
RISCV_GCC_OPTS = -static -mcmodel=medany -fvisibility=hidden -nostdlib -nostartfiles

all: compile-stage2

git-submodule-init:
	git submodule update --init

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

generate-stage1: clean-out git-submodule-init build
	@mkdir -p ${OUTPUT_STAGE1}
	build/generator -VLEN ${VLEN} -XLEN ${XLEN} -integer=${INTEGER} -stage1output ${OUTPUT_STAGE1} -configs ${CONFIGS}

-include build/Makefrag

compile-stage1: generate-stage1
	@mkdir -p ${OUTPUT_STAGE1_BIN} ${OUTPUT_STAGE1_ASM}
	$(MAKE) $(tests)

$(tests): %: ${OUTPUT_STAGE1}%.S
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) -Ienv/riscv-test-env/p -Imacros/general -Tenv/riscv-test-env/p/link.ld $< -o ${OUTPUT_STAGE1_BIN}$@
ifeq ($(MODE),sequencer-vector)
	${SPIKE} --isa ${MARCH} --varch=vlen:${VLEN},elen:${XLEN} ${OUTPUT_STAGE1_BIN}$(shell basename $@)
	$(RISCV_GCC) -Ienv/sequencer-vector -Imacros/sequencer-vector -E ${OUTPUT_STAGE1}$(shell basename $@).S -o ${OUTPUT_STAGE1_ASM}$(shell basename $@).S
endif

tests_patch = $(addsuffix .patch, $(tests))

patching-stage2: build-patcher-spike compile-stage1
	@mkdir -p ${OUTPUT_STAGE2_PATCH}
	$(MAKE) $(tests_patch)

$(tests_patch):
	LD_LIBRARY_PATH=$(SPIKE_INSTALL)/lib ${PATCHER_SPIKE} --isa ${MARCH} --varch=vlen:${VLEN},elen:${XLEN} ${OUTPUT_STAGE1_BIN}$(shell basename $@ .patch) > ${OUTPUT_STAGE2_PATCH}$@

generate-stage2: patching-stage2
	build/merger -stage1output ${OUTPUT_STAGE1} -stage2output ${OUTPUT_STAGE2} -stage2patch ${OUTPUT_STAGE2_PATCH}

compile-stage2: generate-stage2
	@mkdir -p ${OUTPUT_STAGE2_ASM}
	@mkdir -p ${OUTPUT_STAGE2_BIN}
	$(MAKE) $(tests_stage2)

tests_stage2 = $(addsuffix .stage2, $(tests))

$(tests_stage2):
ifeq ($(MODE),user)
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) -Ienv/ps -Imacros/general -Tenv/ps/link.ld ${OUTPUT_STAGE2}$(shell basename $@ .stage2).S -o ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
	${SPIKE} --isa ${MARCH} --varch=vlen:${VLEN},elen:${XLEN} $(shell which pk) ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
else ifeq ($(MODE),sequencer/vector)
	$(RISCV_GCC) -Ienv/sequencer-vector -Imacros/sequencer-vector -E ${OUTPUT_STAGE2}$(shell basename $@ .stage2).S -o ${OUTPUT_STAGE2_ASM}$(shell basename $@ .stage2).S
else # machine
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) -Ienv/riscv-test-env/p -Imacros/general -Tenv/riscv-test-env/p/link.ld ${OUTPUT_STAGE2}$(shell basename $@ .stage2).S -o ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
	${SPIKE} --isa ${MARCH} --varch=vlen:${VLEN},elen:${XLEN} ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
endif


clean-out:
	rm -rf $(OUTPUT)

clean: clean-out
	go clean
	rm -rf build/

.PHONY: all \
		git-submodule-init \
		build build-generator unittest \
		generate-stage1 compile-stage1 $(tests) \
		$(tests_patch) generate-stage2 compile-stage2 $(tests_stage2) \
		clean-out clean
