package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

type Counter [256][2]byte

func main() {
	keylength := flag.Int("l", 1, "key length")
	alphabetSize := flag.Int("N", 256, "alphabet size")
	infile := flag.String("r", "", "file to examine")
	exampleFile := flag.String("e", "", "file to emulate")
	dumpTxp := flag.Bool("D", false, "Dump transposition table")
	flag.Parse()

	byteCount, byteBuffer := readFile(*infile)

	fmt.Fprintf(os.Stderr, "Input file: %q\n", *infile)
	fmt.Fprintf(os.Stderr, "Alphabet size: %v\n", *alphabetSize)
	fmt.Fprintf(os.Stderr, "Read %d bytes\n", byteCount)
	fmt.Fprintf(os.Stderr, "Buffer size %d\n", len(byteBuffer))
	fmt.Fprintf(os.Stderr, "Assumed key length %d\n", *keylength)

	var blocks [][]byte

	blocks = make([][]byte, *keylength)

	for i, b := range byteBuffer {
		blocks[i%*keylength] = append(blocks[i%*keylength], b)
	}

	for i, row := range blocks {
		fmt.Fprintf(os.Stderr, "buffer %d has %d bytes\n", i, len(row))
	}

	blocksFreq := make([]*Counter, len(blocks))

	for i, row := range blocks {
		blocksFreq[i] = frequencyCounter(row)
	}

	exampleBytes, exampleBuffer := readFile(*exampleFile)
	fmt.Fprintf(os.Stderr, "Read %d bytes from example\n", exampleBytes)
	fmt.Fprintf(os.Stderr, "Example buffer size %d\n", len(exampleBuffer))

	exampleFreq := frequencyCounter(exampleBuffer)
	if *dumpTxp {
		fmt.Printf("Example Frequencies\n")
		spacer := ""
		for i := range exampleFreq {
			fmt.Printf("%s%3d", spacer, exampleFreq[i][0])
			spacer = " "
		}
		fmt.Printf("\n")
		spacer = ""
		for i := range exampleFreq {
			fmt.Printf("%s%3d", spacer, exampleFreq[i][1])
			spacer = " "
		}
		fmt.Printf("\n\n")
	}

	transpose := make([][256]byte, len(blocksFreq))

	for i, freq := range blocksFreq {
		transpose[i] = createTransposition(freq, exampleFreq)
	}

	if *dumpTxp {
		for i, row := range transpose {
			fmt.Printf("Transpose %d:\n", i)
			spacer := ""
			for j := range row {
				fmt.Printf("%s%3d", spacer, j)
				spacer = " "
			}
			fmt.Printf("\n")
			spacer = ""
			for _, b := range row {
				fmt.Printf("%s%3d", spacer, b)
				spacer = " "
			}
			fmt.Printf("\n")
		}
		os.Exit(0)
	}

	wrtr := bufio.NewWriter(os.Stdout)

	for i, b := range byteBuffer {
		clearbyte := transpose[i%*keylength][b]
		ew := wrtr.WriteByte(clearbyte)
		if ew != nil {
			log.Fatalf("Problem writing byte %d: %s\n", i, ew)
		}
	}

	wrtr.Flush()
}

// Make *Counter sortable on byte count
func (p *Counter) Len() int { return len(p) }
func (p *Counter) Swap(i, j int) {
	p[i][0], p[j][0] = p[j][0], p[i][0]
	p[i][1], p[j][1] = p[j][1], p[i][1]
}
func (p *Counter) Less(i, j int) bool { return p[i][1] < p[j][1] }

func frequencyCounter(buffer []byte) *Counter {

	var freq Counter

	for _, b := range buffer {
		freq[b][1]++
		freq[b][0] = b
	}

/*
	for i, pair := range freq {
		if pair[1] == 0 {
			freq[i][0] = '_'  // Put in '_' for all uncounted values
		}
	}
*/

	sort.Sort(&freq)

	return &freq
}

// both cipherText and example should comprise *Counter
// instances. Make a [256]byte where the index is from
// ciphertext[N][0], and the value is example[N][0]
func createTransposition(ciphertext *Counter, example *Counter) [256]byte {
	var txp [256]byte

	for i := 0; i < 256; i++ {
		txp[ciphertext[i][0]] = example[i][0]
	}

	return txp
}

func handleFn(file *os.File) func(error) {
	return func(err error) {
		if err != nil {
			file.Close()
			log.Fatal(err)
		}
	}
}
func readFile(fileName string) (int, []byte) {

	fin, err := os.Open(fileName)
	handle := handleFn(fin)
	handle(err)

	fileinfo, err := fin.Stat()
	handle(err)

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := fin.Read(buffer)
	handle(err)

	fin.Close()

	return bytesread, buffer
}
