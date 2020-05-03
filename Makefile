.PHONY: default build run clean install docs

default: build

build:
	go build -o build/bin/text-generation examples/text-generation/main.go

run:
	go run examples/text-generation/main.go

clean:
	go clean
	rm -rf build/bin/*
	touch build/bin/.gitkeep

install: build
	cp ./build/bin/text-generation /usr/local/bin/text-generation
