package decorate

import "fmt"

type Decorator struct {
	Enabled bool
}

func New() *Decorator {
	return &Decorator{Enabled: true}
}

func (d *Decorator) Decorate(s, code string) string {
	if d.Enabled {
		return Decorate(s, code)
	}
	return s
}

func (d *Decorator) Title(s string) string {
	return d.Decorate(d.Decorate(s, "1"), "4")
}

func (d *Decorator) Red(s string) string {
	return d.Decorate(s, "0;31")
}

func (d *Decorator) Green(s string) string {
	return d.Decorate(s, "0;32")
}

func (d *Decorator) LightBlue(s string) string {
	return d.Decorate(s, "1;34")
}

func Decorate(s, code string) string {
	return fmt.Sprintf("\033[%vm%v\033[0m", code, s)
}

func Title(s string) string {
	return Decorate(Decorate(s, "1"), "4")
}

func Red(s string) string {
	return Decorate(s, "0;31")
}

func Green(s string) string {
	return Decorate(s, "0;32")
}

func LightBlue(s string) string {
	return Decorate(s, "1;34")
}
