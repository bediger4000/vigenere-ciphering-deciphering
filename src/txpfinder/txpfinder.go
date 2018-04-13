package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Counter struct {
	value byte
	count int
}

type CounterSlice []*Counter

func main() {
	keylength := flag.Int("l", 1, "key length")
	alphabetSize := flag.Int("N", 256, "alphabet size")
	infile := flag.String("r", "", "file to examine")
	exampleFile := flag.String("e", "", "file to emulate")
	dumpTxp := flag.Bool("D", false, "Dump transposition table")
	dumpFreq := flag.Bool("F", false, "Dump byte-value frequencies")
	readTxp := flag.String("R", "", "read transposition tables from filename")
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

	transpose := make([][256]byte, len(blocksFreq))

	if *readTxp != "" {
		fd, err := os.Open(*readTxp)
		if err != nil {
			log.Fatalf("Couldn't open %q for read: %s\n", *readTxp, err)
		}
		defer fd.Close()

		var txpNumber int
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "Cipher ") {
				continue
			}
			if strings.HasPrefix(line, "Clear ") {
				continue
			}
			if strings.HasPrefix(line, "Transpose ") {
				// This will screw up on 2-digit transpose table indexes
				var e error
				txpNumber, e = strconv.Atoi(line[10:11])
				if e != nil {
					log.Fatalf("Could not read transposition number from %q: %s\n", line, e)
				}
				continue
			}
			// Get here, just read a line like: "   00    e4"
			//                                   0123456789a"
			// first number is cipher byte value, 2nd is clear byte value
			cipherByteValue, cbe := strconv.ParseUint(line[3:5], 0x10, 8)
			if cbe != nil {
				log.Fatalf("Txp %d, %q: %s\n", txpNumber, line, cbe)
			}
			clearByteValue, clbe := strconv.ParseUint(line[9:11], 0x10, 8)
			if clbe != nil {
				log.Fatalf("Txp %d, %q: %s\n", txpNumber, line, clbe)
			}
			transpose[txpNumber][cipherByteValue] = byte(0xff & clearByteValue)
		}

	} else {

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

		for i, freq := range blocksFreq {
			transpose[i] = createTransposition(freq, exampleFreq)
		}
	}

	if *dumpTxp {
		for i, row := range transpose {
			transposeDump(row, fmt.Sprintf("Transpose %d", i))
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
	fmt.Printf("%s\nCount   Value\n", phrase)
	for _, counter := range *cs {
		fmt.Printf("%04x     %02x\n", counter.count, counter.value)
	}
	fmt.Printf("\n\n")
}

func transposeDump(txp [256]byte, phrase string) {
	fmt.Printf("%s\nCipher    Clear\n", phrase)
	for in, out := range txp {
		fmt.Printf("   %02x    %02x\n", in, out)
	}
	fmt.Printf("\n\n")
}

// Make a [256]byte where the index is from
// ciphertext[N].value, and the value is example[N].value
func createTransposition(ciphertext CounterSlice, example CounterSlice) [256]byte {
	var txp [256]byte

	for i := 0; i < 256; i++ {
		if ciphertext[i].count > 0 {
			txp[ciphertext[i].value] = example[i].value
		}
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
