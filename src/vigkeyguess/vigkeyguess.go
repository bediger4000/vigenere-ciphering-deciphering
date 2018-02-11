package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type rating struct {
	count int
	offset int
}

const bestN = 4

type ratings []rating

func (a ratings) Len() int { return len(a) }
func (a ratings) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ratings) Less(i, j int) bool { return a[i].count > a[j].count }

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
	alphabetSize := flag.Int("N", 256, "alphabet size")
	infile := flag.String("r", "", "file to examine")
	flag.Parse()

	byteCount, byteBuffer := readFile(*infile)

	fmt.Fprintf(os.Stderr, "Input file: %q\n", *infile)
	fmt.Fprintf(os.Stderr, "Read %d bytes\n", byteCount)
	fmt.Fprintf(os.Stderr, "Buffer size %d\n", len(byteBuffer))
	fmt.Fprintf(os.Stderr, "Assumed key length %d\n", *keylength)

	findKey(byteBuffer, byteCount, *keylength, *alphabetSize)
}

func findKey(cipherText []byte, cipherTextSize int, keyLength int, alphabetSize int) {

	columns := make([][]byte, keyLength)

	bufferSize := cipherTextSize / keyLength
	fmt.Printf("Each column has %d bytes\n", bufferSize)

	// Got to be a clever way to just use cipherText[] in place,
	// instead of using double the memory this way.
	for i := 0; i < cipherTextSize; i++ {
		columns[i%keyLength] = append(columns[i%keyLength], cipherText[i])
	}

	var keyBytes [bestN][]byte

	for colIdx, col := range columns {

		var  highestRated ratings;

		for offset := 0; offset < alphabetSize; offset++ {
			asciiCount := 0
			for _, b := range col {
				d := modulo(int(b) - offset, alphabetSize)

				if isAscii(byte(d)) {
					asciiCount++
				}
			}
			highestRated = append(highestRated, rating{count: asciiCount, offset: offset})
			sort.Sort(highestRated)
			if len(highestRated) > bestN {
				highestRated = highestRated[:bestN]
			}
		}

		fmt.Printf("column %d\t%d\t", colIdx, len(col))
		for i, m := range highestRated {
			fmt.Printf("%d: ", m.count)
			if isAscii(byte(m.offset)) {
				fmt.Printf("'%c',\t", m.offset);
			} else {
				fmt.Printf("%d,\t", m.offset);
			}
			keyBytes[i] = append(keyBytes[i], byte(m.offset))
		}
		fmt.Printf("\n")
	}

	for _, byteString := range keyBytes {
		if len(byteString) == keyLength {
			fmt.Printf("%q\n", string(byteString))
		} 
		separator := ""
		for _, b := range byteString {
			fmt.Printf("%s%d", separator, b)
			separator = "/"
		}
		fmt.Printf("\n")
	}
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
	if b == '\t' || b == '\n' || b == '\r' || (b >= 32 && b <= 127) {
		return true
	}
	return false
}

func modulo(d, m int) uint8 {
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return uint8(res + m)
	}
	return uint8(res)
}
