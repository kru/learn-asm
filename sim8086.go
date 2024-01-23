package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	// loop over bit16 to decode it to assembly instructions
	for j := 0; j < len(bit16); j++ {
		line := bit16[j][0]
		opcode := line[0:6]
		reg := line[10:13]
		rm := line[13:16]
		d := line[6]
		var inst string
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

func asBits(val uint8) string {
	bits := ""
	for i := 0; i < 8; i++ {
		bits = fmt.Sprintf("%x%s", val&0x1, bits)
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
