package main

import (
	"flag"
	"fmt"
	codes "go-eccodes"
	cio "go-eccodes/io"
	"io"
	"log"
	"runtime/debug"
	"time"
)

func main() {
	filename := flag.String("file", "", "io path, e.g. /tmp/ARPEGE_0.1_SP1_00H12H_201709290000.grib2")

	flag.Parse()
	fmt.Printf("parsing file %v", *filename)

	f, err := cio.OpenFile(*filename, "r")
	if err != nil {
		log.Fatalf("failed to open file on file system: %s", err.Error())
	}
	defer func() { _ = f.Close() }()

	file, err := codes.OpenFile(f)
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

	shortName, err := msg.GetString("shortName")
	if err != nil {
		return fmt.Errorf("failed to get 'shortName' value: %v", err)
	}
	name, err := msg.GetString("name")
	if err != nil {
		return fmt.Errorf("failed to get 'name' value: %w", err)
	}

	log.Printf("Variable = [%s](%s)\n", shortName, name)

	// just to measure timing
	_, _, _, err = msg.Data()
	if err != nil {
		return fmt.Errorf("failed to get data (latitudes, longitudes, values): %w", err)
	}

	log.Printf("elapsed=%.0f ms", time.Since(start).Seconds()*1000)
	log.Printf("============= END MESSAGE N%d ============\n\n", n)

	debug.FreeOSMemory()

	return nil
}
