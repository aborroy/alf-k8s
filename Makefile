EXECUTABLE=alf-k8s
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN_ARM64=$(EXECUTABLE)_darwin_arm64
DARWIN_AMD64=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

.PHONY: all test clean

all: test build

build: windows linux darwin-arm64 darwin-amd64
	@echo version: $(VERSION)

windows: $(WINDOWS) 

linux: $(LINUX) 

darwin-arm64: $(DARWIN_ARM64) 

darwin-amd64: $(DARWIN_AMD64) 

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o $(WINDOWS) -ldflags="-s -w -X main.version=$(VERSION)"  ./main.go

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"  ./main.go

$(DARWIN_ARM64):
	env GOOS=darwin GOARCH=arm64 go build -v -o $(DARWIN_ARM64) -ldflags="-s -w -X main.version=$(VERSION)"  ./main.go

$(DARWIN_AMD64):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(DARWIN_AMD64) -ldflags="-s -w -X main.version=$(VERSION)"  ./main.go

test:
	cd test && ./test_alfresco.sh

clean:
	go clean
	rm -f $(WINDOWS) $(LINUX) $(DARWIN_AMD64) $(DARWIN_ARM64)
