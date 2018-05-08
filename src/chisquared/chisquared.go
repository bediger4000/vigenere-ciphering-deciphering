package main

/*
 * Chi-squared similarity measure.
 * http://practicalcryptography.com/cryptanalysis/text-characterisation/chi-squared-statistic/
 */

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type CountVector struct {
	count *[256]int
	total float64
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "%s: measure vector between byte histograms\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s filename1 filename2\n", os.Args[0])
		os.Exit(1)
	}

	vector1, e1 := NewCountVector(os.Args[1])
	if e1 != nil {
		log.Fatal(e1)
	}
	vector2, e2 := NewCountVector(os.Args[2])
	if e2 != nil {
		log.Fatal(e2)
	}

	measure := chiSquared(vector1, vector2)

	fmt.Printf("%q\t%q\t%f\n", os.Args[1], os.Args[2], measure)
}

func NewCountVector(fileName string) (*CountVector, error) {

	fin, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("Could not open \"" + fileName + "\" for read: " + err.Error())
	}
	rdr := bufio.NewReader(fin)
	var b byte
	var e error
	var v CountVector

	v.count = new([256]int)

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		v.count[b]++
		v.total++
	}

	if e != nil {
		if e != io.EOF {
			return nil, errors.New("Problem reading a byte: " + e.Error())
		}
	}

	return &v, nil
}

func chiSquared(actual, expected *CountVector) float64 {

	var sumOfSquares float64

	ratio := actual.total/expected.total

	for idx, count := range *(expected.count) {

		if count != 0 {
			e := float64(count)*ratio
			n := float64(actual.count[idx]) - e
			sumOfSquares += n*n/e
		}
	}

	return sumOfSquares
}
