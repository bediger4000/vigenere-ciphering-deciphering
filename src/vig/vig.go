package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Vigenere enciphering\n")
		fmt.Fprintf(os.Stderr, "Usage: vig N/M/P... X filename\n")
		fmt.Fprintf(os.Stderr, "           N/M/P... rotations\n")
		fmt.Fprintf(os.Stderr, "           X alphabet size\n")
		fmt.Fprintf(os.Stderr, "           filename  file to encipher\n")
		os.Exit(1)
	}

	keyRotationsStrings := os.Args[1]
	alphabetSizeString := os.Args[2]
	inputFilename := os.Args[3]

	alphabetSize, _ := strconv.Atoi(alphabetSizeString)
	keyRotations := processShifts(keyRotationsStrings, alphabetSize)

	fmt.Fprintf(os.Stderr, "File name %q\n", inputFilename)
	fmt.Fprintf(os.Stderr, "Alphabet size %d\n", alphabetSize)
	fmt.Fprintf(os.Stderr, "Rotations: %v\n", keyRotations)

	rotations := createRotations(keyRotations, alphabetSize)
	rotationCount := len(rotations)

	fin, err := os.Open(inputFilename)
	if err != nil {
		log.Fatalf("Problem opening %q: %s\n", inputFilename, err)
	}

	rdr := bufio.NewReader(fin)
	wrtr := bufio.NewWriter(os.Stdout)

	var b byte
	var e error
	var i int

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		ew := wrtr.WriteByte(rotations[i%rotationCount][(0xff&int(b))%alphabetSize])
		if ew != nil {
			fmt.Fprintf(os.Stderr, "Problem writing byte %d: %s\n", i, ew)
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
			fmt.Printf("Problem with shift %q: %s\n", str, e)
		}
	}

	return shiftList
}

func createRotations(keyRotations []int, alphabetSize int) [][]byte {
	var rotations [][]byte

	rotations = make([][]byte, len(keyRotations))

	for i, n := range keyRotations {

		for n < 0 {
			n += alphabetSize
		}

		rotations[i] = make([]byte, alphabetSize)

		for j := 0; j < alphabetSize; j++ {
			rotations[i][j] = byte((j + n) % alphabetSize)
		}
	}

	return rotations
}
