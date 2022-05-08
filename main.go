package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hlmerscher/hack-assembler-go/assembler"
)

func main() {
	var filename string
	flag.StringVar(&filename, "f", "", "the filename of the asm source file")
	flag.Parse()
	if filename == "" {
		panic("filename is missing")
	}
	fmt.Printf("input:\t%s\n", filename)

	asmFile := openAsmFile(filename)
	defer asmFile.Close()
	binary := assembler.Assemble(asmFile)
	writeToHackFile(filename, binary)
}

func openAsmFile(filename string) *os.File {
	inputFile, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("error opening file <%s>\n", err))
	}
	return inputFile
}

func writeToHackFile(filename string, content string) {
	outputFilename := strings.Replace(filename, ".asm", ".hack", 1)
	fmt.Printf("output:\t%s\n", outputFilename)

	err := os.WriteFile(outputFilename, []byte(content), 0666)
	if err != nil {
		panic(err)
	}
}
