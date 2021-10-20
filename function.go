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
	"math"
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
	argsLength := len(args)

	// cArgs
	cArgs := C.malloc(C.size_t(argsLength) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cArgsSlice := (*[math.MaxUint16]unsafe.Pointer)(cArgs)[:argsLength:argsLength]
	for i, v := range args {
		cArgsSlice[i] = nil

		iVal, ok := v.(int)
		if ok {
			cVal := C.int(iVal)
			cArgsSlice[i] = unsafe.Pointer(&cVal)
			continue
		}

		f32Val, ok := v.(float32)
		if ok {
			cVal := C.float(f32Val)
			cArgsSlice[i] = unsafe.Pointer(&cVal)
			continue
		}

		f64Val, ok := v.(float64)
		if ok {
			cVal := C.float(f64Val)
			cArgsSlice[i] = unsafe.Pointer(&cVal)
			continue
		}
	}
	cArgsPtr := (*unsafe.Pointer)(cArgs)

	// call function
	result := C.m3_Call(f.Ptr(), C.uint(argsLength), cArgsPtr)
	if result != nil {
		return nil, errors.New("error when calling function")
	}

	// cResults
	cResults := C.malloc(C.size_t(f.ptr.funcType.numRets) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cResultsPtr := (*unsafe.Pointer)(&cResults)

	// get results
	result = C.m3_GetResults(f.Ptr(), C.uint(f.ptr.funcType.numRets), cResultsPtr)
	if result != nil {
		return nil, errors.New("error when getting result")
	}

	// parse results
	cResultsSlice := (*[math.MaxUint16]unsafe.Pointer)(unsafe.Pointer(cResults))[:f.ptr.funcType.numRets:f.ptr.funcType.numRets]
	iResults := make([]interface{}, len(cResultsSlice))
	for i, x := range cResultsSlice {
		t := C.m3_GetRetType(f.ptr, C.uint(i))

		switch t {
		case C.c_m3Type_none:
			iResults[i] = nil
		case C.c_m3Type_i32:
			y := *(*int32)(unsafe.Pointer(&x))
			w := (interface{})(y)
			iResults[i] = w
		case C.c_m3Type_i64:
			y := *(*int64)(unsafe.Pointer(&x))
			w := (interface{})(y)
			iResults[i] = w
		case C.c_m3Type_f32:
			y := *(*float32)(unsafe.Pointer(&x))
			w := (interface{})(y)
			iResults[i] = w
		case C.c_m3Type_f64:
			y := *(*float64)(unsafe.Pointer(&x))
			w := (interface{})(y)
			iResults[i] = w
		default:
			iResults[i] = nil
		}
	}

	return iResults, nil
}
