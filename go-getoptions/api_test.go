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
	cmd1 := opt.NewCommand("cmd1", "")
	cmd1.String("cmd1opt1", "")
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
			mode     Mode
			expected *CLIArg
		}{
			{"empty", nil, Normal, &CLIArg{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Args:     []string{},
				Children: []*CLIArg{},
			}},
			{"empty", []string{}, Normal, &CLIArg{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Args:     []string{},
				Children: []*CLIArg{},
			}},
			{"option", []string{"--opt1"}, Normal, &CLIArg{
				Type: argTypeProgname,
				Name: os.Args[0],
				Args: []string{"--opt1"},
				Children: []*CLIArg{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Args:     []string{},
						Children: []*CLIArg{},
					},
				},
			}},
			{"terminator", []string{"--", "--opt1"}, Normal, &CLIArg{
				Type: argTypeProgname,
				Name: os.Args[0],
				Args: []string{"--", "--opt1"},
				Children: []*CLIArg{
					{
						Type:     argTypeText,
						Name:     "--opt1",
						Args:     []string{},
						Children: []*CLIArg{},
					},
				},
			}},
			{"command", []string{"--opt1", "cmd1", "--cmd1opt1", "hello", "world"}, Normal, &CLIArg{
				Type: argTypeProgname,
				Name: os.Args[0],
				Args: []string{"--opt1", "cmd1", "--cmd1opt1", "hello", "world"},
				Children: []*CLIArg{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Args:     []string{},
						Children: []*CLIArg{},
					},
					{
						Type: argTypeCommand,
						Name: "cmd1",
						Args: []string{},
						Children: []*CLIArg{
							{
								Type:     argTypeOption,
								Name:     "cmd1opt1",
								Args:     []string{},
								Children: []*CLIArg{},
							},
							{
								Type:     argTypeText,
								Name:     "hello",
								Args:     []string{},
								Children: []*CLIArg{},
							},
							{
								Type:     argTypeText,
								Name:     "world",
								Args:     []string{},
								Children: []*CLIArg{},
							},
						},
					},
				},
			}},
			// {"arg", []string{"hello", "world"}, Normal, &CLIArg{
			// 	Type: argTypeProgname,
			// 	Name: os.Args[0],
			// 	Args: []string{"text"},
			// 	Children: []*CLIArg{
			// 		{
			// 			Type:     argTypeText,
			// 			Name:     "hello",
			// 			Args:     []string{},
			// 			Children: []*CLIArg{},
			// 		},
			// 		{
			// 			Type:     argTypeText,
			// 			Name:     "world",
			// 			Args:     []string{},
			// 			Children: []*CLIArg{},
			// 		},
			// 	},
			// }},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				argTree := parseCLIArgs(tree, test.args, test.mode)
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
		name     string
		inputArg string
		mode     Mode
		expected []*CLIArg
		isOption bool
	}{
		{"lone dash", "-", Normal, []*CLIArg{NewCLIArg(argTypeOption, "-")}, true},
		{"lone dash", "-", Bundling, []*CLIArg{NewCLIArg(argTypeOption, "-")}, true},
		{"lone dash", "-", SingleDash, []*CLIArg{NewCLIArg(argTypeOption, "-")}, true},

		// TODO: Lets not return an option here
		// Lets let the caller identify this.
		{"double dash", "--", Normal, []*CLIArg{NewCLIArg(argTypeTerminator, "--")}, false},
		{"double dash", "--", Bundling, []*CLIArg{NewCLIArg(argTypeTerminator, "--")}, false},
		{"double dash", "--", SingleDash, []*CLIArg{NewCLIArg(argTypeTerminator, "--")}, false},

		{"no option", "opt", Normal, []*CLIArg{}, false},
		{"no option", "opt", Bundling, []*CLIArg{}, false},
		{"no option", "opt", SingleDash, []*CLIArg{}, false},

		{"Long option", "--opt", Normal, []*CLIArg{NewCLIArg(argTypeOption, "opt")}, true},
		{"Long option", "--opt", Bundling, []*CLIArg{NewCLIArg(argTypeOption, "opt")}, true},
		{"Long option", "--opt", SingleDash, []*CLIArg{NewCLIArg(argTypeOption, "opt")}, true},

		{"Long option with arg", "--opt=arg", Normal, []*CLIArg{NewCLIArg(argTypeOption, "opt", "arg")}, true},
		{"Long option with arg", "--opt=arg", Bundling, []*CLIArg{NewCLIArg(argTypeOption, "opt", "arg")}, true},
		{"Long option with arg", "--opt=arg", SingleDash, []*CLIArg{NewCLIArg(argTypeOption, "opt", "arg")}, true},

		{"short option", "-opt", Normal, []*CLIArg{NewCLIArg(argTypeOption, "opt")}, true},
		{"short option", "-opt", Bundling, []*CLIArg{NewCLIArg(argTypeOption, "o"), NewCLIArg(argTypeOption, "p"), NewCLIArg(argTypeOption, "t")}, true},
		{"short option", "-opt", SingleDash, []*CLIArg{NewCLIArg(argTypeOption, "o", "pt")}, true},

		{"short option with arg", "-opt=arg", Normal, []*CLIArg{NewCLIArg(argTypeOption, "opt", "arg")}, true},
		{"short option with arg", "-opt=arg", Bundling, []*CLIArg{NewCLIArg(argTypeOption, "o"), NewCLIArg(argTypeOption, "p"), NewCLIArg(argTypeOption, "t", "arg")}, true},
		{"short option with arg", "-opt=arg", SingleDash, []*CLIArg{NewCLIArg(argTypeOption, "o", "pt=arg")}, true},
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
