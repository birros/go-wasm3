package kernel

import (
	_ "embed"
	"errors"
	"log"

	"github.com/birros/go-wasm3"
)

//go:embed lib.wasm
var wasmBytes []byte

const (
	fnName = "sum"
)

func Start() {
	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	defer runtime.Destroy()

	module, err := runtime.ParseModule(wasmBytes)
	if err != nil {
		panic(err)
	}

	_, err = runtime.LoadModule(module)
	if err != nil {
		panic(err)
	}

	sum, err := runtime.FindFunction(fnName)
	if err != nil {
		panic(err)
	}

	log.Printf("Sum: 10 + 32")
	result, err := sum(10, 32)
	if err != nil {
		panic(err)
	}

	if len(result) == 0 {
		panic(errors.New("len(result) = 0"))
	}

	value, ok := result[0].(int32)
	if !ok {
		panic(errors.New("not int"))
	}
	log.Printf("Value: %d", value)
}
