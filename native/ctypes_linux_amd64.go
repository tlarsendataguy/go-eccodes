package native

/*
#cgo LDFLAGS: -leccodes -lpng -laec -lopenjp2 -lz -lm
*/
import "C"

type Cint = int32
type Clong = int64
type Culong = uint64
type Cdouble = float64
type CsizeT = int64
