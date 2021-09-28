package wasm3

/*
#cgo CFLAGS: -Iinclude
#cgo darwin,!ios,amd64 LDFLAGS: -L${SRCDIR}/lib/macos/amd64 -lm3
#cgo ios,arm64 LDFLAGS: -L${SRCDIR}/lib/ios/arm64 -lm3 -framework Security
#cgo android,arm64 LDFLAGS: -L${SRCDIR}/lib/android/arm64 -lm3 -lm
#cgo linux,!android,amd64 LDFLAGS: -L${SRCDIR}/lib/linux/amd64 -lm3 -lm
#cgo linux,!android,arm64 LDFLAGS: -L${SRCDIR}/lib/linux/arm64 -lm3 -lm

#include "wasm3.h"
#include "m3_api_libc.h"
#include "m3_api_wasi.h"
#include "m3_env.h"
*/
import "C"

import "errors"

type Environment struct {
	ptr C.IM3Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		ptr: C.m3_NewEnvironment(),
	}
}

func (e *Environment) ParseModule(wasmBytes []byte) (*Module, error) {
	bytes := C.CBytes(wasmBytes)
	length := len(wasmBytes)
	var module C.IM3Module
	result := C.m3_ParseModule(
		e.Ptr(),
		&module,
		(*C.uchar)(bytes),
		C.uint(length),
	)
	if result != nil {
		return nil, errors.New("parse error")
	}
	return NewModule(module), nil
}

func (e *Environment) Ptr() C.IM3Environment {
	return e.ptr
}

func (e *Environment) Destroy() {
	C.m3_FreeEnvironment(e.Ptr())
}
