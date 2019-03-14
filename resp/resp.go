package resp

import (
	"bufio"
	"fmt"

	"github.com/jbert/loris/store"
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

func Parse(r *bufio.Scanner) (Typ, error) {
	if !r.Scan() {
		return nil, r.Err()
	}
	c := r.Bytes()
	switch c[0] {
	case '+':
		panic("simple string")
	case '-':
		panic("error")
	case ':':
		panic("integer")
	case '$':
		return ParseBulkString(r, c[1:])
	case '*':
		return ParseArray(r, c[1:])
	}
	return nil, fmt.Errorf("Unknown RESP prefix: [%02x] [%s]", c, c)
}

func parseNumLine(c []byte) int {
	var n int
	for i := 0; i < len(c); i++ {
		n *= 10
		n += int(rune(c[i]) - '0')
	}
	return n
}

func ParseBulkString(r *bufio.Scanner, c []byte) (Typ, error) {
	strlen := parseNumLine(c)
	if !r.Scan() {
		return nil, r.Err()
	}

	b := r.Bytes()

	if len(b) != strlen {
		return nil, fmt.Errorf("Length mismatch %d != %d", len(b), strlen)
	}

	return b, nil
}

func ParseArray(r *bufio.Scanner, c []byte) (Typ, error) {
	size := parseNumLine(c)
	
	a := make([]Typ, size)
	for i := 0; i < size; i++ {
		t, err := Parse(r)
		if err != nil {
			return nil, fmt.Errorf("Error reading array elt [%d]: %s", i, err)
		}
		a[i] = t
	}
	return a, nil
}
