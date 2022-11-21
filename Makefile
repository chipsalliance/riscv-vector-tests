MODE = machine # or user, sequencer/vector
VLEN = 256
XLEN = 64
INTEGER = 0
OUTPUT = out/
OUTPUT_STAGE1 = ${OUTPUT}tests/stage1/
OUTPUT_STAGE2 = ${OUTPUT}tests/stage2/
OUTPUT_STAGE2_PATCH = ${OUTPUT}patches/stage2/
OUTPUT_STAGE1_ASM = ${OUTPUT}asm/stage1/
OUTPUT_STAGE2_ASM = ${OUTPUT}asm/stage2/
OUTPUT_STAGE1_BIN = ${OUTPUT}bin/stage1/
OUTPUT_STAGE2_BIN = ${OUTPUT}bin/stage2/
CONFIGS = configs/

SPIKE = spike
PATCHED_SPIKE = build/spike/build/spike
MARCH = rv${XLEN}gcv
MABI = lp64d

ifeq ($(XLEN), 32)
MABI = ilp32f
endif

include Makefrag

RISCV_PREFIX = riscv64-unknown-elf-
RISCV_GCC = $(RISCV_PREFIX)gcc
RISCV_GCC_OPTS = -static -mcmodel=medany -fvisibility=hidden -nostdlib -nostartfiles

all: compile-stage2

build: build-generator build-merger

build-generator:
	go build -o build/generator

build-merger:
	go build -o build/merger merger/merger.go

build-spike:
	git clone --depth 1 https://github.com/riscv-software-src/riscv-isa-sim.git build/spike/ 2>/dev/null || true
	cd build/spike/; \
	git reset --hard origin/master; \
	git am --abort || true; \
	git am < ../../patches/0001-RV${XLEN}-Modify-addi-to-generate-test-cases.patch; \
	mkdir -p build; \
	cd build; \
	../configure --prefix=$(.); \
	$(MAKE) -j8 \

unittest:
	go test ./...

generate-stage1: clean-out build
	@mkdir -p ${OUTPUT_STAGE1}
	build/generator -VLEN ${VLEN} -XLEN ${XLEN} -integer=${INTEGER} -stage1output ${OUTPUT_STAGE1} -configs ${CONFIGS}

compile-stage1: generate-stage1
	@mkdir -p ${OUTPUT_STAGE1_BIN} ${OUTPUT_STAGE1_ASM}
	$(MAKE) $(tests)

$(tests): %: ${OUTPUT_STAGE1}%.S
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) -Ienv/p -Imacros/general -Tenv/p/link.ld $< -o ${OUTPUT_STAGE1_BIN}$@
ifeq ($(MODE),sequencer/vector)
	${SPIKE} --isa ${MARCH} --varch=vlen:${VLEN},elen:${XLEN} ${OUTPUT_STAGE1_BIN}$(shell basename $@)
	$(RISCV_GCC) -Ienv/sequencer-vector -Imacros/sequencer-vector -E ${OUTPUT_STAGE1}$(shell basename $@).S -o ${OUTPUT_STAGE1_ASM}$(shell basename $@).S
endif

tests_patch = $(addsuffix .patch, $(tests))

patching-stage2: build-spike compile-stage1
	@mkdir -p ${OUTPUT_STAGE2_PATCH}
	$(MAKE) $(tests_patch)

$(tests_patch):
	${PATCHED_SPIKE} --isa ${MARCH} --varch=vlen:${VLEN},elen:${XLEN} ${OUTPUT_STAGE1_BIN}$(shell basename $@ .patch) > ${OUTPUT_STAGE2_PATCH}$@

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
	$(RISCV_GCC) -march=${MARCH} -mabi=${MABI} $(RISCV_GCC_OPTS) -Ienv/p -Imacros/general -Tenv/p/link.ld ${OUTPUT_STAGE2}$(shell basename $@ .stage2).S -o ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
	${SPIKE} --isa ${MARCH} --varch=vlen:${VLEN},elen:${XLEN} ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
endif


clean-out:
	rm -rf out/

clean: clean-out
	go clean
	rm -rf build/

.PHONY: all \
 		build build-generator unittest \
		generate-stage1 compile-stage1 $(tests) \
		$(tests_patch) generate-stage2 compile-stage2 $(tests_stage2) \
		clean-out clean