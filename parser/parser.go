package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Instruction int

const (
	A_INSTRUCTION = iota
	C_INSTRUCTION
	L_INSTRUCTION
)

var IgnoredLine = errors.New("ignored line")

type ParsedLine struct {
	value string
}

func (pl ParsedLine) IsL() bool {
	return InstructionType(pl.value) == L_INSTRUCTION
}

func (pl ParsedLine) IsA() bool {
	return InstructionType(pl.value) == A_INSTRUCTION
}

func (pl ParsedLine) IsC() bool {
	return InstructionType(pl.value) == C_INSTRUCTION
}

func (pl ParsedLine) Symbol() string {
	if strings.HasPrefix(pl.value, "@") {
		return pl.value[1:]
	}
	pl.value = strings.Replace(pl.value, "(", "", 1)
	pl.value = strings.Replace(pl.value, ")", "", 1)
	return pl.value
}

func (pl ParsedLine) Dest() string {
	command := strings.Split(pl.value, "=")
	if len(command) == 1 {
		return ""
	}
	return command[0]
}

func (pl ParsedLine) Comp() (comp string) {
	command := strings.Split(pl.value, "=")
	if len(command) > 1 {
		comp = command[1]
	} else {
		comp = command[0]
	}
	command = strings.Split(comp, ";")
	if len(command) > 1 {
		comp = command[0]
	}
	return
}

func (pl ParsedLine) Jump() string {
	command := strings.Split(pl.value, ";")
	if len(command) == 1 {
		return ""
	}
	return command[1]
}

var EmptyLine = ParsedLine{}

type Parser struct {
	input *bufio.Reader
}

func (p *Parser) ReadLine() (ParsedLine, error) {
	line, err := p.input.ReadString('\n')
	if err != nil {
		return EmptyLine, err
	}

	// removing comment
	if strings.HasPrefix(line, "//") {
		return EmptyLine, IgnoredLine
	}
	commentFoundAt := strings.Index(line, "//")
	if commentFoundAt > 1 {
		line = line[:commentFoundAt-1]
	}

	line = strings.Replace(line, "\r", "", 1)
	line = strings.Replace(line, "\n", "", 1)
	line = strings.Trim(line, " ")
	if line == "" {
		return EmptyLine, IgnoredLine
	}

	return ParsedLine{line}, nil
}

func New(input io.Reader) Parser {
	reader := bufio.NewReader(input)
	return Parser{reader}
}

func InstructionType(line string) Instruction {
	if strings.HasPrefix(line, "@") {
		return A_INSTRUCTION
	}

	dest := strings.Split(line, "=")
	jump := strings.Split(line, ";")
	if len(dest) >= 2 || len(jump) >= 2 {
		return C_INSTRUCTION
	}

	return L_INSTRUCTION
}
