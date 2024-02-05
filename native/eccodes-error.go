package native

/*
#include <eccodes.h>
*/
import "C"

func CgribGetErrorMessage(res int) string {
	// we do not need to free memory after grib_get_error_message
	return C.GoString(C.grib_get_error_message(C.int(Cint(res))))
}
