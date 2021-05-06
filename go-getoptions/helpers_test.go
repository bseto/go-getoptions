package getoptions

import (
	"reflect"
	"testing"
)

func TestIsOption(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		mode    Mode
		optPair []optionPair
		is      bool
	}{
		{"lone dash", "-", Normal, []optionPair{{Option: "-"}}, true},
		{"lone dash", "-", Bundling, []optionPair{{Option: "-"}}, true},
		{"lone dash", "-", SingleDash, []optionPair{{Option: "-"}}, true},

		// TODO: Lets not return an option here
		// Lets let the caller identify this.
		{"double dash", "--", Normal, []optionPair{{Option: "--"}}, false},
		{"double dash", "--", Bundling, []optionPair{{Option: "--"}}, false},
		{"double dash", "--", SingleDash, []optionPair{{Option: "--"}}, false},

		{"no option", "opt", Normal, []optionPair{}, false},
		{"no option", "opt", Bundling, []optionPair{}, false},
		{"no option", "opt", SingleDash, []optionPair{}, false},

		{"Long option", "--opt", Normal, []optionPair{{Option: "opt"}}, true},
		{"Long option", "--opt", Bundling, []optionPair{{Option: "opt"}}, true},
		{"Long option", "--opt", SingleDash, []optionPair{{Option: "opt"}}, true},

		{"Long option with arg", "--opt=arg", Normal, []optionPair{{Option: "opt", Args: []string{"arg"}}}, true},
		{"Long option with arg", "--opt=arg", Bundling, []optionPair{{Option: "opt", Args: []string{"arg"}}}, true},
		{"Long option with arg", "--opt=arg", SingleDash, []optionPair{{Option: "opt", Args: []string{"arg"}}}, true},

		{"short option", "-opt", Normal, []optionPair{{Option: "opt"}}, true},
		{"short option", "-opt", Bundling, []optionPair{{Option: "o"}, {Option: "p"}, {Option: "t"}}, true},
		{"short option", "-opt", SingleDash, []optionPair{{Option: "o", Args: []string{"pt"}}}, true},

		{"short option with arg", "-opt=arg", Normal, []optionPair{{Option: "opt", Args: []string{"arg"}}}, true},
		{"short option with arg", "-opt=arg", Bundling, []optionPair{{Option: "o"}, {Option: "p"}, {Option: "t", Args: []string{"arg"}}}, true},
		{"short option with arg", "-opt=arg", SingleDash, []optionPair{{Option: "o", Args: []string{"pt=arg"}}}, true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			buf := setupLogging()
			optPair, is := isOption(tt.in, tt.mode)
			if !reflect.DeepEqual(optPair, tt.optPair) || is != tt.is {
				t.Errorf("isOption(%q, %q) == (%q, %v), want (%q, %v)",
					tt.in, tt.mode, optPair, is, tt.optPair, tt.is)
			}
			t.Log(buf.String())
		})
	}
}
