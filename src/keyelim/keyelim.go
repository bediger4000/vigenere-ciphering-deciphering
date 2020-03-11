package main

/*
 * Perform Vigenere cipher "key elimination" as per Wikipedia page.
 * https://en.wikipedia.org/wiki/Vigen%C3%A8re_cipher#Key_elimination
 *
 * I quote:
 * --
 * The Vigenere cipher, with normal alphabets, essentially uses modulo
 * arithmetic, which is commutative. Therefore, if the key length is known (or
 * guessed), subtracting the cipher text from itself, offset by the key length,
 * will produce the plain text encrypted with itself. If any "probable word" in
 * the plain text is known or can be guessed, its self-encryption can be
 * recognized, which allows recovery of the key by subtracting the known
 * plaintext from the cipher text. Key elimination is especially useful against
 * short messages.
 * --
 *
 * This program can do xor-encoded ciphertext key elimination as well.
 * xor is commutative - you don't even have to deal with modulo.
 * ciphertext byte i is Ci = (Mi ^ Ki)
 * ciphertext byte j == i+keylen is Cj = (Mj ^ Ki), because the same
 * key byte value gets used at position i and i+keylen
 * Ci ^ Cj == (Mi ^ Ki) ^ (Mj ^ Ki)
 * Ci ^ Cj == Mi ^ (Ki ^ (Mj ^ Ki))
 * Ci ^ Cj == Mi ^ (Ki ^ (Ki ^ Mj))
 * Ci ^ Cj == Mi ^ ((Ki ^ Ki) ^ Mj)
 * Ci ^ Cj == Mi ^ (0 ^ Mj)
 * Ci ^ Cj == Mi ^ Mj
 * j == i + keylen
 * The same procedure to find offset of a known cleartext works
 * for both Vigenere-encoded bytes and Xor-encoded bytes.
 *
 * This code assumes 1 byte per letter, that Go type byte is unsigned,
 * and that doing (byte - byte) arithmetic is implicitly modulo 256.
 * It also reads the entire ciphertext into memory, so very large texts
 * might cause problems. It should be possible to do a streaming version,
 * using something like KMP string searching as you read bytes in.
 */

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type selfModification func(buffer []byte, keylen int) []byte
type byteRecovery func(byte, byte) byte

func main() {

	keylen := flag.Int("l", 4, "key length")
	filename := flag.String("r", "", "input file")
	str := flag.String("s", "", "known cleartext")
	xorElimination := flag.Bool("x", false, "assume xor-encoded input (default Vigenere)")
	flag.Parse()

	if *filename == "" {
		log.Fatal("need input filename -r\n")
	}

	if *str == "" {
		log.Fatal("need cleartext string -s\n")
	}

	buffer, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatal(err)
	}

	var modificationFunction selfModification = selfSubtract
	var recoveryFunction byteRecovery = subtractByte
	if *xorElimination {
		modificationFunction = selfXor
		recoveryFunction = xorByte
	}

	length := len(buffer)

	// Subtract/Xor the cipher text from itself, offset by the key length.
	// The first *keylen bytes of haystack are not self-subtracted/xored,
	// so if that's the only match to the known cleartext, key elimination
	// won't work.
	haystack := modificationFunction(buffer, *keylen)

	// Subtract the known text from itself, offset by the key length.
	// The only unusual thing here is what bytes of the known text are
	// searchable.
	needle := modificationFunction([]byte(*str), *keylen)
	needle = needle[:len(*str)-*keylen]
	fmt.Fprintf(os.Stderr, "Using known text buffer %v\n", needle)
	if len(needle) <= 3 {
		fmt.Fprintf(os.Stderr, "Known text buffer is only %d bytes, may yield false positives\n", len(needle))
	}

	for i := range haystack {
		if haystack[i] == needle[0] {
			foundit := true
			for j, k := i, 0; k < len(needle) && length-j > 0; j, k = j+1, k+1 {
				if haystack[j] != needle[k] {
					foundit = false
					break
				}
			}
			if foundit {
				fmt.Printf("Found cleartext match at offset %d in ciphertext\n", i)
				// buffer[i:i+*keylen] is ciphertext that derives from
				// cleartext given on command line.
				clearbytes := []byte(*str)
				keybytes := make([]byte, *keylen)

				for j := 0; j < *keylen; j++ {
					keybytes[j] = recoveryFunction(buffer[i+j], clearbytes[j])
				}

				// The key bytes recovered may be in any alignment with cleartext.
				// Re-order the keybytes. This is harder than it should be.
				orderedkey := keybytes
				alignment := i % *keylen
				if alignment != 0 {
					orderedkey = make([]byte, alignment)
					x := *keylen - alignment
					copy(orderedkey, keybytes[x:])
					orderedkey = append(orderedkey, keybytes[:x]...)
				}
				fmt.Printf("  Key: %q\n", string(orderedkey))
			}
		}
	}
}

func selfSubtract(buffer []byte, keylen int) []byte {
	length := len(buffer)
	copied := make([]byte, length)
	for i := 0; i < length-keylen; i++ {
		copied[i] = buffer[i] - buffer[i+keylen]
	}
	return copied
}

func selfXor(buffer []byte, keylen int) []byte {
	length := len(buffer)
	copied := make([]byte, length)
	for i := 0; i < length-keylen; i++ {
		copied[i] = buffer[i] ^ buffer[i+keylen]
	}
	return copied
}

func subtractByte(a, b byte) byte {
	return a - b
}

func xorByte(a, b byte) byte {
	return a ^ b
}
