GOBIN=$(shell go env GOPATH)/bin

.PHONY: build
build:
	CGOENABLE=0 go build -o bin/consul-query -ldflags '-w -s' main.go

install: build
	install -m 755 bin/consul-query $(GOBIN)
