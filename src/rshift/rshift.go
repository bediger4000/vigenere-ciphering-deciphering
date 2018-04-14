/*
 * Generate random transposition arrays, then "encipher"
 * input text with those arrays. For testing Index of
 * Coinidence key-length-guessing.
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	keylength := flag.Int("l", 2, "key length")
	infile := flag.String("r", "", "file to encipher")
	alphabetSize := flag.Int("N", 256, "Alphabet size, characters")
	dumpTxp := flag.Bool("D", false, "Print tranpositions to stderr")
	flag.Parse()

	if infile == nil || len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "%s: need file name on command line, -r <filename>\n", os.Args[0])
	}

	fin, err := os.Open(*infile)

	if err != nil {
		log.Fatal("Open Problem: %s\n", err)
	}

	rand.Seed(time.Now().UnixNano())
	txp := make([][256]int, *keylength)

	for idx := 0; idx < *keylength; idx++ {
		have := make(map[int]bool)
		for i := 0; i < *alphabetSize; i++ {
			N := rand.Intn(*alphabetSize)
			for y := have[N]; y; y = have[N] {
				N = rand.Intn(*alphabetSize)
			}
			have[N] = true
			txp[idx][i] = N
		}
	}

	if *dumpTxp {
		for idx, transpose := range txp {
			fmt.Fprintf(os.Stderr, "Transpose %d\nCipher   Clear\n", idx)
			for clear, cipher := range transpose {
				if clear >= *alphabetSize { break }
				if isPrintable(clear) && isPrintable(cipher) {
					fmt.Fprintf(os.Stderr, "*  %02x %c  %02x %c\n", cipher, cipher, clear, clear)
				} else if isPrintable(clear) {
					fmt.Fprintf(os.Stderr, "*  %02x    %02x %c\n", cipher, clear, clear)
				} else if isPrintable(cipher) {
					fmt.Fprintf(os.Stderr, "*  %02x %c  %02x\n", cipher, cipher, clear)
				} else {
					fmt.Fprintf(os.Stderr, "   %02x    %02x\n", cipher, clear)
				}
			}
		}
	}

	rdr := bufio.NewReader(fin)
	wrtr := bufio.NewWriter(os.Stdout)

	var b byte
	var e error
	var i int

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		ew := wrtr.WriteByte(byte(txp[i%*keylength][b]))
		if ew != nil {
			fmt.Fprintf(os.Stderr, "Problem writing: %s\n", ew)
		}
		i++
	}

	if e != nil && e != io.EOF {
		fmt.Fprintf(os.Stderr, "ReadByte Problem after %d bytes: %s\n", i, e)
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

func isPrintable(b int) bool {
	if b >= 32 && b <= 127 {
		return true
	}
	return false
}
