# RISC-V Vector Tests Generator

## About

This repository hosts unit tests generator for the RISC-V vector extension.

## How it works

The Spike simulator is known as the RISC-V gold standard simulator, and although we don't know how Spike is tested, it does fully support the V extension. So we added [a custom special instruction](https://github.com/ksco/riscv-vector-tests/blob/6a23892a5ab0cc72f4867cc95186b3528c99c2a0/pspike/pspike.cc#L20) to Spike, and for any test, let it automatically generate a reference result for that test. This way, we generate tests for all instructions almost automatically. Under this framework, all we have to do is write a [simple config file for each instruction](configs/).

## Prerequisite

1. `riscv64-unknown-elf-gcc` with RVV 1.0 support
2. The Spike simulator
3. Golang 1.19+
4. `riscv-pk` if you need to generate user-mode binaries

For starters, you can directly download the pre-generated tests from Github Action Artifacts.

## How to use

```
make -j$(nproc)
```

> If you have problems compiling, please refer to the build steps in [build-and-test.yml](.github/workflows/build-and-test.yml).

After `make`, you will find all the generated tests in `out/v[vlen]x[xlen][mode]/bin/stage2/`.

Options:

- `VLEN`, default is 256
- `XLEN`, default is 64, we do not support specifying ELEN yet, ELEN is consistent with XLEN
- `MODE`, default is `machine`, can be `machine`, `virtual` or `user`
- `INTEGER`, default is 0, set to 1 if you don't want float tests (i.e. for Zve32x or Zve64x)
- `PATTERN`, default is `.*`, set to a valid regex to generate the tests of your interests (e.g. `PATTERN='^v[ls].+\.v$'` to generate load/store tests)
- `TESTFLOAT3LEVEL`, default is 2, must be one of 1 or 2, testing level for testfloat3 generated cases.
- `TEST_MODE`, default is `self`, set to `cosim` if you want to generate faster tests without self-verification (to be used with co-simulators).
- `REPEAT`, default is `1`, set to greater value to repeat the same V instruction n times for a better coverage (only valid for float instructions).

For example, to generate `isa=rv32gcv varch=vlen:128,elen:32 mode=machine` tests, use `make -e VLEN=128 XLEN=32 MODE=machine -j$(nproc)`.

> Note: [single/single.go](single/single.go) generates tests directly from stage 1, suitable for targets with co-simulators (or simply use `TEST_MODE=cosim` if you're lazy).

### Nix package

This repository also provides a nix derivation with the following output provided:

- `${riscv-vector-test}/bin/*`: Generator binaries
- `${riscv-vector-test}/include/*`: Necessary headers for runtime usage
- `${riscv-vector-test}/configs/*`: Necessary runtime configs

## TODOs

- [ ] Add tests for sub-extensions (e.g. Zvamo, Zvfh)
- [ ] Support Zve64f and Zve32d
- [ ] Tail/mask agnostic support
- [ ] Fault-only-first tests

## License

This project uses third-party projects, and the licenses of these projects are attached to the corresponding directories.

The code for this project is distributed under the MIT license.

The “RISC-V” trade name is a registered trademark of RISC-V International.
