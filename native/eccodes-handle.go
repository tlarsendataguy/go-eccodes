package native

/*
#include <eccodes.h>
*/
import "C"

import (
	"io"
	"unsafe"

	"errors"
)

func CcodesHandleNewFromIndex(index CcodesIndex) (CcodesHandle, error) {
	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))

	h := C.codes_handle_new_from_index((*C.codes_index)(index), cError)
	if err != 0 {
		if err == Cint(C.CODES_END_OF_INDEX) {
			return nil, io.EOF
		}
		return nil, errors.New(CgribGetErrorMessage(int(err)))
	}
	return unsafe.Pointer(h), nil
}

func CcodesHandleNewFromFile(ctx CcodesContext, file CFILE, product int) (CcodesHandle, error) {
	var cProduct C.int

	cProduct = C.int(product)

	var err Cint
	cError := (*C.int)(unsafe.Pointer(&err))

	h := C.codes_handle_new_from_file((*C.grib_context)(ctx), (*C.FILE)(file), C.ProductKind(cProduct), cError)
	if err != 0 {
		return nil, errors.New(CgribGetErrorMessage(int(err)))
	}

	if h == nil {
		return nil, io.EOF
	}

	return unsafe.Pointer(h), nil
}

func CcodesHandleDelete(handle CcodesHandle) error {
	err := C.codes_handle_delete((*C.codes_handle)(handle))
	if err != 0 {
		return errors.New(CgribGetErrorMessage(int(err)))
	}
	return nil
}
