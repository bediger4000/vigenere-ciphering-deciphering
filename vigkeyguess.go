package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
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
	keylength := flag.Int("l", 1, "key length")
	infile := flag.String("r", "", "file to examine")
	flag.Parse()

	byteCount, byteBuffer := readFile(*infile)

	fmt.Fprintf(os.Stderr, "Input file: %q\n", *infile)
	fmt.Fprintf(os.Stderr, "Read %d bytes\n", byteCount)
	fmt.Fprintf(os.Stderr, "Buffer size %d\n", len(byteBuffer))
	fmt.Fprintf(os.Stderr, "Assumed key length %d\n", *keylength)

	findKey(byteBuffer, byteCount, *keylength)
}

func findKey(cipherText []byte, cipherTextSize int, keyLength int) {

	var outputKey []int
	columns := make([][]byte, keyLength)

	bufferSize :=  cipherTextSize/keyLength
	fmt.Printf("Each column has %d bytes\n", bufferSize)

	// Got to be a clever way to just use cipherText[] in place,
	// instead of using double the memory this way.
	for i := 0; i < cipherTextSize; i++ {
		columns[i%keyLength] = append(columns[i%keyLength], cipherText[i])
	}

	for colIdx, col := range columns {

		maxCount := -1
		maxCountOffset := 0

		for offset := 0; offset < 256; offset++ {
			asciiCount := 0
			for _, b := range col {
				d := (int(b) - offset)%256

				if isAscii(byte(d)) {
					asciiCount++
				}
			}
			if asciiCount > maxCount {
				maxCount = asciiCount
				maxCountOffset = int(offset)
			}
		}

		outputKey = append(outputKey, maxCountOffset)

		fmt.Printf("column %d\t%d\t%d\t%d", colIdx, len(col), maxCount, maxCountOffset)
		if isAscii(byte(maxCountOffset)) {
			fmt.Printf("\t%c", byte(maxCountOffset))
		}
		fmt.Printf("\n")
	}

	separater := ""
	for _, offset := range outputKey {
		fmt.Printf("%s%d", separater, offset)
		separater = "/"
	}
	fmt.Printf("\n")
}

func processShifts(shifts string) []int {

	var shiftList []int

	shiftsAsStrings := strings.Split(shifts, "/")

	for _, str := range shiftsAsStrings {
		if n, e := strconv.Atoi(str); e == nil {
			shiftList = append(shiftList, n)
		} else {
			fmt.Printf("Problem with shift %q: %s\n", str, e)
		}
	}

	return shiftList
}

func isAscii(b byte) bool {
	if b == 9 || b == 10 || b == 13 || (b >= 32 && b <= 127) {
		return true
	}
	return false
}
