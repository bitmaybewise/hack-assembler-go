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

type Parser struct {
	input *bufio.Reader
}

func (p *Parser) ReadLine() (string, error) {
	line, err := p.input.ReadString('\n')
	if err != nil {
		return "", err
	}

	// removing comment
	if strings.HasPrefix(line, "//") {
		return "", IgnoredLine
	}
	commentFoundAt := strings.Index(line, "//")
	if commentFoundAt > 1 {
		line = line[:commentFoundAt-1]
	}

	line = strings.Replace(line, "\r", "", 1)
	line = strings.Replace(line, "\n", "", 1)
	line = strings.Trim(line, "")
	if line == "" {
		return "", IgnoredLine
	}

	return line, nil
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

func Symbol(line string) string {
	if strings.HasPrefix(line, "@") {
		return line[1:]
	}
	return line
}

func Dest(line string) string {
	command := strings.Split(line, "=")
	if len(command) == 1 {
		return ""
	}
	return command[0]
}

func Comp(line string) (comp string) {
	command := strings.Split(line, "=")
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

func Jump(line string) string {
	command := strings.Split(line, ";")
	if len(command) == 1 {
		return ""
	}
	return command[1]
}
