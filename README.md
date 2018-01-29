# Vigenere Cipher Deciphering

I came across an enciphered piece of PHP malware, and I
wanted to figure out the cleartext. I thought the downloader
might have used a Vigenere cipher.

I read the Wikipedia page on it, and wrote some programs.

## Programs

### ic - calculate Index of Coincidence

	$ go build ic.go
	$ ./ic some.filename

`ic` calculates the Index of Coincidence of a file full of bytes.
This Index gets used in key-length estimation.

### shift - Vigener ciphering and deciphering

	$ go build shift.go
	$ ./shift -S 56/67/99/105 -r inputfile > ciphertext
    $ ./shift -u -S 56/67/99/105 -r ciphertext > cleartext

That will use a key length of 4 bytes, with the byte values 56, 67, 99 and 105.
Represented as an ASCII string, the key is "8Cci". You can use longer keys, and
key byte values from 0 to 255.

This could use an ASCII key option, because a lot of times that's what you see used.

This could use an alphabet size in bytes - the Wikipedia page
has a modulo-alphabet-size operation on it, this one effectively
does mod 256

### vigkeylength - estimate key length in bytes

    $ go build vigkeylength.go
    $ ./vigkeylength filename 4 40

That will give Index of Coincidence values for keys between 4 and 40 bytes.
The key length(s) with the lowest Index are probably the correct keylengths.
I find that multiple of the key length end up as low values for some reason.

This could use an alphabet size in bytes

### vigkeyguess - calculate guess of cipher key

    $ go build vigkeyguess.go
    $ ./vigkeyguess -l 5 -r ciphertext

The example finds the most likely 5-byte-long key for a file named "ciphertext".
The longer the file the more accurate the guess will be.

Output is in a format suitable for use in the `shift` program from above, with -u flag.

This could use an alphabet size in bytes

### byteshisto - histogram of byte values on stdin

    $ go build byteshisto.go
    $ ./byteshisto < ciphertext > histo.dat

Build a text histogram (range 0 thru 255) of byte values
appearing on stdin. Output suitable for use in [gnuplot](http://gnuplot.info/)
