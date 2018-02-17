package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Break up a file into pieces by key length\n")
		fmt.Printf("Usage: %s [-f filename] -l n  n an integer, number of pieces\n", os.Args[0])
		fmt.Printf("Don't forget '--' to terminate arguments for use in pipeline\n")
		os.Exit(1)
	}

	lPtr := flag.Int("l", 4, "key length")
	dirPtr := flag.String("d", "temp", "output directory name")
	inFile := flag.String("f", "", "file to encipher")
	flag.Parse()

	fin := os.Stdin
	if *inFile != "" {
		var err error
		fin, err = os.Open(*inFile)
		if err != nil {
			log.Fatal("Problem opening file %q: %s\n", *inFile, err)
		}
	}

	rdr := bufio.NewReader(fin)

	l := *lPtr

	writers := setupOutputFiles(*dirPtr, l)

	var x byte
	var e error
	var byteCount int

	for x, e = rdr.ReadByte(); e == nil; x, e = rdr.ReadByte() {

		ew := writers[byteCount%l].WriteByte(x)
		if ew != nil {
			log.Fatalf("Problem writing a byte to file %d: %s\n", byteCount%l, ew)
		}

		byteCount++
	}

	if e != io.EOF {
		fmt.Fprintf(os.Stderr, "Problem reading a byte: %s\n", e)
	}

	for _, wrtr := range writers {
		wrtr.Flush()
	}

	fin.Close()
}

func setupOutputFiles(directory string, n int) []*bufio.Writer {
	var writers []*bufio.Writer

	e := os.Mkdir(directory, os.ModePerm)
	if e != nil {
		log.Fatalf("Trying to create directory %q: %s\n", directory, e)
	}

	for i := 0; i < n; i++ {
		filename := fmt.Sprintf("%s/%d", directory, i)
		ds, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Trying to create file %q: %s\n", filename, err)
		}
		writers = append(writers, bufio.NewWriter(ds))
	}

	return writers
}
