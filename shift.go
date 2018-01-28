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
	infile := flag.String("r", "", "file to encipher")
	unshift := flag.Bool("u", false, "unshift the key")
	flag.Parse()

	if infile == nil || len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "%s: need file name on command line, -r <filename>\n", os.Args[0])
	}

	fmt.Fprintf(os.Stderr, "Shifts: %q\n", *shiftPtr)
	fmt.Fprintf(os.Stderr, "Input file: %q\n", *infile)

	shifts := processShifts(*shiftPtr, *unshift)
	keylength := len(shifts)

	fmt.Fprintf(os.Stderr, "Key length %d: %v\n", keylength, shifts)

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
		ew := wrtr.WriteByte(byte((int(b) + shifts[i%keylength]) % 256))
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

func processShifts(shifts string, unshift bool) []int {

	var shiftList []int

	shiftsAsStrings := strings.Split(shifts, "/")

	for _, str := range shiftsAsStrings {
		if n, e := strconv.Atoi(str); e == nil {
			if unshift { n = -n }
			shiftList = append(shiftList, n)
		} else {
			fmt.Printf("Problem with shift %q: %s\n", str, e)
		}
	}

	return shiftList
}
