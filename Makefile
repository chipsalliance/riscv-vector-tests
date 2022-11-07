VLEN = 256
ELEN = 64
OUTPUT = out/
OUTPUT_STAGE1 = ${OUTPUT}tests/stage1/
OUTPUT_STAGE2 = ${OUTPUT}tests/stage2/
OUTPUT_STAGE2_PATCH = ${OUTPUT}patches/stage2/
OUTPUT_STAGE1_BIN = ${OUTPUT}bin/stage1/
OUTPUT_STAGE2_BIN = ${OUTPUT}bin/stage2/
CONFIGS = configs/

SPIKE = spike
PATCHED_SPIKE = build/spike/build/spike

include Makefrag

RISCV_PREFIX = riscv64-unknown-elf-
RISCV_GCC = $(RISCV_PREFIX)gcc
RISCV_GCC_OPTS = -static -mcmodel=medany -fvisibility=hidden -nostdlib -nostartfiles

all: compile-stage2

build: build-spike build-generator build-merger

build-generator:
	go build -o build/generator

build-merger:
	go build -o build/merger merger/merger.go

build-spike:
	git clone --depth 1 https://github.com/riscv-software-src/riscv-isa-sim.git build/spike/ 2>/dev/null || true
	cd build/spike/; \
	git am < ../../patches/0001-Modify-addi-to-generate-test-cases.patch 2>/dev/null || true; \
	mkdir -p build; \
	cd build; \
	../configure --prefix=$(.); \
	$(MAKE) -j8

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

tests_patch = $(addsuffix .patch, $(tests))

patching-stage2: compile-stage1
	$(MAKE) $(tests_patch)

$(tests_patch):
	mkdir -p ${OUTPUT_STAGE2_PATCH}
	${PATCHED_SPIKE} --isa rv64gcv --varch=vlen:${VLEN},elen:${ELEN} ${OUTPUT_STAGE1_BIN}$(shell basename $@ .patch) > ${OUTPUT_STAGE2_PATCH}$@

generate-stage2: patching-stage2
	build/merger -stage1output ${OUTPUT_STAGE1} -stage2output ${OUTPUT_STAGE2} -stage2patch ${OUTPUT_STAGE2_PATCH}

compile-stage2: generate-stage2
	@mkdir -p ${OUTPUT_STAGE2_BIN}
	$(MAKE) $(tests_stage2)

tests_stage2 = $(addsuffix .stage2, $(tests))

$(tests_stage2):
	$(RISCV_GCC) -march=rv64gv $(RISCV_GCC_OPTS) -Ienv/p -Imacros -Tenv/p/link.ld ${OUTPUT_STAGE2}$(shell basename $@ .stage2).S -o ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)
	${SPIKE} --isa rv64gcv --varch=vlen:${VLEN},elen:${ELEN} ${OUTPUT_STAGE2_BIN}$(shell basename $@ .stage2)


clean:
	go clean
	rm -rf out/
	rm -rf build/

.PHONY: all \
 		build build-generator unittest \
		generate-stage1 compile-stage1 $(tests) \
		$(tests_patch) generate-stage2 compile-stage2 $(tests_stage2) \
		clean