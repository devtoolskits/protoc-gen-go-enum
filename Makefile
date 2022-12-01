.PHONY: build install

build:
	go build "-ldflags=-X main.version=$(shell git describe --tags --abbrev=0)" -o protoc-gen-go-enum -trimpath .

install:
	go install .