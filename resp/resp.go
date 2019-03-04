package resp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Typ interface {
}

func Parse(r *bufio.Reader) (Typ, error) {
	c := make([]byte, 1)
	n, err := r.Read(c)
	if n != 1 || err != nil {
		return nil, err
	}
	switch c[0] {
	case '+':
		panic("simple string")
	case '-':
		panic("error")
	case ':':
		panic("integer")
	case '$':
		panic("bulk string")
	case '*':
		return ParseArray(r)
	}
	return nil, fmt.Errorf("Unknown RESP prefix: [%s]", c)
}

func readRespLine(r *bufio.Reader) (string, error) {
	l, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Failed to read line: %s", err)
	}
	l = strings.TrimRight(l, "\r\n")
	return l, nil
}

func ParseArray(r *bufio.Reader) (Typ, error) {
	l, err := readRespLine(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to read array size: %s", err)
	}
	size, err := strconv.Atoi(l)
	if err != nil {
		return nil, fmt.Errorf("Array size not numeric: [%s]", l)
	}
	a := make([]Typ, 0, size)
	for i := 0; i < size; i++ {
		a[i], err = Parse(r)
		if err != nil {
			return nil, fmt.Errorf("Error reading array elt [%d]: %s", i, err)
		}
	}
	return a, nil
}
