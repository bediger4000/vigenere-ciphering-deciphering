/*
 * Calculate byte-wise index of coincidence for the whole file.
 *
 * IC = sum(n_i*(n_i - 1)*c / (N*(N-1))
 * n_i = observed letter counts in ciphertext
 *   N = total letters in ciphertext
 *   c = letters in alphabet
 */

package main

import (
	"fmt"
	"log"
	"os"
)

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

func main() {
	byteCount, byteBuffer := readFile(os.Args[1])

	fmt.Printf("Read %d bytes\n", byteCount)
	fmt.Printf("Buffer size %d\n", len(byteBuffer))

	var counts [256]int
	var byteUsed [256]bool

	for i := 0; i < byteCount; i++ {
		counts[byteBuffer[i]]++
		byteUsed[byteBuffer[i]] = true
	}

	var s1 float64

	for i := 0; i < 256; i++ {
		s1 += float64(counts[i]*(counts[i] - 1))
	}

	var N float64
	for _, appears := range byteUsed {
		if appears { N++ }
	}

	ic := s1*N/float64(byteCount*(byteCount - 1))

	fmt.Printf("Alphabet of %.0f size\n", N)
	fmt.Printf("Index of coincidence %f\n", ic)
}
