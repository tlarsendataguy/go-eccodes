package native

/*
#include <eccodes.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func CcodesIndexNewFromFile(ctx CcodesContext, filename string, keys string) (CcodesIndex, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cKeys := C.CString(keys)
	defer C.free(unsafe.Pointer(cKeys))

	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))
	idx := C.codes_index_new_from_file((*C.codes_context)(ctx), cFilename, cKeys, cError)
	if err != 0 {
		return nil, errors.New(CgribGetErrorMessage(int(err)))
	}
	return unsafe.Pointer(idx), nil
}

func Ccodes_index_new(ctx CcodesContext, keys string) (CcodesIndex, error) {
	cKeys := C.CString(keys)
	defer C.free(unsafe.Pointer(cKeys))

	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))
	idx := C.codes_index_new((*C.codes_context)(ctx), cKeys, cError)
	if idx == nil {
		return nil, errors.New(CgribGetErrorMessage(int(err)))
	}
	return unsafe.Pointer(idx), nil
}

func CcodesIndexSelectDouble(index CcodesIndex, key string, value float64) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	err := C.codes_index_select_double((*C.codes_index)(index), cKey, C.double(Cdouble(value)))
	if err != 0 {
		return errors.New(CgribGetErrorMessage(int(err)))
	}
	return nil
}

func CcodesIndexSelectLong(index CcodesIndex, key string, value int64) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	err := C.codes_index_select_long((*C.codes_index)(index), cKey, C.long(Clong(value)))
	if err != 0 {
		return errors.New(CgribGetErrorMessage(int(err)))
	}
	return nil
}

func CcodesIndexSelectString(index CcodesIndex, key string, value string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	err := C.codes_index_select_string((*C.codes_index)(index), cKey, cValue)
	if err != 0 {
		return errors.New(CgribGetErrorMessage(int(err)))
	}
	return nil
}

func CcodesIndexDelete(index CcodesIndex) {
	C.codes_index_delete((*C.codes_index)(index))
}
