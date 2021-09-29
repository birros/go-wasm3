all: run

GO_FILES = go.mod go.sum $(shell find . -type f -name '*.go')
LIB_FILES = go.mod go.sum $(shell find lib -type f -name '*.go')

run: build/bin
	./build/bin

build/bin: ${GO_FILES} kernel/lib.wasm
	go build -o build/bin main.go

build/ios/arm64/Kernel.framework: ${GO_FILES} kernel/lib.wasm
	gomobile bind -target ios/arm64 -o build/ios/arm64/Kernel.framework github.com/birros/go-wasm3/examples/sum/kernel

build/android/arm64/kernel.aar: ${GO_FILES} kernel/lib.wasm
	mkdir -p build/android/arm64
	gomobile bind -target android/arm64 -androidapi 28 -o build/android/arm64/kernel.aar github.com/birros/go-wasm3/examples/sum/kernel

kernel/lib.wasm: ${LIB_FILES}
	tinygo build -o kernel/lib.wasm -target wasm lib/sum.go
