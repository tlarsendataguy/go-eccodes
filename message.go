package codes

import (
	"math"
	"runtime"

	"github.com/tlarsendataguy/go-eccodes/debug"
	"github.com/tlarsendataguy/go-eccodes/native"
)

type Message interface {
	isOpen() bool

	GetString(key string) (string, error)

	GetLong(key string) (int64, error)
	SetLong(key string, value int64) error

	GetDouble(key string) (float64, error)
	SetDouble(key string, value float64) error

	Data() (latitudes []float64, longitudes []float64, values []float64, err error)
	DataUnsafe() (latitudes *Float64ArrayUnsafe, longitudes *Float64ArrayUnsafe, values *Float64ArrayUnsafe, err error)

	Close() error
}

type message struct {
	handle native.CcodesHandle
}

func newMessage(h native.CcodesHandle) Message {
	m := &message{handle: h}
	runtime.SetFinalizer(m, messageFinalizer)

	// set missing value to NaN
	_ = m.SetDouble(parameterMissingValue, math.NaN())

	return m
}

func (m *message) isOpen() bool {
	return m.handle != nil
}

func (m *message) GetString(key string) (string, error) {
	return native.CcodesGetString(m.handle, key)
}

func (m *message) GetLong(key string) (int64, error) {
	return native.CcodesGetLong(m.handle, key)
}

func (m *message) SetLong(key string, value int64) error {
	return native.CcodesSetLong(m.handle, key, value)
}

func (m *message) GetDouble(key string) (float64, error) {
	return native.CcodesGetDouble(m.handle, key)
}

func (m *message) SetDouble(key string, value float64) error {
	return native.CcodesSetDouble(m.handle, key, value)
}

func (m *message) Data() (latitudes []float64, longitudes []float64, values []float64, err error) {
	return native.CcodesGribGetData(m.handle)
}

func (m *message) DataUnsafe() (latitudes *Float64ArrayUnsafe, longitudes *Float64ArrayUnsafe, values *Float64ArrayUnsafe, err error) {
	lats, lons, vals, err := native.CcodesGribGetDataUnsafe(m.handle)
	if err != nil {
		return nil, nil, nil, err
	}
	return newFloat64ArrayUnsafe(lats), newFloat64ArrayUnsafe(lons), newFloat64ArrayUnsafe(vals), nil
}

func (m *message) Close() error {
	defer func() { m.handle = nil }()
	return native.CcodesHandleDelete(m.handle)
}

func messageFinalizer(m *message) {
	if m.isOpen() {
		debug.MemoryLeakLogger.Print("message is not closed")
		_ = m.Close()
	}
}
