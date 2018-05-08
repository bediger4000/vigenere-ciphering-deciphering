package main

/*
 * Differential xor, as per: https://github.com/gchq/CyberChef/issues/17
 */

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	decode := flag.Bool("d", false, "differentially decode")
	infile := flag.String("r", "", "file to differentially xor")
	N := flag.Int("N", 45, "first byte value of key")
	flag.Parse()

	if infile == nil {
		log.Fatalf("%s: need -r <filename> on command line\n", os.Args[0])
	}
	if *N < 0 || *N > 255  {
		log.Fatalf("%s: -N <number> must be between 0 and 255\n", os.Args[0])
	}

	fin, err := os.Open(*infile)
	if err != nil {
		log.Fatalf("Problem opening %q: %s\n", *infile, err)
	}

	rdr := bufio.NewReader(fin)
	wrtr := bufio.NewWriter(os.Stdout)

	var b byte
	var e error
	var i int

	key := uint8(*N)

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		a := key ^ b
		if *decode {
			key = b
		} else {
			key = a
		}
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
