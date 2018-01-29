package main
import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		fmt.Printf("Histogram of byte values on stdin\n")
		fmt.Printf("Output suitable for gnuplot plots\n")
		os.Exit(1)
	}

	var count [256]int

	rdr := bufio.NewReader(os.Stdin)

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

	fmt.Printf("# Total bytes: %f\n", sum)
	fmt.Printf("# bytevalue/count of that value/proportion\n")

	for i, c := range count {
		fmt.Printf("%d\t%d\t%.4f\n", i, c, float64(c)/sum)
	}
}
