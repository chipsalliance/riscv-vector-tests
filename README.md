# RISC-V Vector Tests

The RISC-V V-extension 1.0 has been frozen for some time, but there is currently no public test set available, this project attempts to fill that gap.

The Spike simulator is known as the RISC-V gold standard simulator, and although we don't know how Spike is tested, it does fully support the V extension. So we can make a slight modification to Spike, and for any test, let it automatically generate a reference result for that test. This way, we can generate tests for all instructions almost automatically. Under this framework, all we have to do is write a simple config file for each instruction.

## How to use

```
make -j8
```

After `make`, you will find all the generated tests in `out/bin/stage2/`.

The default VLEN is 256, if you want to generate tests for a different VLEN/ELEN, you can use `make -e VLEN=512 ELEN=128 -j8`.



## How it works

It uses Spike to generate reference results, which are then combined into tests to form the final result.



## License

This project uses third-party projects, and the licenses of these projects are attached to the corresponding directories.

The code for this project is distributed under the MIT license.
