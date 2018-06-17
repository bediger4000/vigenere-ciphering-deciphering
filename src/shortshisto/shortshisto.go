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

	goArrayOutput := flag.Bool("g", false, "Go array output")
	infile := flag.String("r", "", "File name")
	flag.Parse()

	fin := os.Stdin
	if *infile != "" {
		var err error
		fin, err = os.Open(*infile)
		if err != nil {
			log.Fatalf("Opening input file %q: %s\n", *infile, err)
		}
	}

	var count [65536]int

	rdr := bufio.NewReader(fin)

	var b1, b2 byte
	var e error
	var x uint16

	for b1, e = rdr.ReadByte(); e == nil; b1, e = rdr.ReadByte() {
		b2, e = rdr.ReadByte()
		if e == nil {
			x = 0xffff & (uint16(b1) | (uint16(b2) << 8))
			count[x]++
		} else {
			break
		}
	}

	if e != nil {
		if e != io.EOF {
			fmt.Fprintf(os.Stderr, "Problem reading a byte: %s\n", e)
		}
	}

	var sum float64

	for _, c := range count {
		sum += float64(c)
	}

	if !*goArrayOutput {
		fmt.Printf("# Total shorts: %.0f\n", sum)
		fmt.Printf("# Value  Count  Proportion\n")
	} else {
		fmt.Printf("var vector = []int{\n")
	}

	var sumOfSquares uint64

	for i, c := range count {
		if c > 0 {
			if *goArrayOutput {
				sumOfSquares += uint64(c) * uint64(c)
				fmt.Printf("%d,\n", c)
			} else {
				fmt.Printf("%04x\t%d\t%.4f\n", i, c, float64(c)/sum)
			}
		}
	}

	if *goArrayOutput {
		fmt.Printf("}\n")
		fmt.Printf("var Sum float64 = %.1f\n", sum)
		fmt.Printf("var SumOfSquares float64 = %d.0\n", sumOfSquares)
	}
}
