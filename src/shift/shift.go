/*
 * Tool to do Vignere enciphering.
 * Key is "N/M/O/P.." where N, M, O, P are the offsets
 * to be applied to characters modulo the keylength,
 * which the program finds out from the key representation.
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	shiftPtr := flag.String("S", "", "Use N/M/... for shifts, N < 256")
	asciiShiftPtr := flag.String("s", "", "Use string of ASCII bytes for key")
	infile := flag.String("r", "", "file to encipher")
	unshift := flag.Bool("u", false, "unshift the key")
	alphabetSize := flag.Int("N", 256, "Alphabet size, characters")
	flag.Parse()

	if infile == nil || len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "%s: need file name on command line, -r <filename>\n", os.Args[0])
	}

	var keylength int
	var shifts []int

	if shiftPtr != nil && len(*shiftPtr) > 0 {
		fmt.Fprintf(os.Stderr, "Shifts: %q\n", *shiftPtr)
		shifts = processShifts(*shiftPtr, *unshift, *alphabetSize)
		keylength = len(shifts)
	}

	if asciiShiftPtr != nil && len(*asciiShiftPtr) > 0 {
		fmt.Fprintf(os.Stderr, "ASCII key: %q\n", *asciiShiftPtr)
		keylength, shifts = processAsciiKey(*asciiShiftPtr, *unshift, *alphabetSize)
	}

	fmt.Fprintf(os.Stderr, "Input file: %q\n", *infile)

	fmt.Fprintf(os.Stderr, "Key length %d: %v\n", keylength, shifts)

	fmt.Fprintf(os.Stderr, "Alphabet size %d\n", *alphabetSize)

	if keylength <= 0 {
		log.Fatal("keylength <= 0\n")
	}

	fin, err := os.Open(*infile)

	if err != nil {
		log.Fatal("Open Problem: %s\n", err)
	}

	rdr := bufio.NewReader(fin)
	wrtr := bufio.NewWriter(os.Stdout)

	var b byte
	var e error
	var i int

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		ew := wrtr.WriteByte(modulo(int(b) + shifts[i%keylength], *alphabetSize))
		if ew != nil {
			fmt.Fprintf(os.Stderr, "Problem writing: %s\n", ew)
		}
		i++
	}

	if e != nil && e != io.EOF {
		fmt.Fprintf(os.Stderr, "ReadByte Problem after %d bytes: %s\n", i, e)
	} else {
		fmt.Fprintf(os.Stderr, "Read and wrote %d bytes\n", i)
	}

	wrtr.Flush()

	fin.Close()
}

func processAsciiKey(asciiShifts string, unshift bool, alphabetSize int) (int, []int) {
	var shiftList []int

	n := 1
	if unshift {
		n = -1
	}

	bytes := []byte(asciiShifts)

	shiftList = make([]int, len(bytes))

	keylength := 0

	for i, c := range bytes {
		shiftList[i] = n * int(c)
		keylength++
	}

	return keylength, shiftList
}

func processShifts(shifts string, unshift bool, alphabetSize int) []int {

	var shiftList []int

	factor := 1

	if unshift {
		factor = -1
	}

	shiftsAsStrings := strings.Split(shifts, "/")

	for _, str := range shiftsAsStrings {
		if n, e := strconv.Atoi(str); e == nil {
			if n > alphabetSize {
				fmt.Fprintf(os.Stderr, "Shift value %d greater than slphabet size %d\n", n, alphabetSize)
			}
			shiftList = append(shiftList, factor*n)
		} else {
			fmt.Fprintf(os.Stderr, "Problem with shift %q: %s\n", str, e)
		}
	}

	return shiftList
}

func modulo(d, m int) uint8 {
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return uint8(res + m)
	}
	return uint8(res)
}
