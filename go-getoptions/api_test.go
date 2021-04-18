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
		root := &programTree{
			Type:     argTypeProgname,
			Name:     os.Args[0],
			Children: []*programTree{},
		}
		opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "opt1",
			Parent:   root,
			Children: []*programTree{},
		}
		cmd1 := &programTree{
			Type:     argTypeCommand,
			Name:     "cmd1",
			Parent:   root,
			Children: []*programTree{},
		}
		cmd1opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "cmd1opt1",
			Parent:   cmd1,
			Children: []*programTree{},
		}
		sub1cmd1 := &programTree{
			Type:     argTypeCommand,
			Name:     "sub1cmd1",
			Parent:   cmd1,
			Children: []*programTree{},
		}
		sub1cmd1opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "sub1cmd1opt1",
			Parent:   sub1cmd1,
			Children: []*programTree{},
		}
		cmd2 := &programTree{
			Type:     argTypeCommand,
			Name:     "cmd2",
			Parent:   root,
			Children: []*programTree{},
		}
		cmd2opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "cmd2opt1",
			Parent:   cmd2,
			Children: []*programTree{},
		}
		root.Children = append(root.Children, []*programTree{
			opt1,
			cmd1,
			cmd2,
		}...)
		cmd1.Children = append(cmd1.Children, []*programTree{
			opt1,
			cmd1opt1,
			sub1cmd1,
		}...)
		sub1cmd1.Children = append(sub1cmd1.Children, []*programTree{
			opt1,
			cmd1opt1,
			sub1cmd1opt1,
		}...)
		cmd2.Children = append(cmd2.Children, []*programTree{
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
			expected *programTree
		}{
			{"empty", nil, Normal, &programTree{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Option:   Option{Args: []string{}},
				Children: []*programTree{},
			}},
			{"empty", []string{}, Normal, &programTree{
				Type:     argTypeProgname,
				Name:     os.Args[0],
				Option:   Option{Args: []string{}},
				Children: []*programTree{},
			}},
			{"option", []string{"--opt1"}, Normal, &programTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--opt1"}},
				Children: []*programTree{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Option:   Option{Args: []string{}},
						Children: []*programTree{},
					},
				},
			}},
			{"terminator", []string{"--", "--opt1"}, Normal, &programTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--", "--opt1"}},
				Children: []*programTree{
					{
						Type:     argTypeText,
						Name:     "--opt1",
						Option:   Option{Args: []string{}},
						Children: []*programTree{},
					},
				},
			}},
			{"command", []string{"--opt1", "cmd1", "--cmd1opt1"}, Normal, &programTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--opt1", "cmd1", "--cmd1opt1"}},
				Children: []*programTree{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Option:   Option{Args: []string{}},
						Children: []*programTree{},
					},
					{
						Type:   argTypeCommand,
						Name:   "cmd1",
						Option: Option{Args: []string{}},
						Children: []*programTree{
							{
								Type:     argTypeOption,
								Name:     "cmd1opt1",
								Option:   Option{Args: []string{}},
								Children: []*programTree{},
							},
						},
					},
				},
			}},
			{"subcommand", []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"}, Normal, &programTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"}},
				Children: []*programTree{
					{
						Type:     argTypeOption,
						Name:     "opt1",
						Option:   Option{Args: []string{}},
						Children: []*programTree{},
					},
					{
						Type:   argTypeCommand,
						Name:   "cmd1",
						Option: Option{Args: []string{}},
						Children: []*programTree{
							{
								Type:     argTypeOption,
								Name:     "cmd1opt1",
								Option:   Option{Args: []string{}},
								Children: []*programTree{},
							},
							{
								Type:   argTypeCommand,
								Name:   "sub1cmd1",
								Option: Option{Args: []string{}},
								Children: []*programTree{
									{
										Type:     argTypeOption,
										Name:     "sub1cmd1opt1",
										Option:   Option{Args: []string{}},
										Children: []*programTree{},
									},
								},
							},
						},
					},
				},
			}},
			{"arg", []string{"hello", "world"}, Normal, &programTree{
				Type:   argTypeProgname,
				Name:   os.Args[0],
				Option: Option{Args: []string{"hello", "world"}},
				Children: []*programTree{
					{
						Type:     argTypeText,
						Name:     "hello",
						Option:   Option{Args: []string{}},
						Children: []*programTree{},
					},
					{
						Type:     argTypeText,
						Name:     "world",
						Option:   Option{Args: []string{}},
						Children: []*programTree{},
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
