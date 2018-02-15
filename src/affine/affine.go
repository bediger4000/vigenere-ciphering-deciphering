/* Affine cipher */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Find greatest common divisor, Euclidean algorithm\n")
		fmt.Printf("Usage: %s -a n -b m [-m M]    n, m, M integers\n", os.Args[0])
		fmt.Printf("n, M have to be coprime. M is alphabet size\n")
		fmt.Printf("Don't forget '--' to terminate arguments for use in pipeline\n")
		os.Exit(1)
	}

	decipherPtr := flag.Bool("u", false, "decipher input")
	aPtr := flag.Int("a", 12, "a key")
	bPtr := flag.Int("b", 1, "b key")
	mPtr := flag.Int("m", 127, "alphabet size")
	inFile := flag.String("f", "", "file to encipher")
	flag.Parse()

	if !coprime(*mPtr, *aPtr) {
		fmt.Fprintf(os.Stderr, "a key (%d) and alphabet size (%d) must be coprime\n", *aPtr, *mPtr)
		os.Exit(1)
	}

	fin := os.Stdin
	if *inFile != "" {
		var err error
		fin, err = os.Open(*inFile)
		if err != nil {
			log.Fatal("Problem opening file %q: %s\n", *inFile, err)
		}
	}

	rdr := bufio.NewReader(fin)
	wrtr := bufio.NewWriter(os.Stdout)

	a, b, m := *aPtr, *bPtr, *mPtr
	decipher := *decipherPtr

	/* enciphering */
	fn := func (z byte) uint8 {
		return modulo(a*int(z) + b, m)
	}

	if decipher {
		mmi := modularMultiplicativeInverse(a, m)
		/* replace cipher with decipher */
		fn = func (z byte) uint8 {
			return modulo(mmi*(int(z) - b), m)
		}
	}

	var x byte
	var e error

	for x, e = rdr.ReadByte(); e == nil; x, e = rdr.ReadByte() {

		ew := wrtr.WriteByte(fn(x))
		if ew != nil {
			log.Fatalf("Problem writing a byte: %s\n", ew)
		}

	}

	if e != io.EOF {
		fmt.Fprintf(os.Stderr, "Problem reading a byte: %s\n", e)
	}

	wrtr.Flush()

	fin.Close()
}

func modulo(d, m int) uint8 {
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return uint8(res + m)
	}
	return uint8(res)
}

func coprime(a, b int) bool {

	for {
		if a < b {
			a, b = b, a
		}
		r := a % b
		if r == 0 {
			if b == 1 {
				return true
			}
			return false
		}
		a, b = b, r
	}
	return true
}

/* Find modular multiplicative inverse by brute force */
func modularMultiplicativeInverse(a, m int) (mmi int) {

	for mmi = 0; mmi < m; mmi++ {
		if (mmi*a)%m == 1 {
			fmt.Fprintf(os.Stderr, "Modular multiplicative inverse of %d: %d\n", a, mmi)
			return mmi
		}
	}
	return mmi
}
