package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Vector struct {
	vector       []int
	sumOfSquares float64
}

type Digit struct {
	values  []byte
	current int
}

type Digits []*Digit

func main() {

	alphabetSize := flag.Int("N", 256, "Alphabet size, characters")
	fileName := flag.String("r", "", "input file name")
	fileType := flag.String("t", "php", "input file type guess, 'php' or 'english'")
	flag.Parse()

	if flag.NArg() <= 0 {
		log.Fatalf("Need possible keybytes on command line (1a/1b/1c  2a/2b ...)\n")
	}

	// Read in the possible bytes for each position in key
	var keyBytes Digits
	a := flag.Args()
	for _, str := range a {
		d := new(Digit)
		d.values = convertKeybytes(str)
		keyBytes = append(keyBytes, d)
	}

	var compareVector Vector

	if *fileType == "php" {
		compareVector.vector = phpVector
		compareVector.sumOfSquares = phpSumOfSquares
	} else {
		compareVector.vector = englishVector
		compareVector.sumOfSquares = englishSumOfSquares
	}

	_, ciphertext := readFile(*fileName)

	findAngles(ciphertext, &compareVector, keyBytes, *alphabetSize)
}

type KeyKeeper struct {
	theta float64
	key []int
}

type KeyKeepers []*KeyKeeper

func (p KeyKeepers) Len() int { return len(p) }
func (p KeyKeepers) Swap(i, j int) { p[i], p[j] = p[j], p[i]}
func (p KeyKeepers) Less(i, j int) bool { return p[i].theta < p[j].theta }

func findAngles(ciphertext []byte, compareVector *Vector, keyBytes Digits, alphabetSize int) {

	var keyCount uint64 = 1
	for _, d := range keyBytes {
		keyCount *= uint64(len(d.values))
	}
	fmt.Fprintf(os.Stderr, "Looking at %d possible keys\n", keyCount)

	var minTheta = math.MaxFloat32

	var keykeeper KeyKeepers
	var quit bool
	var key []int
	var bestKey []int
	var count int

	start := time.Now()

	for !quit {
		count++

		if count%1000 == 0 {
			elapsed := time.Now().Sub(start)
			fmt.Fprintf(os.Stderr, "%d keys checked, min angle %.1f, %v elapsed\n", count, minTheta, elapsed)
		}

		key, quit = keyBytes.enumerateDigits()

		decodedVector := countDecodedBytes(ciphertext, key, alphabetSize)

		theta := vectorAngle(decodedVector, compareVector)
		// traceOutput(theta, key)

		if theta < minTheta {
			minTheta, bestKey = theta, key
		}

		k := new(KeyKeeper)
		k.theta = theta
		k.key = make([]int, len(key))
		copy(k.key, key)
		keykeeper = append(keykeeper, k)
	}

	fmt.Printf("Best key at %f:\n", minTheta)
	for _, b := range bestKey {
		fmt.Printf("\t%d", b)
		if isAscii(byte(b)) {
			fmt.Printf("\t%c", b)
		}
		fmt.Printf("\n")
	}

	sort.Sort(keykeeper)
	n := 0
	for _, k := range keykeeper {
		fmt.Printf("%.4f\t%v\n", k.theta, k.key)
		n++
		if n >= 10 { break }
	}

}

func traceOutput(theta float64, key []int) {
	fmt.Fprintf(os.Stderr, "%v: %.4f\n", key, theta)
}

func isAscii(b byte) bool {
	if b == '\t' || b == '\n' || b == '\r' || (b >= 32 && b <= 127) {
		return true
	}
	return false
}

func countDecodedBytes(ciphertext []byte, key []int, alphabetSize int) *Vector {

	keylength := len(key)

	var encoded Vector
	encoded.vector = make([]int, 256)

	for i, x := range ciphertext {
		encoded.vector[modulo(int(x)+key[i%keylength], alphabetSize)]++
	}

	for _, x := range encoded.vector {
		encoded.sumOfSquares += float64(x * x)
	}

	return &encoded
}

// Return configuration ([]int) and a bool that
// indicates "finished" if true
func (p *Digits) enumerateDigits() ([]int, bool) {
	var r []int

	carry := true

	for _, digit := range *p {

		r = append(r, int(digit.values[digit.current]))

		if carry {
			digit.current++
			carry = false
		}

		if digit.current == len(digit.values) {
			digit.current = 0
			carry = true
		}
	}

	return r, carry
}

func modulo(d, m int) uint8 {
	res := d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return uint8(res + m)
	}
	return uint8(res)
}

/*
func readFile(filename string) ([]int, float64, float64) {

	fin := os.Stdin
	if filename != "" {
		var err error
		fin, err = os.Open(filename)
		if err != nil {
			log.Fatalf("Opening input file %q: %s\n", filename, err)
		}
	}

	rdr := bufio.NewReader(fin)

	var b byte
	var e error

	vector := make([]int, 256)

	for b, e = rdr.ReadByte(); e == nil; b, e = rdr.ReadByte() {
		vector[b]++
	}

	if e != nil {
		if e != io.EOF {
			fmt.Fprintf(os.Stderr, "Problem reading a byte: %s\n", e)
		}
	}

	var sum uint64
	var sumOfSquares uint64

	for _, c := range vector {
		x := uint64(c)
		sum += x
		sumOfSquares += x * x

	}

	return vector, float64(sum), float64(sumOfSquares)
}
*/

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

/* Turn a string like "101/127/134/..." into an array
 * of bytes that have those numerical values.  */
func convertKeybytes(str string) []byte {

	var keyBytes []byte

	shiftsAsStrings := strings.Split(str, "/")

	for _, shft := range shiftsAsStrings {
		if n, e := strconv.Atoi(shft); e == nil {
			keyBytes = append(keyBytes, byte(n))
		} else {
			fmt.Fprintf(os.Stderr, "Problem with shift %q: %s\n", shft, e)
		}
	}

	return keyBytes
}

func vectorAngle(vector1, vector2 *Vector) float64 {
	var dotProduct float64

	if len(vector1.vector) != len(vector2.vector) {
		log.Fatalf("Vectors not of same dimension: %d != %d\n", len(vector1.vector), len(vector2.vector))
	}

	for i, v1 := range vector1.vector {
		dotProduct += float64(v1 * vector2.vector[i])
	}

	magA := math.Sqrt(vector1.sumOfSquares)
	magB := math.Sqrt(vector2.sumOfSquares)
	z := dotProduct / (magA * magB)

	// math.Acos() undefined for argument -1 <= x >= 1,
	// and we know that z is positive.
	return math.Acos(z - 0.0000001)
}

var phpVector = []int{
	56,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	621136,
	548071,
	0,
	0,
	127054,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	2,
	0,
	0,
	0,
	0,
	0,
	3123392,
	25305,
	244473,
	12657,
	401123,
	8595,
	19836,
	442267,
	316939,
	315698,
	58844,
	21610,
	193932,
	141891,
	250422,
	199635,
	103920,
	106886,
	90341,
	107117,
	63592,
	62273,
	75678,
	52547,
	40477,
	41334,
	55481,
	231735,
	128623,
	249523,
	221584,
	28568,
	36028,
	75354,
	28089,
	43516,
	27413,
	70139,
	32110,
	22891,
	24404,
	43505,
	9130,
	11073,
	40876,
	27692,
	40642,
	81754,
	88299,
	11611,
	54075,
	109753,
	89416,
	28132,
	33229,
	18803,
	11915,
	12899,
	14728,
	99586,
	205455,
	101878,
	941,
	226254,
	2669,
	598479,
	161159,
	423983,
	323876,
	1020634,
	270813,
	133997,
	281705,
	605646,
	17580,
	74677,
	432070,
	256552,
	500212,
	525263,
	390768,
	40961,
	607450,
	668141,
	821533,
	290712,
	89184,
	96185,
	209069,
	124786,
	22459,
	81949,
	6388,
	82738,
	1540,
	0,
	152,
	74,
	38,
	48,
	33,
	10,
	8,
	12,
	51,
	10,
	2,
	9,
	46,
	0,
	6,
	10,
	5,
	4,
	2,
	3,
	116,
	3,
	284,
	8,
	6,
	6,
	10,
	14,
	8,
	23,
	2,
	98,
	4,
	5,
	4,
	23,
	6,
	4,
	4,
	6,
	5,
	14,
	4,
	4,
	4,
	4,
	4,
	4,
	101,
	4,
	21,
	17,
	54,
	87,
	12,
	16,
	48,
	27,
	36,
	79,
	10,
	53,
	103,
	73,
	3,
	2,
	87,
	139,
	4,
	12,
	2,
	3,
	4,
	0,
	2,
	0,
	2,
	3,
	2,
	2,
	692,
	304,
	5,
	0,
	0,
	1,
	0,
	1,
	0,
	1,
	0,
	0,
	1,
	20,
	0,
	4,
	155,
	15,
	422,
	15,
	64,
	201,
	25,
	36,
	131,
	60,
	52,
	91,
	108,
	93,
	216,
	72,
	89,
	102,
	157,
	35,
	24,
	21,
	0,
	35,
	22,
	14,
	1,
	25,
	44,
	5,
	5,
	35,
}
var phpSum float64 = 19622992.0
var phpSumOfSquares float64 = 16516921766410.0
var englishVector = []int{
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	209,
	2989,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	17105,
	17,
	276,
	53,
	0,
	3,
	13,
	59,
	353,
	370,
	235,
	26,
	884,
	0,
	1561,
	793,
	310,
	228,
	187,
	92,
	61,
	51,
	66,
	59,
	82,
	73,
	384,
	10,
	75,
	402,
	86,
	7,
	8,
	348,
	177,
	344,
	246,
	323,
	97,
	296,
	86,
	480,
	85,
	12,
	175,
	143,
	238,
	207,
	292,
	4,
	286,
	393,
	425,
	149,
	54,
	77,
	88,
	55,
	6,
	307,
	9,
	309,
	1,
	304,
	0,
	5771,
	1861,
	3353,
	2887,
	9982,
	1769,
	2140,
	2663,
	6190,
	290,
	494,
	4378,
	2147,
	5410,
	6203,
	2120,
	78,
	5312,
	5498,
	7266,
	2649,
	999,
	928,
	292,
	1106,
	90,
	22,
	0,
	22,
	3,
	0,
	0,
	0,
	0,
	0,
	3,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	3,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	3,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
}
var englishSum float64 = 115075.0
var englishSumOfSquares float64 = 734987429.0
