package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	// "sort"
	// "strconv"
	// "strings"
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

	byteCount, byteBuffer := readFile(*infile)

	fmt.Fprintf(os.Stderr, "Input file: %q\n", *infile)
	fmt.Fprintf(os.Stderr, "Read %d bytes\n", byteCount)
	fmt.Fprintf(os.Stderr, "Buffer size %d\n", len(byteBuffer))

	processBytes(*substringSize, byteBuffer, byteCount)
}

func processBytes(substringSize int, buffer []byte, bufsize int) {


	for i := 0 ; i < bufsize - substringSize; i++ {

		substring := buffer[i:substringSize+i]
		var matches []int

		for j := i+substringSize; j < bufsize - substringSize; j++ {
			if bytes.Equal(substring, buffer[j:j+substringSize]) {
				matches = append(matches, j)
				j += substringSize
			}
		}

		if len(matches) > 1 {
			fmt.Printf("%d:%d", i, substringSize)
			last := matches[0]
			for _, idx := range matches[1:] {
				fmt.Printf(" %d", idx - last)
				last = idx
			}
			fmt.Printf("\n")
		}
	}
}
