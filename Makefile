VLEN = 256
OUTPUT = "out/"
CONFIGS = "configs/"

build:
	go build -o build/generator

unittest:
	go test ./...

run: build
	build/generator -VLEN ${VLEN} -output ${OUTPUT} -configs ${CONFIGS}

clean:
	go clean
	rm -rf out/
	rm -rf build/

.PHONY: build run clean