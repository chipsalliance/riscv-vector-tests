# RISC-V Vector Tests

The RISC-V V-extension 1.0 has been frozen for some time, but there is currently no public test set available, this project attempts to fill that gap.

The Spike simulator is known as the RISC-V gold standard simulator, and although we don't know how Spike is tested, it does fully support the V extension. So we can make a slight modification to Spike, and for any test, let it automatically generate a reference result for that test. This way, we can generate tests for all instructions almost automatically. Under this framework, all we have to do is write a simple config file for each instruction.

## How to use

```
make
```

After `make`, you will find all the generated tests in `out/bin/stage2/`.
