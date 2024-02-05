package native

/*
#include <eccodes.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

func CcodesKeysIteratorNew(handle CcodesHandle, flags int, namespace string) CcodesKeysIterator {
	var cNamespace *C.char

	if len(namespace) > 0 {
		cNamespace = C.CString(namespace)
		defer C.free(unsafe.Pointer(cNamespace))
	}

	return unsafe.Pointer(C.codes_keys_iterator_new((*C.codes_handle)(handle), C.ulong(Culong(flags)), nil))
}

func CcodesKeysIteratorNext(kiter CcodesKeysIterator) int {
	return int(C.codes_keys_iterator_next((*C.codes_keys_iterator)(kiter)))
}

func CcodesKeysIteratorGetName(kiter CcodesKeysIterator) string {
	return C.GoString(C.codes_keys_iterator_get_name((*C.codes_keys_iterator)(kiter)))
}

func CcodesKeysIteratorDelete(kiter CcodesKeysIterator) error {
	err := C.codes_keys_iterator_delete((*C.codes_keys_iterator)(kiter))
	if err != 0 {
		return errors.New(CgribGetErrorMessage(int(err)))
	}
	return nil
}
