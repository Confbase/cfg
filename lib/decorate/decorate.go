package decorate

import "fmt"

func Decorate(s, code string) string {
	return fmt.Sprintf("\033[%vm%v\033[0m", code, s)
}

func Title(s string) string {
	return Decorate(Decorate(s, "1"), "4")
}

func Red(s string) string {
	return Decorate(s, "0;31")
}

func Black(s string) string {
        return Decorate(s, "0;30")
}

func DarkGrey(s string) string {
	return Decorate(s, "1;30")
}

func LightRed(s string) string {
	return Decorate(s, "1;31")
}

func Green(s string) string {
        return Decorate(s, "0;32")
}

func LightGreen(s string) string {
	return Decorate(s, "1;32")
}

func Orange(s string) string {
	return Decorate(s, "0;33")
}

func Yellow(s string) string {
        return Decorate(s, "1;33")
}

func Blue(s string) string {
	return Decorate(s, "0;34")
}

func LightBlue(s string) string {
	return Decorate(s, "1;34")
}

func Purple(s string) string {
	return Decorate(s, "0;35")
}

func LightPurple(s string) string {
	return Decorate(s, "1;35")
}

func Cyan(s string) string {
	return Decorate(s, "0;36")
}

func LightCyan(s string) string {
	return Decorate(s, "1;36")
}

func LightGray(s string) string {
	return Decorate(s, "0;37")
}
