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

	var count [256]int

	rdr := bufio.NewReader(fin)

	var b byte
	var e error
	
	for b, e = rdr.ReadByte(); e == nil ;b, e = rdr.ReadByte() {
		count[b]++
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
		fmt.Printf("# Total bytes: %f\n", sum)
	} else {
		fmt.Printf("var vector = []int{\n")
	}

	var sumOfSquares uint64

	for i, c := range count {
		if *goArrayOutput {
			sumOfSquares += uint64(c) * uint64(c)
			fmt.Printf("%d,\n", c)
		} else {
			fmt.Printf("%d\t%d\t%.4f\n", i, c, float64(c)/sum)
		}
	}

	if *goArrayOutput {
		fmt.Printf("}\n")
		fmt.Printf("var Sum float64 = %.1f\n", sum)
		fmt.Printf("var SumOfSquares float64 = %d.0\n", sumOfSquares)
	}
}
