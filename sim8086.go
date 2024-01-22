package main

import (
	"bufio"
	"fmt"
	"os"
)

var opcodes = map[string]string{
	"100010": "mov",
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

	bit16 := make([][]uint8, 0)
	var bits []uint8
	for i := 0; i < len(bytes); i++ {
		bits = append(bits, asBits(bytes[i])...)
		if i%2 != 0 {
			bit16 = append(bit16, append([]uint8{}, bits...))
			bits = bits[:0]
		}
	}

	// loop over bit16 to decode it to assembly instructions
	for j := 0; j < len(bit16); j++ {
		line := bit16[j]
		var inst, opcode, reg, rm string
		var d, w int
		for k := 0; k < len(line); k++ {
			if len(opcode) < 6 {
				if line[k] == 1 {
					opcode += "1"
				} else {
					opcode += "0"
				}
			} else if k == 6 {
				d = int(line[k])
			} else if k == 7 {
				w = int(line[k])
			} else if k == 8 || k == 9 {
				//do MOD later
			} else if k > 9 && k <= 12 {
				if line[k] == 1 {
					reg += "1"
				} else {
					reg += "0"
				}
			} else if k > 12 && k <= 15 {
				if line[k] == 1 {
					rm += "1"
				} else {
					rm += "0"
				}
			}
		}
		dest := regs[rm][w]
		source := regs[reg][w]
		if d == 0 {
			inst += fmt.Sprintf("%s %s, %s", opcodes[opcode], dest, source)
		} else {
			inst += fmt.Sprintf("%s %s, %s", opcodes[opcode], source, dest)
		}
		fmt.Println(inst)
	}
}

func asBits(val uint8) []uint8 {
	bits := []uint8{}
	for i := 0; i < 8; i++ {
		bits = append([]uint8{val & 0x1}, bits...)
		val = val >> 1
	}
	return bits
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}
