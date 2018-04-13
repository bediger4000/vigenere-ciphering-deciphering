package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

type CountVector [256]int

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

	theta := vectorAngle(vector1, vector2)

	fmt.Printf("%q\t%q\t%f\n", os.Args[1], os.Args[2], theta)
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

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		v[b]++
	}

	if e != nil {
		if e != io.EOF {
			return nil, errors.New("Problem reading a byte: " + e.Error())
		}
	}

	return &v, nil
}

func vectorAngle(vector1, vector2 *CountVector) float64 {
	var dotProduct float64
	var sumOfSquares1, sumOfSquares2 float64

	for i := 0; i < 256; i++ {
		dotProduct += float64(vector1[i] * vector2[i])
		sumOfSquares1 += float64(vector1[i] * vector1[i])
		sumOfSquares2 += float64(vector2[i] * vector2[i])
	}

	/*
	   	fmt.Printf("dot product: %f\n", dotProduct)
	   	fmt.Printf("sum of squares 1: %f\n", sumOfSquares1)
	   	fmt.Printf("sum of squares 2: %f\n", sumOfSquares2)
	   	fmt.Printf("sqrt1: %f\n", math.Sqrt(sumOfSquares1))
	   	fmt.Printf("sqrt2: %f\n", math.Sqrt(sumOfSquares2))
	       fmt.Printf("zork: %f\n", dotProduct/(math.Sqrt(sumOfSquares1) * math.Sqrt(sumOfSquares2)))

	   	fmt.Printf("%f\n", math.Acos(1.+math.SmallestNonzeroFloat64))
	*/

	// math.Acos() undefined for argument -1 <= x >= 1,
	// and we know that z is positive.
	cosTheta := dotProduct / (math.Sqrt(sumOfSquares1) * math.Sqrt(sumOfSquares2))

	return 1.00 - cosTheta
}
