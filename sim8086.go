package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

var regormem = "100010"
var itoreg = "1011"
var opcodes = map[string]string{
	regormem: "mov",
	itoreg:   "mov",
}

var regs = map[string][]string{
	"000": {"AL", "AX"},
	"001": {"CL", "CX"},
	"010": {"DL", "DX"},
	"011": {"BL", "BX"},
	"100": {"AH", "SP"},
	"101": {"CH", "BP"},
	"110": {"DH", "SI"},
	"111": {"BH", "DI"},
}

func main() {
	argv := os.Args
	if len(argv) < 2 {
		fmt.Println("please specify input file")
		return
	}
	f, err := os.Open(argv[1])
	checkErr(err)
	defer f.Close()

	fmt.Printf("opening binary file %s\n", argv[1])
	fmt.Println("bits 16")

	// read in chunks
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var bytes []byte
	for scanner.Scan() {
		bytes = append(bytes, scanner.Bytes()...)
	}

	bit16 := make([][]string, 0)
	var bits []string
	i := 0
	for i < len(bytes) {
		bits = append(bits, asBits(bytes[i])+asBits(bytes[i+1]))
		if i%2 == 0 {
			bit16 = append(bit16, append([]string{}, bits...))
			bits = bits[:0]
		}
		i += 2
	}
	fmt.Println(bytes, bit16)

	// loop over bit16 to decode it to assembly instructions
	for j := 0; j < len(bit16); j++ {
		line := bit16[j][0]
		opcode := line[0:4]
		if opcode != itoreg {
			opcode = line[0:6]
		}

		var inst string
		switch opcode {
		case itoreg:
			w, err := strconv.Atoi(string(line[4]))
			if err != nil {
				fmt.Println(line, err)
			}
			reg := line[5:8]
			data := toDecimal(line[8:16])
			inst += fmt.Sprintf("%s %s %d", opcodes[opcode], regs[reg][w], data)
			fmt.Println(inst)

		case regormem:
			reg := line[10:13]
			rm := line[13:16]
			d := line[6]
			w, err := strconv.Atoi(string(line[7]))
			if err != nil {
				fmt.Println(line, err)
			}
			dest := regs[rm][w]
			source := regs[reg][w]
			if d == '0' {
				inst += fmt.Sprintf("%s %s, %s", opcodes[opcode], dest, source)
			} else {
				inst += fmt.Sprintf("%s %s, %s", opcodes[opcode], source, dest)
			}
			fmt.Println(inst)
		}
	}
}

func asBits(val uint8) string {
	bits := ""
	for i := 0; i < 8; i++ {
		bits = fmt.Sprintf("%x%s", val&0x1, bits)
		val = val >> 1
	}
	return bits
}

func toDecimal(bits string) int {
	var total int

	for i := 0; i < len(bits); i++ {
		if bits[i] == '1' {
			y := float64(len(bits) - 1 - i)
			total += int(math.Pow(2, y))
		}
	}

	return total
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}
