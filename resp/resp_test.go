package resp

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	a := assert.New(t)
	type testCase struct {
		s string
		t Typ
	}
	testCases := []testCase{
		{"$3\r\nfoo\r\n", []byte{'f', 'o', 'o'}},
	}

	for _, tc := range testCases {
		r := strings.NewReader(tc.s)
		t, err := Parse(bufio.NewReader(r))
		a.NoError(err, "Can parse [%s]", tc.s)
		a.Equal(tc.t, t)

		buf := TypToBuf(t)
		a.Equal(tc.s, string(buf))
	}
}
