package main

import (
	"flag"
	"fmt"
	codes "github.com/tlarsendataguy/go-eccodes"
	"io"
	"log"
	"runtime/debug"
	"time"
	"unsafe"
)

func main() {
	filename := flag.String("file", "", "io path, e.g. /tmp/ARPEGE_0.1_SP1_00H12H_201709290000.grib2")

	flag.Parse()

	// set filter: get 'tp' variable messages
	filter := map[string]interface{}{
		"shortNameECMF": "tp",
	}

	file, err := codes.OpenFileByPathWithFilter(*filename, filter)
	if err != nil {
		log.Fatalf("failed to open file: %s", err.Error())
	}
	defer file.Close()

	n := 0
	for {
		err = process(file, n)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to get message (#%d) from index: %s", n, err.Error())
		}

		n++
	}
}

func process(file codes.File, n int) error {
	start := time.Now()

	msg, err := file.Next()
	if err != nil {
		return err
	}
	defer func() { _ = msg.Close() }()

	log.Printf("============= BEGIN MESSAGE N%d ==========\n", n)

	startStep, err := msg.GetString("startStep")
	if err != nil {
		return fmt.Errorf("failed to get 'startStep' value: %w", err)
	}
	log.Printf("startStep = %s\n", startStep)

	endStep, err := msg.GetString("endStep")
	if err != nil {
		return fmt.Errorf("failed to get 'endStep' value: %w", err)
	}
	log.Printf("endStep = %s\n", endStep)

	stepRange, err := msg.GetString("stepRange")
	if err != nil {
		return fmt.Errorf("failed to get 'stepRange' value: %w", err)
	}
	log.Printf("stepRange = %s\n", stepRange)

	forecastTime, err := msg.GetString("forecastTime")
	if err != nil {
		return fmt.Errorf("failed to get 'forecastTime' value: %w", err)
	}
	log.Printf("forecastTime = %s\n", forecastTime)

	shortName, err := msg.GetString("shortName")
	if err != nil {
		return fmt.Errorf("failed to get 'shortName' value: %w", err)
	}
	name, err := msg.GetString("name")
	if err != nil {
		return fmt.Errorf("failed to get 'name' value: %w", err)
	}

	log.Printf("Variable = [%s](%s)\n", shortName, name)

	size, err := msg.GetLong("numberOfDataPoints")
	if err != nil {
		return fmt.Errorf("failed to get 'numberOfDataPoints' value: %w", err)
	}

	// just to measure timing
	lats, lons, vals, err := msg.DataUnsafe()
	if err != nil {
		return fmt.Errorf("failed to get data (latitudes, longitudes, values): %w", err)
	}
	defer lats.Free()
	defer lons.Free()
	defer vals.Free()

	var lat, lon, val float64
	for i := int64(0); i < size; i++ {
		uptr := uintptr(lats.Data) + uintptr(uintptr(i)*unsafe.Sizeof(lat))
		ptr := (*float64)(unsafe.Pointer(uptr))
		lat = *ptr

		uptr = uintptr(lons.Data) + uintptr(uintptr(i)*unsafe.Sizeof(lon))
		ptr = (*float64)(unsafe.Pointer(uptr))
		lon = *ptr

		uptr = uintptr(vals.Data) + uintptr(uintptr(i)*unsafe.Sizeof(val))
		ptr = (*float64)(unsafe.Pointer(uptr))
		val = *ptr

		if i < 6 {
			log.Printf("[%fx%f]=%f", lat, lon, val)
		}
	}

	log.Printf("elapsed=%.0f ms", time.Since(start).Seconds()*1000)
	log.Printf("============= END MESSAGE N%d ============\n\n", n)

	debug.FreeOSMemory()

	return nil
}
