.PHONY: build clean

build:
	mkdir -p ./build
	go build -o ./build/minicache ./cmd

clean:
	rm -rf ./build

