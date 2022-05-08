package assembler

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hlmerscher/hack-assembler-go/code"
	"github.com/hlmerscher/hack-assembler-go/parser"
)

func Assemble(input io.Reader) string {
	psr := parser.New(input)

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

		if instruction == parser.A_INSTRUCTION || instruction == parser.L_INSTRUCTION {
			symbol := parser.Symbol(line)
			parsedContent = code.Symbol(symbol)
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
