# Vigenere Cipher Deciphering

I came across an [enciphered piece of PHP malware](https://github.com/bediger4000/php-malware-analysis/NxAaGc),
and I wanted to figure out the cleartext.
I thought the downloader might have used a Vigenere cipher.
I was wrong, it was base64 encoded, then XOR-encoded.

I read the Wikipedia page on it, and wrote some programs.

## Programs

### ic - calculate Index of Coincidence

	$ GOPATH=$PWD go build ic
	$ ./ic some.filename

`ic` calculates the Index of Coincidence of a file full of bytes.
This Index gets used in key-length estimation.
This calculates the Index of Coincidence for the entire file.

### shift - Vigenere ciphering and deciphering

	$ GOPATH=$PWD go build shift
	$ ./shift -S 56/67/99/105 -r inputfile > ciphertext
    $ ./shift -u -S 56/67/99/105 -r ciphertext > cleartext

Or alternately, for a printable ASCII key:

	$ ./shift -s '8Cci' -r inputfile > ciphertext
    $ ./shift -u -s '8Cci' -r ciphertext > cleartext


That will use a key length of 4 bytes, with the byte values 56, 67, 99 and 105.
Represented as an ASCII string, the key is "8Cci". You can use longer keys, and
key byte values from 0 to 255.

### vigkeylength - estimate key length in bytes

    $ GOPATH=$PWD go build vigkeylength
    $ ./vigkeylength filename 4 40

That will give Index of Coincidence values for keys between 4 and 40 bytes.
The key length(s) with the lowest Index are probably the correct keylengths.
I find that multiple of the key length end up as low values for some reason.

### vigkeyguess - calculate guess of cipher key

    $ GOPATH=$PWD go build vigkeyguess
    $ ./vigkeyguess -N 127 -l 5 -r ciphertext

The example finds the most likely 5-byte-long key for a file named "ciphertext",
for a 127-value (values 0 - 126) alphabet.
The longer the file the more accurate the guess will be.

Output is in a format suitable for use in the `shift` program from above, with -u flag.

### byteshisto - histogram of byte values on stdin

    $ GOPATH=$PWD go build byteshisto
    $ ./byteshisto < ciphertext > histo.dat

Build a text histogram (range 0 thru 255) of byte values appearing on stdin.
Output suitable for use in [gnuplot](http://gnuplot.info/)

### shortshisto - histogram of 2-byte values on stdin

    $ GOPATH=$PWD go build shortshisto
    $ ./shortshisto < ciphertext > histo.dat

Build a text histogram (range 0 thru 65536)
of 2-byte values appearing on stdin.
Output suitable for use in [gnuplot](http://gnuplot.info/)

### affine - [affine enciphering/deciphering](https://en.wikipedia.org/wiki/Affine_cipher)

    $ GOPATH=$PWD go build affine
    $ ./affine -m 256 -a 11 -b 120 -f cleartext | ./affine -u -m 256 -a 11 -b 120 > deciphered
	$ diff cleartext deciphered

That illustrates enciphering and deciphering in a single pipeline.
Affine ciphers seem like a variant of Vigenere ciphers,
so I wanted this to try on my mystery data.

### Differential xor encoding

### kasiski - [Kasiski method](https://en.wikipedia.org/wiki/Kasiski_examination)

`kasiski` counts distance between repeating blocks of bytes.
Key length should be a factor of the distances between repeating blocks.
This should help confirm the key length
derived from Index of Coindidence by `vigkeylength`

    $ GOPATH=$PWD go build kasiski
	$ ./kasiski -n substring-length -r filename > distances

File `distances` will have all the distances between repeating substring-length sized blocks of bytes in the file.
You probably will have to do some post-processing on the output,
like remove duplicates, sort numerically, etc etc.
The more ciphertext you've got the better this will work.
The key length will be a factor of the distances
between repeating blocks of bytes.
Some distances between repeats will almost certainly
not have key-length as a factor because of bad luck.
You'll have to weed them out.

Output (on stdout) has one row per block of bytes:

    3492:7 8008 56112 576 56088

The 7-byte-long block of bytes starting at index 3492 in the
input file has repeats with distances of 8008, 56112, 576, 56088
between the repetitions.

#### substrings - find repeating substrings of given length

To go along with Kasiski testing,
this program finds all repeating substrings of bytes of a given length
in a file,
and their offsets in the file.
This was to look for encrypted `"function "` strings in cipertext.

### rshift - make random polytransposition ciphertext

This program creates N random transpositions,
then runs a cleartext file through the transpositions.
I used this to create ciphertext that I could try hamming distance,
kasiski test and index of coincidence key length guessing on.

### keyguess - try to find most likely key

    $ GOPATH=$PWD go build keyguess
	$ ./keyguess [-N alphabet-size] [-t php] [-t english] -r filename a1/a2/a3/...  b1/b2/...  c1/c2/...

From alternative byte values at every position in key, find the "best" key.
This is based on calculating a 256-D angle between byte-value histograms of a comparison
(english or PHP text) and potentially decoded bytes from the input file.

Alternative bytes at each position are described by a string like `53/67/89/104`.
You provide one such string for each position in a key. It can have only one value, or
it can have multiple values. Each position can have a different count of alternatives.
`keyguess` iterates through all of the possible combinations of byte values for
each position, deciphering the input ciphertext with each possible key.

### vig.php - PHP Vigenere encoding.

Did Vigenere encoding in PHP to see how the code looks in PHP,
how hard it is to code, etc etc.

    $ php vig.php 'keystring' cleartextfilename > ciphertextfilename

You can use `vigkeylength`, `kasiski` and `vigkeyguess` to re-calculate
the key string.

### vectormeasure - calculate 256-dimension angle between byte value histograms

Gives some idea of the "distance" between two files full of bytes.
Calculates a byte-value histogram of each file.
It calculates an angle via dot product by
treating the byte-value histograms as 256-dimensional vectors.

    $ GOPATH=$PWD go build vectormeasure
    $ ./vectormeasure filename1 filename2
    "filename1"  "filename2"        0.013280

The closer to zero the angle is,
the "closer" the two files are by this measure.

This could maybe be improved by using more than a single
second filename,
then calculating angles between "filename1" and each subsequent file.

### chisquared - 

### txpfinder - transposition finder

Try to find N-byte transpositions for a given key length
that make resulting byte-value histograms match that of a target file.

This needs to be able to read transpositions to allow human fine tuning.
