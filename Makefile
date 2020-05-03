.PHONY: default build run clean install

default: build

build:
	go build -o build/bin/hosted-model cmd/hosted-model/main.go

run:
	go run cmd/hosted-model/main.go

clean:
	go clean
	rm -rf build/bin/*
	touch build/bin/.gitkeep

install: build
	cp ./build/bin/hosted-model /usr/local/bin/hosted-model
