package getoptions

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func setupLogging() *bytes.Buffer {
	spew.Config = spew.ConfigState{
		Indent:                  "  ",
		MaxDepth:                0,
		DisableMethods:          false,
		DisablePointerMethods:   false,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		ContinueOnMethod:        false,
		SortKeys:                true,
		SpewKeys:                false,
	}
	s := ""
	buf := bytes.NewBufferString(s)
	Logger.SetOutput(buf)
	return buf
}

func TestTrees(t *testing.T) {
	buf := setupLogging()
	opt := New()
	opt.String("opt1", "")
	cmd := opt.NewCommand("cmd1", "")
	cmd.String("cmd1opt1", "")
	cmd2 := opt.NewCommand("cmd2", "")
	cmd2.String("cmd2opt1", "")
	tree := opt.cliTree
	t.Run("CLITree", func(t *testing.T) {
		expectedTree := &CLITree{
			Type:     argTypeProgname,
			Name:     os.Args[0],
			Children: []*CLITree{},
		}
		expectedTreeCmd1 := &CLITree{
			Type:     argTypeCommand,
			Name:     "cmd1",
			Parent:   expectedTree,
			Children: []*CLITree{},
		}
		expectedTreeCmd1.Children = append(expectedTreeCmd1.Children, &CLITree{
			Type:     argTypeOption,
			Name:     "cmd1opt1",
			Parent:   expectedTreeCmd1,
			Children: []*CLITree{},
		})
		expectedTreeCmd2 := &CLITree{
			Type:     argTypeCommand,
			Name:     "cmd2",
			Parent:   expectedTree,
			Children: []*CLITree{},
		}
		expectedTreeCmd2.Children = append(expectedTreeCmd2.Children, &CLITree{
			Type:     argTypeOption,
			Name:     "cmd2opt1",
			Parent:   expectedTreeCmd2,
			Children: []*CLITree{},
		})
		expectedTree.Children = append(expectedTree.Children, []*CLITree{
			{
				Type:     argTypeOption,
				Name:     "opt1",
				Parent:   expectedTree,
				Children: []*CLITree{},
			},
			expectedTreeCmd1,
			expectedTreeCmd2,
		}...)

		if !reflect.DeepEqual(expectedTree, tree) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(expectedTree), spew.Sdump(tree))
		}
	})

	t.Run("CLIArg", func(t *testing.T) {

		tests := []struct {
			name     string
			args     []string
			expected *CLIArg
		}{
			{"empty", nil, &CLIArg{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Args:     []string{},
				Children: []*CLIArg{},
			}},
			{"empty", []string{}, &CLIArg{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Args:     []string{},
				Children: []*CLIArg{},
			}},
			{"arg", []string{"opt1"}, &CLIArg{
				Type: argTypeProgname,
				Name: os.Args[0],
				Args: []string{"opt1"},
				Children: []*CLIArg{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Args:     []string{},
						Children: []*CLIArg{},
					},
				},
			}},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				argTree := parseCLIArgs(tree, test.args)
				if !reflect.DeepEqual(test.expected, argTree) {
					t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(test.expected), spew.Sdump(argTree))
				}
			})
		}
	})

	t.Cleanup(func() { t.Log(buf.String()) })
}

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
				t.Errorf("isOption(%q, %d) == (%q, %v), want (%q, %v)",
					tt.in, tt.mode, optPair, is, tt.optPair, tt.is)
			}
			t.Log(buf.String())
		})
	}
}
