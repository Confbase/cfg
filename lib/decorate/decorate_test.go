package decorate

import "testing"

type input struct {
	s    string
	code string
}

type pair struct {
	expect string
	input
}

func TestDecorate(t *testing.T) {
	pairs := []pair{
		{
			input:  input{code: "0;31", s: "hello, world!"},
			expect: "\033[0;31mhello, world!\033[0m",
		},
		{
			input:  input{code: "0;31", s: "some escaped char\n"},
			expect: "\033[0;31msome escaped char\n\033[0m",
		},
	}
	for i, p := range pairs {
		got := Decorate(p.s, p.code)
		if got != p.expect {
			t.Logf("on test %v, expected %v but got %v", i, p.expect, got)
			t.Fail()
		}
	}

}

func TestTitle(t *testing.T) {
	pairs := []struct {
		input  string
		expect string
	}{
		{
			input:  "hello, world!",
			expect: "\033[4m\033[1mhello, world!\033[0m\033[0m",
		},
	}
	for i, p := range pairs {
		got := Title(p.input)
		if got != p.expect {
			t.Logf("on test %v, expected %v but got %v", i, p.expect, got)
			t.Fail()
		}
	}

}
