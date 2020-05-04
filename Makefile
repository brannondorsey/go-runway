.PHONY: default build run clean install docs

default: build

build:
	go build -o build/bin/basic examples/basic/main.go
	go build -o build/bin/text-generation examples/text-generation/main.go
	go build -o build/bin/hosted-model examples/hosted-model/main.go

clean:
	go clean
	rm -rf build/bin/*
	touch build/bin/.gitkeep

install: build
	cp ./build/bin/basic /usr/local/bin/basic
	cp ./build/bin/text-generation /usr/local/bin/text-generation
	cp ./build/bin/hosted-model /usr/local/bin/hosted-model
