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

import (
	"errors"
	"unsafe"
)

type Function struct {
	ptr  C.IM3Function
	Name string
}

func (f *Function) Ptr() C.IM3Function {
	return f.ptr
}

func (f *Function) Call(args ...interface{}) ([](interface{}), error) {
	length := len(args)
	cArray := C.malloc(C.size_t(length) * C.size_t(unsafe.Sizeof(uintptr(0))))
	c := (*[1<<30 - 1]unsafe.Pointer)(cArray)[:length:length]
	for i, v := range args {
		a, ok := v.(int)
		if ok {
			cVal := C.int(a)
			c[i] = unsafe.Pointer(&cVal)
			continue
		}
	}
	cArrayPtr := (*unsafe.Pointer)(cArray)

	result := C.m3_Call(f.Ptr(), C.uint(length), cArrayPtr)
	if result != nil {
		return nil, errors.New("error when calling function")
	}

	fooArray := C.malloc(C.size_t(f.ptr.funcType.numRets) * C.size_t(unsafe.Sizeof(uintptr(0))))
	fooArrayPtr := (*unsafe.Pointer)(&fooArray)

	result = C.m3_GetResults(f.Ptr(), C.uint(f.ptr.funcType.numRets), fooArrayPtr)
	if result != nil {
		return nil, errors.New("error when getting result")
	}

	d := (*[1<<30 - 1]unsafe.Pointer)(unsafe.Pointer(fooArray))[:f.ptr.funcType.numRets:f.ptr.funcType.numRets]
	rets := make([]interface{}, len(d))
	for i, x := range d {
		t := C.m3_GetRetType(f.ptr, C.uint(i))
		switch t {
		case C.c_m3Type_none:
			rets[i] = nil
		case C.c_m3Type_i32:
			y := *(*int32)(unsafe.Pointer(&x))
			w := (interface{})(y)
			rets[i] = w
		case C.c_m3Type_i64:
			y := *(*int64)(unsafe.Pointer(&x))
			w := (interface{})(y)
			rets[i] = w
		case C.c_m3Type_f32:
			y := *(*float32)(unsafe.Pointer(&x))
			w := (interface{})(y)
			rets[i] = w
		case C.c_m3Type_f64:
			y := *(*float64)(unsafe.Pointer(&x))
			w := (interface{})(y)
			rets[i] = w
		default:
			rets[i] = nil
		}
	}

	return rets, nil
}
