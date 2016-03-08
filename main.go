package main

import (
    "fmt"
    "io"
    "os"
    "debug/elf"
	"bufio"
	"strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ioReader(file string) io.ReaderAt {
    r, err := os.Open(file)
    check(err)
    return r
}

type Stack struct {
	Fuction uint32
	Pc      uint32
	Stack   uint32
}

func readLine(path string) []Stack {
	stacks := make([]Stack, 0)
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 
	
	for scanner.Scan() {
		text := scanner.Text()
		if text == "Function  PC        Stack" {
			fmt.Println(text)
			for scanner.Scan() {
				var stack Stack
				n, err := fmt.Fscanf(strings.NewReader(scanner.Text()),"%X %X %X", &stack.Fuction, &stack.Pc, &stack.Stack)
				if err != nil {
					break;
				}
				fmt.Printf("n = %d, fucntion = 0x%X, PC = 0x%X, Stack = 0x%X\n", n, stack.Fuction, stack.Pc, stack.Stack)
				stacks = append(stacks, stack)
			}
			return stacks
		}
	}
	
	return stacks
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: elftest elf_file log_file")
        os.Exit(1)
    }
	
	stacks := readLine(os.Args[2])
	//fmt.Println(stacks)
    f := ioReader(os.Args[1])
    _elf, err := elf.NewFile(f)
    check(err)

    // Read and decode ELF identifier
    var ident [16]uint8
    f.ReadAt(ident[0:], 0)
    check(err)

    if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
        fmt.Printf("Bad magic number at %d\n", ident[0:4])
        os.Exit(1)
    }

    fmt.Printf("File Header: ")
    fmt.Println(_elf.FileHeader)
    fmt.Printf("ELF Class: %s\n", _elf.Class.String())
    fmt.Printf("Machine: %s\n",  _elf.Machine.String())
    fmt.Printf("ELF Type: %s\n", _elf.Type)
    fmt.Printf("ELF Data: %s\n", _elf.Data)
    fmt.Printf("Entry Point: %d\n", _elf.Entry)

	for _, stack := range stacks {
		symbols, _ := _elf.Symbols()
		for _, symbol := range(symbols) {
			if symbol.Value == uint64(stack.Fuction) {
				fmt.Printf("Function name = %s\n", symbol.Name)
			}
		}
	}
	
}