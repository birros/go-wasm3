package wasm3

import (
	"reflect"
	"unsafe"
)

func cArrayToSlice(cArray unsafe.Pointer, arrayLen int) []unsafe.Pointer {
	var slice []unsafe.Pointer

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sh.Data = (uintptr)(cArray)
	sh.Len = arrayLen
	sh.Cap = arrayLen

	return slice
}
