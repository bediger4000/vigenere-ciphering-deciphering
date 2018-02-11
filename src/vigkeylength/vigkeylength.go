/*
 * A better approach for repeating-key ciphers is to copy the ciphertext
 * into rows of a matrix having as many columns as an assumed key length,
 * then compute the average index of coincidence with each column
 * considered separately; when this is done for each possible key length,
 * the highest average I.C. then corresponds to the most likely key
 * length.[12] Such tests may be supplemented by information from the
 * Kasiski examination.
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
	"strconv"
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

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "%s: estimate Vigenere-ciphertext key length\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s <filename> N M\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Use <filename> as ciphertext, N lowest key length, M highest key length\n")
		os.Exit(1)
	}

	byteCount, byteBuffer := readFile(os.Args[1])

	fmt.Printf("Read %d bytes\n", byteCount)
	fmt.Printf("Buffer size %d\n", len(byteBuffer))

	smallKeyLen, _ := strconv.Atoi(os.Args[2])
	largeKeyLen, _ := strconv.Atoi(os.Args[3])

	counts := make([][]int, largeKeyLen)

	for i := 0; i < largeKeyLen; i++ {
		counts[i] = make([]int, 256)
	}

	for keyLen := smallKeyLen; keyLen <= largeKeyLen; keyLen++ {

		countsIndex := keyLen - smallKeyLen // upper limit for this key length

		// Zero counts indexes
		for j := 0; j < countsIndex; j++ {
			for k := 0; k < 256; k++ {
				counts[j][k] = 0
			}
		}

		// byte value counts in 0 thru keyLen-1 columns
		for j := 0; j < byteCount; j++ {
			counts[j % keyLen][byteBuffer[j]]++
		}

		// Calculate Indices of Coincidence
		var indexOfCoincidenceSum float64

		for j := 0; j < keyLen; j++ {
			var indexOfCoincidence float64
			var byteSum float64
			for i := 0; i < 256; i++ {
				cnt := float64(counts[j][i])
				byteSum += cnt
				indexOfCoincidence += cnt*(cnt - 1.)
			}
			indexOfCoincidenceSum += indexOfCoincidence/(byteSum*(byteSum-1))
		}

		aveIndexOfCoincidence := indexOfCoincidenceSum/float64(keyLen)

		fmt.Printf("%d\t%f\n", keyLen, aveIndexOfCoincidence)
	}
}
