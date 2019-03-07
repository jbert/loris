package resp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/jbert/goredis/store"
)

var TheNullBulkString []byte = nil

type Typ interface{}

func TypToBuf(t Typ) []byte {
	switch v := t.(type) {
	case [][]byte:
		buf := []byte{}
		buf = append(buf, '*')
		buf = append(buf, fmt.Sprintf("%d\r\n", len(v))...)
		for _, rsp := range v {
			buf = append(buf, TypToBuf(rsp)...)
		}
		return buf
	case string:
		return []byte(fmt.Sprintf("+%s\r\n", v))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	case store.Val:
		return TypToBuf([]byte(v))
	case []byte:
		// 20 is bigger than log10(len(buf)) + 3
		// even if it isn't, it just means a realloc
		buf := make([]byte, 0, len(v)+20)

		l := len(v)
		// TheNullBulkString
		if v == nil {
			l = -1
		}
		buf = append(buf, '$')
		buf = append(buf, fmt.Sprintf("%d\r\n", l)...)
		if v != nil {
			buf = append(buf, v...)
			buf = append(buf, "\r\n"...)
		}
		return buf
	case int:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	default:
		panic(fmt.Sprintf("TODO - impl type %T", t))
	}
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
		return ParseBulkString(r)
	case '*':
		return ParseArray(r)
	}
	return nil, fmt.Errorf("Unknown RESP prefix: [%02x] [%s]", c, c)
}

func readNumLine(parseContext string, r *bufio.Reader) (int, error) {
	l, err := readRespLine(r)
	if err != nil {
		return 0, fmt.Errorf("Failed to read %s: %s", parseContext, err)
	}
	size, err := strconv.Atoi(l)
	if err != nil {
		return 0, fmt.Errorf("%s not numeric: [%s]", parseContext, l)
	}
	return size, nil
}

func readRespLine(r *bufio.Reader) (string, error) {
	l, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Failed to read line: %s", err)
	}
	l = strings.TrimRight(l, "\r\n")
	return l, nil
}

func ParseBulkString(r *bufio.Reader) (Typ, error) {
	strlen, err := readNumLine("bulk string length", r)
	if err != nil {
		return nil, err
	}
	s := make([]byte, strlen, strlen)
	num_read := 0
	for num_read < strlen {
		n, err := r.Read(s[num_read:strlen])
		if err != nil {
			return nil, fmt.Errorf("Read [%d/%d] of bulk string: %s", num_read, strlen, err)
		}
		num_read += n
	}
	l, err := readRespLine(r)
	if l != "" || err != nil {
		return nil, fmt.Errorf("Failed to find CRLF at end of bulk string: %s", err)
	}
	return s, nil
}

func ParseArray(r *bufio.Reader) (Typ, error) {
	size, err := readNumLine("array size", r)
	if err != nil {
		return nil, err
	}
	a := make([]Typ, size)
	for i := 0; i < size; i++ {
		a[i], err = Parse(r)
		if err != nil {
			return nil, fmt.Errorf("Error reading array elt [%d]: %s", i, err)
		}
	}
	return a, nil
}
