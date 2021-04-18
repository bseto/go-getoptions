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
		expected []*ProgramTree
		isOption bool
	}{
		{"lone dash", "-", Normal, []*ProgramTree{NewCLIArg(argTypeOption, "-")}, true},
		{"lone dash", "-", Bundling, []*ProgramTree{NewCLIArg(argTypeOption, "-")}, true},
		{"lone dash", "-", SingleDash, []*ProgramTree{NewCLIArg(argTypeOption, "-")}, true},

		// TODO: Lets not return an option here
		// Lets let the caller identify this.
		{"double dash", "--", Normal, []*ProgramTree{NewCLIArg(argTypeTerminator, "--")}, false},
		{"double dash", "--", Bundling, []*ProgramTree{NewCLIArg(argTypeTerminator, "--")}, false},
		{"double dash", "--", SingleDash, []*ProgramTree{NewCLIArg(argTypeTerminator, "--")}, false},

		{"no option", "opt", Normal, []*ProgramTree{}, false},
		{"no option", "opt", Bundling, []*ProgramTree{}, false},
		{"no option", "opt", SingleDash, []*ProgramTree{}, false},

		{"Long option", "--opt", Normal, []*ProgramTree{NewCLIArg(argTypeOption, "opt")}, true},
		{"Long option", "--opt", Bundling, []*ProgramTree{NewCLIArg(argTypeOption, "opt")}, true},
		{"Long option", "--opt", SingleDash, []*ProgramTree{NewCLIArg(argTypeOption, "opt")}, true},

		{"Long option with arg", "--opt=arg", Normal, []*ProgramTree{NewCLIArg(argTypeOption, "opt", "arg")}, true},
		{"Long option with arg", "--opt=arg", Bundling, []*ProgramTree{NewCLIArg(argTypeOption, "opt", "arg")}, true},
		{"Long option with arg", "--opt=arg", SingleDash, []*ProgramTree{NewCLIArg(argTypeOption, "opt", "arg")}, true},

		{"short option", "-opt", Normal, []*ProgramTree{NewCLIArg(argTypeOption, "opt")}, true},
		{"short option", "-opt", Bundling, []*ProgramTree{NewCLIArg(argTypeOption, "o"), NewCLIArg(argTypeOption, "p"), NewCLIArg(argTypeOption, "t")}, true},
		{"short option", "-opt", SingleDash, []*ProgramTree{NewCLIArg(argTypeOption, "o", "pt")}, true},

		{"short option with arg", "-opt=arg", Normal, []*ProgramTree{NewCLIArg(argTypeOption, "opt", "arg")}, true},
		{"short option with arg", "-opt=arg", Bundling, []*ProgramTree{NewCLIArg(argTypeOption, "o"), NewCLIArg(argTypeOption, "p"), NewCLIArg(argTypeOption, "t", "arg")}, true},
		{"short option with arg", "-opt=arg", SingleDash, []*ProgramTree{NewCLIArg(argTypeOption, "o", "pt=arg")}, true},
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
