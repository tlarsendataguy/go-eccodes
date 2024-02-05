package codes

import (
	"fmt"
	"io"
	"runtime"

	"github.com/tlarsendataguy/go-eccodes/debug"
	cio "github.com/tlarsendataguy/go-eccodes/io"
	"github.com/tlarsendataguy/go-eccodes/native"
)

type Reader interface {
	Next() (Message, error)
}

type Writer interface {
}

type File interface {
	Reader
	Writer
	Close()
}

type file struct {
	file cio.File
}

type fileIndexed struct {
	index native.CcodesIndex
}

var emptyFilter = map[string]interface{}{}

func OpenFile(f cio.File) (File, error) {
	return &file{file: f}, nil
}

func OpenFileByPathWithFilter(path string, filter map[string]interface{}) (File, error) {
	if filter == nil {
		filter = emptyFilter
	}

	var k string
	for key, value := range filter {
		if len(k) > 0 {
			k += ","
		}
		k += key
		if value != nil {
			switch value.(type) {
			case int64, int:
				k += ":l"
			case float64, float32:
				k += ":d"
			case string:
				k += ":s"
			}
		}
	}

	i, err := native.CcodesIndexNewFromFile(native.DefaultContext, path, k)
	if err != nil {
		return nil, fmt.Errorf("failed to create filtered index: %w", err)
	}

	for key, value := range filter {
		if value != nil {
			err = nil
			switch value.(type) {
			case int64:
				err = native.CcodesIndexSelectLong(i, key, value.(int64))
				if err != nil {
					err = fmt.Errorf("failed to set filter condition '%s'=%d: %w", key, value.(int64), err)
				}
			case int:
				err = native.CcodesIndexSelectLong(i, key, int64(value.(int)))
				if err != nil {
					err = fmt.Errorf("failed to set filter condition '%s'=%d: %w", key, value.(int64), err)
				}
			case float64:
				err = native.CcodesIndexSelectDouble(i, key, value.(float64))
				if err != nil {
					err = fmt.Errorf("failed to set filter condition '%s'=%f: %w", key, value.(float64), err)
				}
			case float32:
				err = native.CcodesIndexSelectDouble(i, key, float64(value.(float32)))
				if err != nil {
					err = fmt.Errorf("failed to set filter condition '%s'=%f: %w", key, value.(float64), err)
				}
			case string:
				err = native.CcodesIndexSelectString(i, key, value.(string))
				if err != nil {
					err = fmt.Errorf("failed to set filter condition '%s'='%s': %w", key, value.(string), err)
				}
			}
			if err != nil {
				native.CcodesIndexDelete(i)
				return nil, err
			}
		}
	}

	file := &fileIndexed{index: i}
	runtime.SetFinalizer(file, fileIndexedFinalizer)

	return file, nil
}

func (f *file) Next() (Message, error) {
	handle, err := native.CcodesHandleNewFromFile(native.DefaultContext, f.file.Native(), native.ProductAny)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed create new handle from file: %w", err)
	}

	return newMessage(handle), nil
}

func (f *file) Close() {
	f.file = nil
}

func (f *fileIndexed) isOpen() bool {
	return f.index != nil
}

func (f *fileIndexed) Next() (Message, error) {
	handle, err := native.CcodesHandleNewFromIndex(f.index)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create handle from index: %w", err)
	}

	return newMessage(handle), nil
}

func (f *fileIndexed) Close() {
	if f.isOpen() {
		defer func() { f.index = nil }()
		native.CcodesIndexDelete(f.index)
	}
}

func fileIndexedFinalizer(f *fileIndexed) {
	if f.isOpen() {
		debug.MemoryLeakLogger.Print("file is not closed")
		f.Close()
	}
}
