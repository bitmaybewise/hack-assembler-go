package assembler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hlmerscher/hack-assembler-go/code"
	"github.com/hlmerscher/hack-assembler-go/parser"
)

func Assemble(input *os.File) string {
	psr := parser.New(input)
	symbols, parsedLines := firstPass(psr)
	binary := secondPass(symbols, parsedLines)
	return binary
}

// firstPass builds the symbol table and returns the parsed lines to avoid going over it again
func firstPass(psr parser.Parser) (symbols map[string]int, lines []parser.ParsedLine) {
	symbols = make(map[string]int)
	for k, v := range code.Symbols {
		symbols[k] = v
	}

	lines = make([]parser.ParsedLine, 0)
	var i int
	for {
		line, err := psr.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if errors.Is(err, parser.IgnoredLine) {
			continue
		}

		lines = append(lines, line)

		if line.IsL() {
			symbols[line.Symbol()] = i
		} else {
			i++
		}
	}

	return
}

// secondPass goes over the symbol table to look up for the symbol value and transforms the content into binary
func secondPass(symbols map[string]int, parsedLines []parser.ParsedLine) string {
	variables := 16 // reserved memory space starts at address 16
	var binary strings.Builder

	for _, line := range parsedLines {
		var parsedContent string

		if line.IsL() {
			continue
		}
		if line.IsA() {
			n, err := strconv.Atoi(line.Symbol())
			if err == nil {
				parsedContent = fmt.Sprintf("%.16b", n)
			} else {
				val, ok := symbols[line.Symbol()]
				if ok {
					parsedContent = fmt.Sprintf("%.16b", val)
				} else {
					symbols[line.Symbol()] = variables
					parsedContent = fmt.Sprintf("%.16b", variables)
					variables++
				}
			}
		}
		if line.IsC() {
			binDest := code.Dest(line.Dest())
			binComp := code.Comp(line.Comp())
			binJump := code.Jump(line.Jump())
			parsedContent = fmt.Sprintf("111%s%s%s", binComp, binDest, binJump)
		}

		binary.WriteString(parsedContent)
		binary.WriteByte('\n')
	}

	return binary.String()
}
