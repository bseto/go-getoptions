package getoptions

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestIsOption(t *testing.T) {
	cases := []struct {
		name     string
		inputArg string
		mode     Mode
		expected []*programTree
		isOption bool
	}{
		{"lone dash", "-", Normal, []*programTree{newCLIArg(nil, argTypeOption, "-")}, true},
		{"lone dash", "-", Bundling, []*programTree{newCLIArg(nil, argTypeOption, "-")}, true},
		{"lone dash", "-", SingleDash, []*programTree{newCLIArg(nil, argTypeOption, "-")}, true},

		// TODO: Lets not return an option here
		// Lets let the caller identify this.
		{"double dash", "--", Normal, []*programTree{newCLIArg(nil, argTypeTerminator, "--")}, false},
		{"double dash", "--", Bundling, []*programTree{newCLIArg(nil, argTypeTerminator, "--")}, false},
		{"double dash", "--", SingleDash, []*programTree{newCLIArg(nil, argTypeTerminator, "--")}, false},

		{"no option", "opt", Normal, []*programTree{}, false},
		{"no option", "opt", Bundling, []*programTree{}, false},
		{"no option", "opt", SingleDash, []*programTree{}, false},

		{"Long option", "--opt", Normal, []*programTree{newCLIArg(nil, argTypeOption, "opt")}, true},
		{"Long option", "--opt", Bundling, []*programTree{newCLIArg(nil, argTypeOption, "opt")}, true},
		{"Long option", "--opt", SingleDash, []*programTree{newCLIArg(nil, argTypeOption, "opt")}, true},

		{"Long option with arg", "--opt=arg", Normal, []*programTree{newCLIArg(nil, argTypeOption, "opt", "arg")}, true},
		{"Long option with arg", "--opt=arg", Bundling, []*programTree{newCLIArg(nil, argTypeOption, "opt", "arg")}, true},
		{"Long option with arg", "--opt=arg", SingleDash, []*programTree{newCLIArg(nil, argTypeOption, "opt", "arg")}, true},

		{"short option", "-opt", Normal, []*programTree{newCLIArg(nil, argTypeOption, "opt")}, true},
		{"short option", "-opt", Bundling, []*programTree{newCLIArg(nil, argTypeOption, "o"), newCLIArg(nil, argTypeOption, "p"), newCLIArg(nil, argTypeOption, "t")}, true},
		{"short option", "-opt", SingleDash, []*programTree{newCLIArg(nil, argTypeOption, "o", "pt")}, true},

		{"short option with arg", "-opt=arg", Normal, []*programTree{newCLIArg(nil, argTypeOption, "opt", "arg")}, true},
		{"short option with arg", "-opt=arg", Bundling, []*programTree{newCLIArg(nil, argTypeOption, "o"), newCLIArg(nil, argTypeOption, "p"), newCLIArg(nil, argTypeOption, "t", "arg")}, true},
		{"short option with arg", "-opt=arg", SingleDash, []*programTree{newCLIArg(nil, argTypeOption, "o", "pt=arg")}, true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			buf := setupLogging()
			output, is := isOption(tt.inputArg, tt.mode)
			if !reflect.DeepEqual(output, tt.expected) || is != tt.isOption {
				t.Errorf("input: %s, mode: %d\nexpected (%v) tree: %s\n got: (%v) %s\n",
					tt.inputArg, tt.mode, tt.isOption, spew.Sdump(tt.expected), is, spew.Sdump(output))
			}
			t.Log(buf.String())
		})
	}
}
