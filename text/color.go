package text

import "fmt"

func Green(s string) string {
	return fmt.Sprint("\033[32m", s, "\033[0m")
}

func Red(s string) string {
	return fmt.Sprint("\033[31m", s, "\033[0m")
}

func Default(s string) string {
	return fmt.Sprint("\033[39m", s, "\033[0m")
}
