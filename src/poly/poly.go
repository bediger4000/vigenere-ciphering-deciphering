package main

/*
 * Poly-transposition enciphering program.
 */

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

	shiftPtr1 := flag.String("h", "", "Use N/M/... for horizontal extraction, N, M... < 256")
	shiftPtr2 := flag.String("v", "", "Use N/M/... for vertical extraction, N, M... < 256")
	sizePtr := flag.Int("N", 127, "Alphabet size")
	dumpTxpPtr := flag.Bool("D", false, "Dump transpositions")
	inFileName := flag.String("r", "", "file to encipher")
	flag.Parse()

	if inFileName == nil || len(os.Args) == 1 {
		log.Fatalf("%s: need file name on command line, -r <filename>\n", os.Args[0])
	}

	txp := constructTranspositions(*shiftPtr1, *shiftPtr2, *sizePtr)

	if *dumpTxpPtr {
		for i, row := range txp {
			fmt.Printf("Transpose %d:\nClear   Cipher\n", i)
			for clearByte, cipherByte := range row {
				fmt.Printf("   %02x       %02x\n", clearByte, cipherByte)
			}
		}
		os.Exit(0)
	}

	fin, err := os.Open(*inFileName)

	if err != nil {
		log.Fatal("Problem opening %q: %s\n", *inFileName, err)
	}

	rdr := bufio.NewReader(fin)
	wrtr := bufio.NewWriter(os.Stdout)

	var b byte
	var e error
	var i int

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		if int(b) < *sizePtr {
			ew := wrtr.WriteByte(txp[i%len(txp)][b])
			if ew != nil {
				log.Fatalf("Problem writing byte %d: %s\n", i, ew)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Byte %d has value %d\n", i, b)
		}
		i++
	}

	wrtr.Flush()

	if e != nil && e != io.EOF {
		log.Fatalf("ReadByte Problem after %d bytes: %s\n", i, e)
	} else {
		fmt.Fprintf(os.Stderr, "Read and wrote %d bytes\n", i)
	}
}

func constructTranspositions(horizontalString, verticalString string, alphabetSize int) [][]byte {
	var txp [][]byte

	horizontal := processShifts(horizontalString, alphabetSize)
	vertical := processShifts(verticalString, alphabetSize)

	verticalBytes := make(map[int]bool)
	for _, c := range vertical {
		verticalBytes[c] = true
	}

	horizontalBytes := make(map[int]bool)
	for _, c := range horizontal {
		horizontalBytes[c] = true
	}

	txp = make([][]byte, len(vertical))

	for i, c := range vertical {
		txp[i] = make([]byte, alphabetSize)
		txp[i][0] = byte(c)
		for j := 0; j < len(horizontal); j++ {
			txp[i][j+1] = byte(horizontal[j])
		}

		idx := len(horizontal) + 1
		for k := 0; k < alphabetSize; k++ {
			if k == c {
				continue
			}
			if horizontalBytes[k] {
				continue
			}
			if idx >= alphabetSize {
				fmt.Fprintf(os.Stderr, "i %d, k %d, idx = %d, alphbet size %d\n", i, k, idx, alphabetSize)
				os.Exit(10)
			}
			txp[i][idx] = byte(k)
			idx++
		}
	}

	return txp
}

func processShifts(shifts string, alphabetSize int) []int {

	var shiftList []int

	shiftsAsStrings := strings.Split(shifts, "/")

	for _, str := range shiftsAsStrings {
		if n, e := strconv.Atoi(str); e == nil {
			if n > alphabetSize {
				fmt.Fprintf(os.Stderr, "Shift value %d greater than slphabet size %d\n", n, alphabetSize)
			}
			shiftList = append(shiftList, n)
		} else {
			fmt.Fprintf(os.Stderr, "Problem with shift %q: %s\n", str, e)
		}
	}

	return shiftList
}
