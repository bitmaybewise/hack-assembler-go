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

	// first pass
	symbols := make(map[string]int)
	for k, v := range code.Symbols {
		symbols[k] = v
	}

	var i int
	for {
		line, err := psr.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if errors.Is(err, parser.IgnoredLine) {
			continue
		}
		instruction := parser.InstructionType(line)
		if instruction == parser.A_INSTRUCTION || instruction == parser.C_INSTRUCTION {
			i++
		}
		if instruction == parser.L_INSTRUCTION {
			symbol := parser.Symbol(line)
			symbols[symbol] = i
		}
	}

	// second pass
	_, err := input.Seek(0, io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("rewind file error <%s>", err))
	}
	psr = parser.New(input)

	variables := 16
	var b strings.Builder
	for {
		line, err := psr.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if errors.Is(err, parser.IgnoredLine) {
			continue
		}

		instruction := parser.InstructionType(line)
		var parsedContent string

		if instruction == parser.L_INSTRUCTION {
			continue
		}
		if instruction == parser.A_INSTRUCTION {
			symbol := parser.Symbol(line)
			n, err := strconv.Atoi(symbol)
			if err == nil {
				parsedContent = fmt.Sprintf("%.16b", n)
			} else {
				val, ok := symbols[symbol]
				if ok {
					parsedContent = fmt.Sprintf("%.16b", val)
				} else {
					symbols[symbol] = variables
					parsedContent = fmt.Sprintf("%.16b", variables)
					variables++
				}
			}
		}
		if instruction == parser.C_INSTRUCTION {
			dest := parser.Dest(line)
			binDest := code.Dest(dest)

			comp := parser.Comp(line)
			binComp := code.Comp(comp)

			jump := parser.Jump(line)
			binJump := code.Jump(jump)

			parsedContent = fmt.Sprintf("111%s%s%s", binComp, binDest, binJump)
		}

		b.WriteString(parsedContent)
		b.WriteByte('\n')
	}

	return b.String()
}
