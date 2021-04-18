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

	setupOpt := func() *GetOpt {
		opt := New()
		opt.String("opt1", "")

		cmd1 := opt.NewCommand("cmd1", "")
		cmd1.String("cmd1opt1", "")
		cmd2 := opt.NewCommand("cmd2", "")
		cmd2.String("cmd2opt1", "")

		sub1cmd1 := cmd1.NewCommand("sub1cmd1", "")
		sub1cmd1.String("sub1cmd1opt1", "")
		return opt
	}

	t.Run("programTree", func(t *testing.T) {
		tree := setupOpt().programTree
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
		tree := setupOpt().programTree

		tests := []struct {
			name     string
			args     []string
			mode     Mode
			expected *ProgramTree
		}{
			{"empty", nil, Normal, &ProgramTree{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Option:   Option{Args: []string{}},
				Children: []*ProgramTree{},
			}},
			{"empty", []string{}, Normal, &ProgramTree{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Option:   Option{Args: []string{}},
				Children: []*ProgramTree{},
			}},
			{"option", []string{"--opt1"}, Normal, &ProgramTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--opt1"}},
				Children: []*ProgramTree{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Option:   Option{Args: []string{}},
						Children: []*ProgramTree{},
					},
				},
			}},
			{"terminator", []string{"--", "--opt1"}, Normal, &ProgramTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--", "--opt1"}},
				Children: []*ProgramTree{
					{
						Type:     argTypeText,
						Name:     "--opt1",
						Option:   Option{Args: []string{}},
						Children: []*ProgramTree{},
					},
				},
			}},
			{"command", []string{"--opt1", "cmd1", "--cmd1opt1"}, Normal, &ProgramTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--opt1", "cmd1", "--cmd1opt1"}},
				Children: []*ProgramTree{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Option:   Option{Args: []string{}},
						Children: []*ProgramTree{},
					},
					{
						Type:   argTypeCommand,
						Name:   "cmd1",
						Option: Option{Args: []string{}},
						Children: []*ProgramTree{
							{
								Type:     argTypeOption,
								Name:     "cmd1opt1",
								Option:   Option{Args: []string{}},
								Children: []*ProgramTree{},
							},
						},
					},
				},
			}},
			{"subcommand", []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"}, Normal, &ProgramTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"}},
				Children: []*ProgramTree{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Option:   Option{Args: []string{}},
						Children: []*ProgramTree{},
					},
					{
						Type:   argTypeCommand,
						Name:   "cmd1",
						Option: Option{Args: []string{}},
						Children: []*ProgramTree{
							{
								Type:     argTypeOption,
								Name:     "cmd1opt1",
								Option:   Option{Args: []string{}},
								Children: []*ProgramTree{},
							},
							{
								Type:   argTypeCommand,
								Name:   "sub1cmd1",
								Option: Option{Args: []string{}},
								Children: []*ProgramTree{
									{
										Type:     argTypeOption,
										Name:     "sub1cmd1opt1",
										Option:   Option{Args: []string{}},
										Children: []*ProgramTree{},
									},
								},
							},
						},
					},
				},
			}},
			{"arg", []string{"hello", "world"}, Normal, &ProgramTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"hello", "world"}},
				Children: []*ProgramTree{
					{
						Type:     argTypeText,
						Name:     "hello",
						Option:   Option{Args: []string{}},
						Children: []*ProgramTree{},
					},
					{
						Type:     argTypeText,
						Name:     "world",
						Option:   Option{Args: []string{}},
						Children: []*ProgramTree{},
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
