package filetype

import "testing"

func TestGuess(t *testing.T) {
	pairs := []struct {
		input  string
		expect string
	}{
		{input: "file.json", expect: Json},
		{input: "file.JSON", expect: Json},
		{input: "some_other_file.json", expect: Json},
		{input: "So99082me_other_file.JsOn", expect: Json},
		{input: "So99082me_other_file.jSoN", expect: Json},
		{input: "So99082me_other_file.jsOn", expect: Json},
		{input: "some/path/to/some/file.jsOn", expect: Json},
		{input: "some/path/to/some/file.json", expect: Json},
		{input: "some/path/to/some/file.JSON", expect: Json},
		{input: "some\\path\\to\\some\\file.json", expect: Json},
		{input: "file.yaml", expect: Yaml},
		{input: "file.YAML", expect: Yaml},
		{input: "some_other_file.yaml", expect: Yaml},
		{input: "So99082me_other_file.Yaml", expect: Yaml},
		{input: "So99082me_other_file.yaml", expect: Yaml},
		{input: "So99082me_other_file.yaml", expect: Yaml},
		{input: "some/path/to/some/file.yaml", expect: Yaml},
		{input: "some/path/to/some/file.yaml", expect: Yaml},
		{input: "some/path/to/some/file.YAML", expect: Yaml},
		{input: "some\\path\\to\\some\\file.yaml", expect: Yaml},
		{input: "file.toml", expect: Toml},
		{input: "file.TOML", expect: Toml},
		{input: "some_other_file.toml", expect: Toml},
		{input: "So99082me_other_file.Toml", expect: Toml},
		{input: "So99082me_other_file.toml", expect: Toml},
		{input: "So99082me_other_file.toml", expect: Toml},
		{input: "some/path/to/some/file.toml", expect: Toml},
		{input: "some/path/to/some/file.toml", expect: Toml},
		{input: "some/path/to/some/file.TOML", expect: Toml},
		{input: "some\\path\\to\\some\\file.toml", expect: Toml},
		{input: "file.exe", expect: Unknown},
		{input: "file.jar", expect: Unknown},
		{input: "some_other_file.dll", expect: Unknown},
		{input: "So99082me_other_file.H", expect: Unknown},
		{input: "So99082me_other_file.go", expect: Unknown},
		{input: "So99082me_other_file.rs", expect: Unknown},
		{input: "some/path/to/some/file.hs", expect: Unknown},
		{input: "some/path/to/some/file.gov", expect: Unknown},
		{input: "some/path/to/some/file.fbi", expect: Unknown},
		{input: "some\\path\\to\\some\\file.edu", expect: Unknown},
	}
	for i, p := range pairs {
		got := Guess(p.input)
		if got != p.expect {
			t.Logf("on test %v, expected %v but got %v", i, p.expect, got)
			t.Fail()
		}
	}
}
