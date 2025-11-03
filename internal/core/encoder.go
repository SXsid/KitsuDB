package core

import "fmt"

func Encode(value any, isSimple bool) (b []byte) {
	switch v := value.(type) {
	case string:
		if isSimple {
			return fmt.Appendf(b, "+%s\r\n", v)
		} else {
			return fmt.Appendf(b, "$%d\r\n%s\r\n", len(v), v)
		}
	}
	return b
}
