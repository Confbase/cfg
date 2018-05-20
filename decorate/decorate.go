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
		return fmt.Sprintf("\033[%vm%v\033[0m", code, s)
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
