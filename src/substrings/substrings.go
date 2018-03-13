package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

func handleFn(file *os.File) func(error) {
	return func(err error) {
		if err != nil {
			file.Close()
			log.Fatal(err)
		}
	}
}

func readFile(fileName string) (int, []byte) {

	fin, err := os.Open(fileName)
	handle := handleFn(fin)
	handle(err)

	fileinfo, err := fin.Stat()
	handle(err)

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := fin.Read(buffer)
	handle(err)

	fin.Close()

	return bytesread, buffer
}

func main() {

	infile := flag.String("r", "", "file to examine")
	substringSize := flag.Int("n", 4, "Substring length")
	flag.Parse()

	if len(os.Args) < 3 || infile == nil || *infile == "" {
		fmt.Fprintf(os.Stderr, "%s: find all repeated substrings of a given length in a file\n", os.Args[0]);
		fmt.Fprintf(os.Stderr, "%s [-n N] -r filename\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "-n N  - use N as substring length, default 4]\n")
		fmt.Fprintf(os.Stderr, "-r filename  - file in which to find substrings, no default\n")
		os.Exit(1)
	}

	byteCount, byteBuffer := readFile(*infile)

	fmt.Fprintf(os.Stderr, "Input file: %q\n", *infile)
	fmt.Fprintf(os.Stderr, "Read %d bytes\n", byteCount)
	fmt.Fprintf(os.Stderr, "Buffer size %d\n", len(byteBuffer))
	fmt.Fprintf(os.Stderr, "Substring length %d\n", *substringSize)

	processBytes(*substringSize, byteBuffer, byteCount)
}

func processBytes(substringSize int, buffer []byte, bufsize int) {

	for i := 0; i < bufsize-substringSize; i++ {

		// i holds index of where substring starts
		substring := buffer[i : substringSize+i]
		var matches []int

		// walk through the buffer one byte at a time, find
		// all the substringSize pieces of buffer that
		// match substring
		for j := i + substringSize; j < bufsize-substringSize; j++ {
			if bytes.Equal(substring, buffer[j:j+substringSize]) {
				// substring and buffer[j:something] match
				matches = append(matches, j)
			}
		}
		if len(matches) > 1 {
			// Output 1-indexed offsets for human convenience
			spacer := ""
			for _, b := range substring {
				fmt.Printf("%s%02x", spacer, b)
				spacer = " "
			}
			fmt.Printf(" %d", i+1)
			for _, n := range matches {
				fmt.Printf("\t%d", n+1)
			}
			fmt.Printf("\n")
		}
	}
}
