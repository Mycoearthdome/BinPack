package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

func RebuildMap(sOdd string) map[int64]bool {
	var odd = make(map[int64]bool)
	var value string
	var ValueSubString string
	var result float64
	//var bits []string
	//fmt.Println("Buffer -->", sOdd)
	//fmt.Println(len(sOdd))

	for i := 0; i < len(sOdd); i = i + 64 {
		result = 0.0
		if i+64 <= len(sOdd) {
			ValueSubString = sOdd[i : i+64]
			//fmt.Println(ValueSubString)
			for j := len(ValueSubString); j > 0; j-- {
				value = ValueSubString[j-1 : j]
				//fmt.Print(value)
				if value == "1" {
					result += math.Pow(2.0, float64(len(ValueSubString)-j))
				}
			}
			//fmt.Println("")
			odd[int64(result)] = true
		}

	}
	return odd
}

func bytes2binary(bytes []byte, Encode bool) string {
	var Count_ONE int64
	var bits string
	var odd = make(map[int64]bool)
	var j int64 = 0
	var sOdd string

	for _, byte := range bytes {
		for i := 1; i <= 8; i++ {
			if (byte>>i)&1 == 1 {
				bits += "1"
				Count_ONE++
			} else {
				bits += "0"
			}
			if byte%2 != 0 { //ODD
				odd[j] = true
			}
		}
		j++
	}

	for j = ((j * 8) / 2); j%8 != 0; j++ { //the where!
		continue
	}

	if Encode {
		before := bits[:j]
		after := bits[j:]

		for key, _ := range odd {
			//fmt.Println(key, "--->", strconv.FormatInt(key, 2))
			sOdd += fmt.Sprintf("%032s%032s", "", strconv.FormatInt(key, 2)) // max filesize is 4.295 GB!!!!!
		}

		//fmt.Println("Buffer-->", sOdd)

		before += "1111010101010100" + sOdd + "1010101010111110" + after //encapsulated

		bits = before
	}
	//fmt.Println(bytes)
	//fmt.Println(bits)
	//fmt.Println(len(bits))
	//fmt.Println(Count_ONE)
	//fmt.Println(odd)
	return bits
}

func binary2bytes(binary string) []byte {
	var Count_ONE int64
	var bytes []byte
	var pos int64 = 0
	var sByte []string
	var i int64
	var j int64 = 0
	var dec float64
	var before string
	var after string
	var sOdd string
	var odd = make(map[int64]bool)

	before = strings.Split(binary, "1111010101010100")[0]
	after = strings.Split(binary, "1010101010111110")[1]
	sOdd = strings.Split(binary, before)[1]
	sOdd = strings.Split(sOdd, after)[0]
	//fmt.Println(sOdd)
	//fmt.Println("_______________________________________")
	sOdd = strings.Split(sOdd, "1111010101010100")[1]
	sOdd = strings.Split(sOdd, "1010101010111110")[0]
	//fmt.Println(sOdd)

	binary = before + after
	//fmt.Println("Buffer-->", sOdd)

	odd = RebuildMap(sOdd)

	for i = 8; i <= int64(len(binary)); i = i + 8 {
		sByte = append(sByte, binary[pos:i])
		pos = i
	}

	for _, bString := range sByte {
		//fmt.Println(bString)
		dec = 0
		for i := 7; i > 0; i-- {
			if strings.Compare(string(bString[i-1:i]), "1") == 0 {
				dec += math.Pow(2.0, float64(i+1))
				Count_ONE++
			}
		}
		if odd[j] {
			dec += 2.0
		}
		bytes = append(bytes, byte(dec/2)) // / 2 twist
		j++
	}
	//bytes = append(bytes, byte(10))
	//fmt.Println(bytes)
	//fmt.Println(Count_ONE)
	return bytes
}

func LoadRAM(file *os.File, stdin bool) []byte {
	var bytes []byte
	var err error

	if stdin {
		bytes, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

	} else {
		bytes, err = io.ReadAll(file)
		if err != nil {
			panic(err)
		}
	}
	return bytes
}

func saveFile(filename string, bin string) {
	var pos int64 = 0
	var value int
	var result float64
	var bytes []byte
	var i, j int64
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	for i = 7; i < int64(len(bin)); i = i + 8 {
		result = 0
		for j = 0; j < 8; j++ {
			value, err = strconv.Atoi(bin[pos+j : pos+j+1])
			if err != nil {
				panic(err)
			}
			if value == 1 {
				result += math.Pow(2.0, float64(7-j))
			}
		}
		//fmt.Println(result)
		//os.Exit(0)
		bytes = append(bytes, byte(uint8(result)))
		pos = i + 1
	}
	//os.Exit(0)
	//fmt.Println(bytes)
	file.Write(bytes)
	//err = binary.Write(file, binary.LittleEndian, &bytes)
	//if err != nil {
	//	panic(err)
	//}
	file.Close()
}

func readFile(filename string) string {
	var bin string
	//fileStats, err := os.Stat(filename)
	//if err != nil {
	//	panic(err)
	//}
	//data := make([]byte, fileStats.Size())
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//err = binary.Read(file, binary.LittleEndian, &data)
	//if err != nil {
	//	panic(err)
	//}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	for _, byte := range data {
		bin += fmt.Sprintf("%08s", strconv.FormatInt(int64(byte), 2))
		//fmt.Println(bin)
		//os.Exit(0)
	}

	//return bytes2binary(data, false)
	return bin
}

func main() {
	var data []byte
	var outfile string = "out.file"
	var Filename string
	var Encoded string
	var Decoded string
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		//data is being piped to stdin"
		data = LoadRAM(nil, true)
	} else {
		//is from a terminal
		Arguments := os.Args[1:]
		if len(Arguments) > 0 {
			Filename = Arguments[0]
		}

		file, err := os.OpenFile(Filename, os.O_RDONLY, 0)
		if errors.Is(err, os.ErrNotExist) {
			// handle the case where the file doesn't exist
			fmt.Println("Pipe me a file OR use the filename as an argument without a switch please.")
			os.Exit(0)
		}
		data = LoadRAM(file, false)
		file.Close()
	}
	Encoded = bytes2binary(data, true)
	//fmt.Println(Encoded)
	//fmt.Println(len(Encoded))
	//fmt.Println(string(binary2bytes(Encoded)))
	saveFile(outfile, Encoded)
	fmt.Println("----------------------")
	Decoded = readFile(outfile)
	//fmt.Println(string(binary2bytes((Encoded))))
	fmt.Println(Decoded)
	fmt.Println(string(binary2bytes((Decoded))))
	//fmt.Println(len(readFile(outfile)))
}
