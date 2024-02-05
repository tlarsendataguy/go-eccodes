package native

/*
#include <eccodes.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/tlarsendataguy/go-eccodes/debug"
)

const MaxStringLength = 1030
const ParameterNumberOfPoints = "numberOfDataPoints"

func CcodesGetLong(handle CcodesHandle, key string) (int64, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var value Clong
	cValue := (*C.long)(unsafe.Pointer(&value))
	err := C.codes_get_long((*C.codes_handle)(handle), cKey, cValue)
	if err != 0 {
		return 0, errors.New(CgribGetErrorMessage(int(err)))
	}

	return int64(value), nil
}

func CcodesSetLong(handle CcodesHandle, key string, value int64) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	err := C.codes_set_long((*C.codes_handle)(handle), cKey, C.long(Clong(value)))
	if err != 0 {
		return errors.New(CgribGetErrorMessage(int(err)))
	}

	return nil
}

func CcodesGetDouble(handle CcodesHandle, key string) (float64, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var value Cdouble
	cValue := (*C.double)(unsafe.Pointer(&value))
	err := C.codes_get_double((*C.codes_handle)(handle), cKey, cValue)
	if err != 0 {
		return 0, errors.New(CgribGetErrorMessage(int(err)))
	}

	return float64(value), nil
}

func CcodesSetDouble(handle CcodesHandle, key string, value float64) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	err := C.codes_set_double((*C.codes_handle)(handle), cKey, C.double(Cdouble(value)))
	if err != 0 {
		return errors.New(CgribGetErrorMessage(int(err)))
	}

	return nil
}

func CcodesGetString(handle CcodesHandle, key string) (string, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	length := CsizeT(MaxStringLength)
	cLength := (*C.size_t)(unsafe.Pointer(&length))

	err := C.codes_get_length((*C.codes_handle)(handle), cKey, cLength)
	if err != 0 {
		return "", errors.New(CgribGetErrorMessage(int(err)))
	}
	// +1 byte for '\0'
	length++

	var cBytes *C.char
	var result []byte

	if length > MaxStringLength {
		debug.MemoryLeakLogger.Printf("unnecessary memory allocation - length of '%s' value is %d greater than MaxStringLength=%d",
			key, int(length), MaxStringLength)
		result = make([]byte, length)
	} else {
		var buffer [MaxStringLength]byte
		result = buffer[:]
	}

	cBytes = (*C.char)(unsafe.Pointer(&result[0]))
	err = C.codes_get_string((*C.codes_handle)(handle), cKey, cBytes, cLength)
	if err != 0 {
		return "", errors.New(CgribGetErrorMessage(int(err)))
	}

	if length == 0 {
		return "", nil
	}
	return string(result[:length-1]), nil
}

func CcodesGribGetData(handle CcodesHandle) (latitudes []float64, longitudes []float64, values []float64, err error) {

	size, err := CcodesGetLong(handle, ParameterNumberOfPoints)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get long value of '%s': %w", ParameterNumberOfPoints, err)
	}

	latitudes = make([]float64, size)
	cLatitudes := (*C.double)(unsafe.Pointer(&latitudes[0]))

	longitudes = make([]float64, size)
	cLongitudes := (*C.double)(unsafe.Pointer(&longitudes[0]))

	values = make([]float64, size)
	cValues := (*C.double)(unsafe.Pointer(&values[0]))

	res := C.codes_grib_get_data((*C.codes_handle)(handle), cLatitudes, cLongitudes, cValues)
	if res != 0 {
		return nil, nil, nil, errors.New(CgribGetErrorMessage(int(res)))
	}

	return latitudes, longitudes, values, nil
}

func CcodesGribGetDataUnsafe(handle CcodesHandle) (latitudes unsafe.Pointer, longitudes unsafe.Pointer, values unsafe.Pointer, err error) {

	size, err := CcodesGetLong(handle, ParameterNumberOfPoints)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get long value of '%s': %w", ParameterNumberOfPoints, err)
	}

	latitudes = Cmalloc(CsizeT(size * SizeOfFloat64))
	cLatitudes := (*C.double)(latitudes)

	longitudes = Cmalloc(CsizeT(size * SizeOfFloat64))
	cLongitudes := (*C.double)(longitudes)

	values = Cmalloc(CsizeT(size * SizeOfFloat64))
	cValues := (*C.double)(values)

	res := C.codes_grib_get_data((*C.codes_handle)(handle), cLatitudes, cLongitudes, cValues)
	if res != 0 {
		Cfree(latitudes)
		Cfree(longitudes)
		Cfree(values)
		return nil, nil, nil, errors.New(CgribGetErrorMessage(int(res)))
	}

	return latitudes, longitudes, values, nil
}
