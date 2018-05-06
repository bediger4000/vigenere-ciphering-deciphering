package main

/*
 * Differential xor, as per: https://github.com/gchq/CyberChef/issues/17
 */

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Differential XOR\n")
		fmt.Fprintf(os.Stderr, "Usage: %s N filename\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "N in decimal, initial key byte\n")
		fmt.Fprintf(os.Stderr, "Output bytes on stdout\n")
		os.Exit(1)
	}

	initialByteStr := os.Args[1]
	inputFilename := os.Args[2]

	initialByte, err := strconv.Atoi(initialByteStr)
	if err != nil {
		log.Fatal(err)
	}

	fin, err := os.Open(inputFilename)
	if err != nil {
		log.Fatalf("Problem opening %q: %s\n", inputFilename, err)
	}

	rdr := bufio.NewReader(fin)
	wrtr := bufio.NewWriter(os.Stdout)

	var b byte
	var e error
	var i int

	key := uint8(initialByte)

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		a := key ^ b
		key = b
		ew := wrtr.WriteByte(a)
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
