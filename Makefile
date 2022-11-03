build:
	go build -o build/main

run: build
	build/main

clean:
	go clean

.PHONY: build run clean