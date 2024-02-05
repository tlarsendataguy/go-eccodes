package native

/*
#include <eccodes.h>
*/
import "C"
import "unsafe"

var DefaultContext = CcodesContextGetDefault()

func CcodesContextGetDefault() CcodesContext {
	ctx := C.codes_context_get_default()
	return unsafe.Pointer(ctx)
}

func CcodesContextDelete(ctx CcodesContext) {
	C.codes_context_delete((*C.codes_context)(ctx))
}
