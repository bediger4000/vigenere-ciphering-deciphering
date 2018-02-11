/*
 * Calculate byte-wise index of coincidence for the whole file.
 *
 * IC = sum(n_i*(n_i - 1)*c / (N*(N-1))
 * n_i = observed letter counts in ciphertext
 *   N = total letters in ciphertext
 *   c = size of alphabet
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("%s: need options or a file name on command line\n", os.Args[0])
		os.Exit(1)
	}

	fileName := os.Args[1]

	fin, err := os.Open(fileName)
	if err != nil {
		fin.Close()
		log.Fatal(err)
	}
	defer fin.Close()

	rdr := bufio.NewReader(fin)

	var counts [256]int
	var byteUsed [256]bool

	var b byte
	var e error
	var byteCount int

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		byteCount++
		counts[b]++
		byteUsed[b] = true
	}

	if e != nil && e != io.EOF {
		fmt.Fprintf(os.Stderr, "Problem reading byte at offset %d: %s\n", byteCount, e)
	}

	var s1 float64

	for i := 0; i < 256; i++ {
		s1 += float64(counts[i] * (counts[i] - 1))
	}

	var N float64
	for _, appears := range byteUsed {
		if appears {
			N++
		}
	}

	ic := s1 / float64(byteCount*(byteCount-1))

	fmt.Printf("%s\t%d\t%.0f\t%.3f\n", fileName, byteCount, N, ic)
}
