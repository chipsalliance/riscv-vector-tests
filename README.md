# RISC-V Vector Tests

RISC-V V-extension 1.0 has been frozen for a while, but there is currently no serious open source test suite available, and this project tries to fill that void.

The Spike simulator is known as the RISC-V gold standard simulator, and although we don't know how Spike is tested, it does fully support the V extension. So we can make a slight modification to Spike, and for any test, let it automatically generate a reference result for that test. This way, we can generate tests for all instructions almost automatically. Under this framework, all we have to do is write a [simple config file for each instruction](configs/).

## Plan

- [ ] `[303/346]` Add tests for all insns (only basic tests, no coverage required)
- [ ] Improve test cases for existing tests
  - Add more test cases, the more, the better! 
  - Add NaN tests for float insns
  - Add Inf tests for float insns
  - ...
- [ ] Add check mechanism for CSR register
- [ ] Add V register coverage test
- [x] Support generating user mode tests (for gem5)
- [ ] Add test coverage statistics
- [ ] Add negative tests
- [ ] Add tests for sub extensions (e.g. Zvamo, Zvfh).
- [ ] Support Zve64f.
- [ ] Add simple sanity tests.

## Prerequisite

1. `riscv64-unknown-elf-gcc` with RVV 1.0 support
2. The Spike simulator
3. Golang 1.19+
4. `riscv-pk` if you need to generate user mode binaries

## How to use

```
make -j8
```

After `make`, you will find all the generated tests in `out/bin/stage2/`.

The default VLEN is 256, if you want to generate tests for a different VLEN/XLEN, you can use `make -e VLEN=512 XLEN=64 -j8`.

> NOTE:
> 1. We do not support specifying ELEN yet, ELEN is consistent with XLEN.
> 2. When changing VLEN and XLEN, you need to run `make` twice. The first run will regenerate the Makefrag file (and then fails), the second run will generate the tests.

If you want to generate user mode binaries, you can use `make -e USERMODE=1 -j8`.

If you don't want float tests (i.e. for Zve32x or Zve64x), you can use `make -e INTEGER=1 -j8`

## License

This project uses third-party projects, and the licenses of these projects are attached to the corresponding directories.

The code for this project is distributed under the MIT license.
