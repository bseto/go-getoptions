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

	sub1cmd1 := cmd1.NewCommand("sub1cmd1", "")
	sub1cmd1.String("sub1cmd1opt1", "")

	tree := opt.programTree
	t.Run("programTree", func(t *testing.T) {
		root := &ProgramTree{
			Type:     argTypeProgname,
			Name:     os.Args[0],
			Children: []*ProgramTree{},
		}
		opt1 := &ProgramTree{
			Type:     argTypeOption,
			Name:     "opt1",
			Parent:   root,
			Children: []*ProgramTree{},
		}
		cmd1 := &ProgramTree{
			Type:     argTypeCommand,
			Name:     "cmd1",
			Parent:   root,
			Children: []*ProgramTree{},
		}
		cmd1opt1 := &ProgramTree{
			Type:     argTypeOption,
			Name:     "cmd1opt1",
			Parent:   cmd1,
			Children: []*ProgramTree{},
		}
		sub1cmd1 := &ProgramTree{
			Type:     argTypeCommand,
			Name:     "sub1cmd1",
			Parent:   cmd1,
			Children: []*ProgramTree{},
		}
		sub1cmd1opt1 := &ProgramTree{
			Type:     argTypeOption,
			Name:     "sub1cmd1opt1",
			Parent:   sub1cmd1,
			Children: []*ProgramTree{},
		}
		cmd2 := &ProgramTree{
			Type:     argTypeCommand,
			Name:     "cmd2",
			Parent:   root,
			Children: []*ProgramTree{},
		}
		cmd2opt1 := &ProgramTree{
			Type:     argTypeOption,
			Name:     "cmd2opt1",
			Parent:   cmd2,
			Children: []*ProgramTree{},
		}
		root.Children = append(root.Children, []*ProgramTree{
			opt1,
			cmd1,
			cmd2,
		}...)
		cmd1.Children = append(cmd1.Children, []*ProgramTree{
			opt1,
			cmd1opt1,
			sub1cmd1,
		}...)
		sub1cmd1.Children = append(sub1cmd1.Children, []*ProgramTree{
			opt1,
			cmd1opt1,
			sub1cmd1opt1,
		}...)
		cmd2.Children = append(cmd2.Children, []*ProgramTree{
			opt1,
			cmd2opt1,
		}...)

		if !reflect.DeepEqual(root, tree) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
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
			{"command", []string{"--opt1", "cmd1", "--cmd1opt1"}, Normal, &CLIArg{
				Type: argTypeProgname,
				Name: os.Args[0],
				Args: []string{"--opt1", "cmd1", "--cmd1opt1"},
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
						},
					},
				},
			}},
			{"subcommand", []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"}, Normal, &CLIArg{
				Type: argTypeProgname,
				Name: os.Args[0],
				Args: []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"},
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
								Type: argTypeCommand,
								Name: "sub1cmd1",
								Args: []string{},
								Children: []*CLIArg{
									{
										Type:     argTypeOption,
										Name:     "sub1cmd1opt1",
										Args:     []string{},
										Children: []*CLIArg{},
									},
								},
							},
						},
					},
				},
			}},
			{"arg", []string{"hello", "world"}, Normal, &CLIArg{
				Type: argTypeProgname,
				Name: os.Args[0],
				Args: []string{"hello", "world"},
				Children: []*CLIArg{
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
			}},
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
