package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

type Counter struct {
	value byte
	count int
}

type CounterSlice []*Counter
// type ByteSlice []byte

func main() {
	keylength := flag.Int("l", 1, "key length")
	alphabetSize := flag.Int("N", 256, "alphabet size")
	infile := flag.String("r", "", "file to examine")
	exampleFile := flag.String("e", "", "file to emulate")
	dumpTxp := flag.Bool("D", false, "Dump transposition table")
	dumpFreq := flag.Bool("F", false, "Dump byte-value frequencies")
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

	blocksFreq := make([]CounterSlice, len(blocks))

	for i, row := range blocks {
		blocksFreq[i] = frequencyCounter(row)
	}

	exampleBytes, exampleBuffer := readFile(*exampleFile)
	fmt.Fprintf(os.Stderr, "Read %d bytes from example\n", exampleBytes)
	fmt.Fprintf(os.Stderr, "Example buffer size %d\n", len(exampleBuffer))

	exampleFreq := frequencyCounter(exampleBuffer)

	if *dumpFreq {
		frequencyDump(&exampleFreq, "Example text")
		for i := range blocksFreq {
			frequencyDump(&blocksFreq[i], fmt.Sprintf("Block %d", i))
		}
	}

	transpose := make([][256]byte, len(blocksFreq))

	for i, freq := range blocksFreq {
		transpose[i] = createTransposition(freq, exampleFreq)
	}

	if *dumpTxp {
		for i, row := range transpose {
			fmt.Printf("Transpose %d:\nCiphertext: ", i)
			spacer := ""
			for j := range row {
				fmt.Printf("%s%3d", spacer, j)
				spacer = " "
			}
			fmt.Printf("\nCleartext:  ")
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

// Make CounterSlice sortable on byte count
func (p *CounterSlice) Len() int           { return len(*p) }
func (p *CounterSlice) Swap(i, j int)      { (*p)[i], (*p)[j] = (*p)[j], (*p)[i] }
func (p *CounterSlice) Less(i, j int) bool { return (*p)[i].count > (*p)[j].count }

/*
func (p ByteSlice) Len() int           { return len(p) }
func (p ByteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByteSlice) Less(i, j int) bool { return p[i] > p[j] }
*/

func (p Counter) String() string {
	return fmt.Sprintf("<%d, %d>", p.value, p.count)
}

func frequencyCounter(buffer []byte) CounterSlice {

	var freq CounterSlice

	freq = make([]*Counter, 256)

	for i := range freq {
		freq[i] = new(Counter)
	}

	for _, b := range buffer {
		if freq[b] == nil {
			freq[b] = new(Counter)
		}
		freq[b].count++
		freq[b].value = b
	}

	sort.Sort(&freq)

	/*
		for i := len(freq) - 1; freq[i].count == 0 && i >= 0; i-- {
			freq = freq[0:i]
		}
	*/

	return freq
}

func frequencyDump(cs *CounterSlice, phrase string) {
	fmt.Printf("%s\nCount: ", phrase)
	spacer := ""
	for i := range *cs {
		fmt.Printf("%s%3d", spacer, (*cs)[i].count)
		spacer = " "
	}
	fmt.Printf("\nValue: ")
	spacer = ""
	for i := range *cs {
		fmt.Printf("%s%3d", spacer, (*cs)[i].value)
		spacer = " "
	}
	fmt.Printf("\n\n")
}

// Make a [256]byte where the index is from
// ciphertext[N].value, and the value is example[N].value
func createTransposition(ciphertext CounterSlice, example CounterSlice) [256]byte {
	var txp [256]byte

	for i := 0; i < 256; i++ {
		txp[ciphertext[i].value] = example[i].value
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
