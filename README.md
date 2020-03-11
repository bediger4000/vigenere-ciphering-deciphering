# Vigenere Cipher Deciphering

I came across an [enciphered piece of PHP malware](https://github.com/bediger4000/php-malware-analysis/NxAaGc),
and I wanted to figure out the cleartext.
I thought the downloader might have used a Vigenere cipher.
I was wrong, it was base64 encoded, then XOR-encoded.

In any case,
I read the Wikipedia pages on [Vigenere ciphers](https://en.wikipedia.org/wiki/Caesar_cipher),
and [classical ciphers](https://en.wikipedia.org/wiki/Classical_cipher)
and wrote some programs.
Here they are.

## Programs

I wrote a mix of ciphering/deciphering and analysis programs.
I learned as I went along.

### ic - calculate Index of Coincidence

	$ go build src/ic/ic.go
	$ ./ic some.filename

`ic` calculates the Index of Coincidence of a file full of bytes.
This Index gets used in key-length estimation.
This calculates the Index of Coincidence for the entire file.

### shift - Vigenere ciphering and deciphering

	$ go build src/shift/shift.go
	$ ./shift -S 56/67/99/105 -r inputfile > ciphertext
    $ ./shift -u -S 56/67/99/105 -r ciphertext > cleartext

Or alternately, for a printable ASCII key:

	$ ./shift -s '8Cci' -r inputfile > ciphertext
    $ ./shift -u -s '8Cci' -r ciphertext > cleartext


That will use a key length of 4 bytes, with the byte values 56, 67, 99 and 105.
Represented as an ASCII string, the key is "8Cci". You can use longer keys, and
key byte values from 0 to 255.

### vigkeylength - estimate key length in bytes

    $ go build src/vigkeylength/vigkeylength.go
    $ ./vigkeylength filename 4 40

That will give Index of Coincidence values for keys between 4 and 40 bytes.
The key length(s) with the lowest Index are probably the correct keylengths.
I find that multiple of the key length end up as low values for some reason.

### vigkeyguess - calculate guess of cipher key

    $ go build src/vigkeyguess/vigkeyguess.go
    $ ./vigkeyguess -N 127 -l 5 -r ciphertext

The example finds the most likely 5-byte-long key
for a file named "ciphertext",
for a 127-value (values 0 - 126) alphabet.  The
longer the file the more accurate the guess will be.

`vigkeyguess` divides the input file into keylength number of bins
by counting off bytes: every Nth byte goes in bin N.
`vigkeyguess` finds a "rotation" for each bin that
maximizes the number of valid ASCII characters for that bin.
A rotation is a value that gets added to each of the bytes
in a bin, modulo alphabet size.

Output is in a format suitable for use in the `shift` program from above, with -u flag.

### keyelim - find key of Vigenere ciphertext by "key elimination"

    $ go build src/keyelim/keyelim.go
    $ ./keyelim -l 5 -r ciphertext -s knownplaintext

Does [key elimination](https://en.wikipedia.org/wiki/Vigen%C3%A8re_cipher#Key_elimination)
for ciphertext from a true Vigenere cipher.
See [shift](https://github.com/bediger4000/vigenere-ciphering-deciphering#shift---vigenere-ciphering-and-deciphering)
in this repo for just such an enciphering.
You do have to know the length of the key, and some known plaintext that appears in
the original clear text.
The known plaintext has to have a greate length than the key,
the greater the better.

### byteshisto - histogram of byte values on stdin

    $ go build src/byteshisto/byteshisto.go
    $ ./byteshisto < ciphertext > histo.dat

Build a text histogram (range 0 thru 255) of byte values appearing on stdin.
Output suitable for use in [gnuplot](http://gnuplot.info/)

### shortshisto - histogram of 2-byte values on stdin

    $ go build src/shortshisto/shortshisto.go
    $ ./shortshisto < ciphertext > histo.dat

Build a text histogram (range 0 thru 65536)
of 2-byte values appearing on stdin.
Output suitable for use in [gnuplot](http://gnuplot.info/)

### affine - [affine enciphering/deciphering](https://en.wikipedia.org/wiki/Affine_cipher)

    $ go build src/affine/affine.go
    $ ./affine -m 256 -a 11 -b 120 -f cleartext | ./affine -u -m 256 -a 11 -b 120 > deciphered
	$ diff cleartext deciphered

That illustrates enciphering and deciphering in a single pipeline.
Affine ciphers seem like a variant of Vigenere ciphers,
so I wanted this to try on my mystery data.

### Differential xor encoding

Xor encoding where first input byte xor-ed with an initialisation byte.
Every subsequent input byte gets xor-ed with the preceding byte.
Two interpretations of that are possible: xor with the next input byte,
or xor with the input byte itself xored with the previous byte.

    $ go build src/diffxor/diffxor.go
    $ ./diffxor -N 123 -r some.file > ciphertext
    $ ./diffxor -N 123 -d -r ciphertext > cleartext

You could actually use the "-d" flag on to encode.
An invocation without "-d" would decode in that case.

### kasiski - [Kasiski method](https://en.wikipedia.org/wiki/Kasiski_examination)

`kasiski` counts distance between repeating blocks of bytes.
Key length should be a factor of the distances between repeating blocks.
This should help confirm the key length
derived from Index of Coindidence by `vigkeylength`

    $ go build src/kasiski/kasiski.go
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

    $ go build src/rshift/rshift.go

This program creates N random transpositions,
then runs a cleartext file through the transpositions.
I used this to create ciphertext that I could try hamming distance,
kasiski test and index of coincidence key length guessing on.

You can generate transposition tables for [txpfinder](#txpfinder)
to work on.

    $ ./rshift -D -l 4  -r some.file > ciphertext 2> transpositions

### keyguess - try to find most likely key

    $ go build src/keyguess/keyguess.go
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

    $ go build src/vectormeasure/vectormeasure.go
    $ ./vectormeasure filename1 filename2
    "filename1"  "filename2"        0.013280

The closer to zero the angle is,
the "closer" the two files are by this measure.

This could maybe be improved by using more than a single
second filename,
then calculating angles between "filename1" and each subsequent file.

The files don't have to be full of text: `vectormeasure` doesn't
assume anything about encoding or representation.
On the other hand, it doesn't know anything about headers or
other parts of files that are very similar from file to file.
It might give similarity measures that don't make sense for PNG
or GIF files.

### chisquared - chi-squared similarity measure of two files

    $ go build src/chisquared/chisquared.go
    $ ./chisquared filename1 filename2
    "filename1" "filename2"  523.975279

The smaller the measure,
the less difference between the two files,
with 0.00 as a minimum for identical files.

### txpfinder - transposition finder

Try to find N-byte transpositions for a given key length
that result in byte-value histograms matching those
of a target file.

    $ go build src/txpfinder/txpfinder.go

I wanted a program to find transpositions,
not just "rotations",
of characters in an enciphered file.
My `vigkeyguess` program did not decipher the malware
files correctly,
so I thought that perhaps the encoding was
jumbling the bytes in the files,
rather than just shifting them.

The idea behind `txpfinder` is to divide the input bytes
into "keylength" number of bins, byte number x going into
bin x modulo keylength.
Bytes in a given bin theoretically got jumbled using the
same transposition.

`txpfinder` creates a byte-value count of each bin,
and it also creates a byte-value count of an example file.
Because PHP source will have a different character count
than English text, one histogram does not fit all.
`txpfinder` matches a bin's byte-value counts against
the byte-value counts of the example file,
deriving a transposition that should match enciphered
bytes to cleartext bytes.

An invocation like this:

    $ ./txpfinder -N 128 -e lang/english -l 2 -r ciphertext

would put odd numbered bytes of file `ciphertext` in one bin,
even numbered in another (`-l 2`, keylength of 2).
It would create an example byte-value count of file `lang/english`,
but only bytes with values less than 128.
It would count byte values in each of the bins,
then order the byte-value counts for each of the bins,
matching them against the byte-value count for the example file.
`txpfinder` would ultimately output its best guess at cleartext.

Unfortunately, the ciphertext bytes often have an ambiguous
sort: two byte values will have exactly the same count,
which leaves `txpfinder` probably making an incorrect
correspondence between cipher byte value and example byte value.

An invocation like this:

    $ ./txpfinder -D -N 128 -e lang/english -l 2 -r ciphertext > transpose

performs the same work, but outputs 2 byte-value transpositions.
Humans can adjust the correspondence from cipher byte value
to cleartext byte value.

    $ ./txpfinder -R transpose -N 128 -e lang/english -l 2 -r ciphertext > cleartext

causes `txpfinder` to use the byte-value correspondence
in file named `transpose` to create cleartext on stdout.

Hopefully, the `rshift` program can output transpositions
that `txpfinder` can use to create cleartext.
`rshift` puts byte-value correspondence on stderr,
so that it can put ciphertext using that correspondence on stdout.

#### Transposition Finder experience

I created a 110Kb file containing English text by
concatenating README files from a large variety of projects,
mine and others.
I could only get `txpfinder` to approximately correctly
decipher transpositions of keylength 1,
using the English text that got enciphered as the example.

    $ ./rshift -N 128 -l 1  -r lang/english > ciphertext
    $ ./txpfinder -N 128 -l 1  -e lang/english -r ciphertext > cleartext

`txpfinder` can decipher keylengths of 2 or more
only approximately, often yeilding only amusing gibberish.
Apparently deciphering character transpositions takes
very large example texts and enciphered text.
