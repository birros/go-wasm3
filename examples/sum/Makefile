all: run

GO_FILES = go.mod go.sum $(shell find . -type f -name '*.go')
LIB_FILES = go.mod go.sum $(shell find lib -type f -name '*.go')

run: build/bin
	./build/bin

build/bin: ${GO_FILES} kernel/lib.wasm
	go build -o build/bin main.go

build/Kernel-arm64.framework: ${GO_FILES} kernel/lib.wasm
	gomobile bind -target ios/arm64 -o build/Kernel-arm64.framework github.com/birros/go-wasm3/examples/sum/kernel

build/kernel-arm64.aar: ${GO_FILES} kernel/lib.wasm
	mkdir -p build
	gomobile bind -target android/arm64 -androidapi 28 -o build/kernel-arm64.aar github.com/birros/go-wasm3/examples/sum/kernel

kernel/lib.wasm: ${LIB_FILES}
	tinygo build -o kernel/lib.wasm -target wasm lib/sum.go
