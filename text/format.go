package text

import "fmt"

func Bold(s string) string {
	return fmt.Sprintf("%s%s%s", "\033[1m", s, "\033[0m")
}
