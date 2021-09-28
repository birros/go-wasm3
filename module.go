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

type Module struct {
	ptr          C.IM3Module
	numFunctions int
}

func NewModule(ptr C.IM3Module) *Module {
	return &Module{
		ptr:          ptr,
		numFunctions: -1,
	}
}

func (m *Module) Ptr() C.IM3Module {
	return m.ptr
}
