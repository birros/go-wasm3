package wasm3

/*
#cgo CFLAGS: -Iinclude
#cgo darwin,!ios,amd64 LDFLAGS: -L${SRCDIR}/lib/macos/amd64 -lm3
#cgo ios,arm64 LDFLAGS: -L${SRCDIR}/lib/ios/arm64 -lm3 -framework Security
#cgo android,arm64 LDFLAGS: -L${SRCDIR}/lib/android/arm64 -lm3 -lm
#cgo linux,!android,amd64 LDFLAGS: -L${SRCDIR}/lib/linux/amd64 -lm3 -lm
#cgo linux,!android,arm64 LDFLAGS: -L${SRCDIR}/lib/linux/arm64 -lm3 -lm
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/lib/windows/amd64 -lm3 -lm

#include "wasm3.h"
#include "m3_api_libc.h"
#include "m3_api_wasi.h"
#include "m3_env.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

type FunctionWrapper func(args ...interface{}) ([](interface{}), error)

type Config struct {
	Environment *Environment
	StackSize   uint
	EnableWASI  bool
}

type Runtime struct {
	ptr C.IM3Runtime
	cfg *Config
}

func NewRuntime(cfg *Config) *Runtime {
	ptr := C.m3_NewRuntime(
		cfg.Environment.Ptr(),
		C.uint(cfg.StackSize),
		nil,
	)
	return &Runtime{
		ptr: ptr,
		cfg: cfg,
	}
}

func (r *Runtime) FindFunction(funcName string) (FunctionWrapper, error) {
	var f C.IM3Function
	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))
	result := C.m3_FindFunction(
		&f,
		r.Ptr(),
		cFuncName,
	)
	if result != nil {
		return nil, errors.New("function lookup failed")
	}
	fn := &Function{
		ptr: (C.IM3Function)(f),
	}
	return FunctionWrapper(fn.Call), nil
}

func (r *Runtime) LoadModule(module *Module) (*Module, error) {
	result := C.m3_LoadModule(
		r.Ptr(),
		module.Ptr(),
	)
	if result != nil {
		return nil, errors.New("load error")
	}
	result = C.m3_LinkSpecTest(r.Ptr().modules)
	if result != nil {
		return nil, errors.New("LinkSpecTest failed")
	}
	if r.cfg.EnableWASI {
		C.m3_LinkWASI(r.Ptr().modules)
	}
	return module, nil
}

func (r *Runtime) ParseModule(wasmBytes []byte) (*Module, error) {
	return r.cfg.Environment.ParseModule(wasmBytes)
}

func (r *Runtime) Ptr() C.IM3Runtime {
	return r.ptr
}

func (r *Runtime) Destroy() {
	C.m3_FreeRuntime(r.Ptr())
	r.cfg.Environment.Destroy()
}
