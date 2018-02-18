# Vigenere Cipher Deciphering

I came across an enciphered piece of PHP malware, and I
wanted to figure out the cleartext. I thought the downloader
might have used a Vigenere cipher.

I read the Wikipedia page on it, and wrote some programs.

## Programs

### ic - calculate Index of Coincidence

	$ go build ic
	$ ./ic some.filename

`ic` calculates the Index of Coincidence of a file full of bytes.
This Index gets used in key-length estimation.

### shift - Vigener ciphering and deciphering

	$ go build shift
	$ ./shift -S 56/67/99/105 -r inputfile > ciphertext
    $ ./shift -u -S 56/67/99/105 -r ciphertext > cleartext

Or alternately, for a printable ASCII key:

	$ ./shift -s '8Cci' -r inputfile > ciphertext
    $ ./shift -u -s '8Cci' -r ciphertext > cleartext


That will use a key length of 4 bytes, with the byte values 56, 67, 99 and 105.
Represented as an ASCII string, the key is "8Cci". You can use longer keys, and
key byte values from 0 to 255.

### vigkeylength - estimate key length in bytes

    $ go build vigkeylength
    $ ./vigkeylength filename 4 40

That will give Index of Coincidence values for keys between 4 and 40 bytes.
The key length(s) with the lowest Index are probably the correct keylengths.
I find that multiple of the key length end up as low values for some reason.

### vigkeyguess - calculate guess of cipher key

    $ go build vigkeyguess
    $ ./vigkeyguess -N 127 -l 5 -r ciphertext

The example finds the most likely 5-byte-long key for a file named "ciphertext",
for a 127-value (values 0 - 126) alphabet.
The longer the file the more accurate the guess will be.

Output is in a format suitable for use in the `shift` program from above, with -u flag.

### byteshisto - histogram of byte values on stdin

    $ go build byteshisto
    $ ./byteshisto < ciphertext > histo.dat

Build a text histogram (range 0 thru 255) of byte values
appearing on stdin. Output suitable for use in [gnuplot](http://gnuplot.info/)

### affine - [affine enciphering/deciphering](https://en.wikipedia.org/wiki/Affine_cipher)

    $ go build affine
    $ ./affine -m 256 -a 11 -b 120 -f cleartext | ./affine -u -m 256 -a 11 -b 120 > deciphered
	$ diff cleartext deciphered
	$

That illustrates enciphering and deciphering in a single pipeline.
Affine ciphers seem like a variant of Vigenere ciphers, so I wanted this to try on my
mystery data. I don't think this is the cipher used.

### kasiski - [Kasiski method](https://en.wikipedia.org/wiki/Kasiski_examination)

`kasiski` counts distance between repeating blocks of bytes. Key length should be
a factor of the distances between repeating blocks. This should help confirm the
key length derived from Index of Coindidence by `vigkeylength`

    $ GOPATH=$PWD go build kasiski
	$ ./kasiski -n substring-length -r filename > distances

File `distances` will have all the distances between repeating substring-length sized
blocks of bytes in the file. You probably will have to do some post-processing
on the output, like remove duplicates, sort numerically, etc etc. The more ciphertext
you've got the better this will work. The key length will be a factor of the distances
between repeating blocks of bytes. Some distances between repeats will almost certainly
not have key-length as a factor because of bad luck. You'll have to weed them out.

Output (on stdout) has one row per block of bytes:

    3492:7 8008 56112 576 56088

The 7-byte-long block of bytes starting at index 3492 in the
input file has repeats with distances of 8008, 56112, 576, 56088
between the repetitions.
