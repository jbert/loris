package loris

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkParseCommand(b *testing.B) {
	commands := make([]io.Reader, b.N)
	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i)
		c := fmt.Sprintf("*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$%d\r\n%s\r\n", len(s), s)
		commands[i] = strings.NewReader(c)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseCommand(commands[i])
	}
}
