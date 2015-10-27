.PHONY: all
all:
	go vet .
	go fmt .
	go test .
