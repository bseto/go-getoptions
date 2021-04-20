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

func setupOpt() *GetOpt {
	opt := New()
	opt.String("rootopt1", "")

	cmd1 := opt.NewCommand("cmd1", "")
	cmd1.String("cmd1opt1", "")
	cmd2 := opt.NewCommand("cmd2", "")
	cmd2.String("cmd2opt1", "")

	sub1cmd1 := cmd1.NewCommand("sub1cmd1", "")
	sub1cmd1.String("sub1cmd1opt1", "")
	return opt
}

func TestTrees(t *testing.T) {
	buf := setupLogging()

	t.Run("programTree", func(t *testing.T) {
		tree := setupOpt().programTree
		root := &programTree{
			Type:     argTypeProgname,
			Name:     os.Args[0],
			Children: []*programTree{},
		}
		rootopt1 := &programTree{
			Type:     argTypeOption,
			Name:     "rootopt1",
			Parent:   root,
			Children: []*programTree{},
		}
		cmd1 := &programTree{
			Type:     argTypeCommand,
			Name:     "cmd1",
			Parent:   root,
			Children: []*programTree{},
			Level:    1,
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
			Level:    2,
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
			Level:    1,
		}
		cmd2opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "cmd2opt1",
			Parent:   cmd2,
			Children: []*programTree{},
		}
		root.Children = append(root.Children, []*programTree{
			rootopt1,
			cmd1,
			cmd2,
		}...)
		cmd1.Children = append(cmd1.Children, []*programTree{
			rootopt1,
			cmd1opt1,
			sub1cmd1,
		}...)
		sub1cmd1.Children = append(sub1cmd1.Children, []*programTree{
			rootopt1,
			cmd1opt1,
			sub1cmd1opt1,
		}...)
		cmd2.Children = append(cmd2.Children, []*programTree{
			rootopt1,
			cmd2opt1,
		}...)

		if !reflect.DeepEqual(root, tree) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err := getNode(tree)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(root, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, []string{}...)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(root, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, "cmd1")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(cmd1, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, "cmd1", "sub1cmd1")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(sub1cmd1, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, "cmd2")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(cmd2, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

	})

	t.Run("CLIArg", func(t *testing.T) {

		tests := []struct {
			name     string
			args     []string
			mode     Mode
			expected *programTree
		}{
			{"empty", nil, Normal, setupOpt().programTree},
			{"empty", []string{}, Normal, setupOpt().programTree},
			{"text", []string{"txt"}, Normal, func() *programTree {
				tree := setupOpt().programTree
				tree.Children = append(tree.Children, newCLIArg(argTypeText, "txt"))
				return tree
			}()},
			{"command", []string{"cmd1"}, Normal, func() *programTree {
				n, err := getNode(setupOpt().programTree, "cmd1")
				if err != nil {
					panic(err)
				}
				return n
			}()},
			{"text to command", []string{"cmd1", "txt"}, Normal, func() *programTree {
				tree := setupOpt().programTree
				n, err := getNode(tree, "cmd1")
				if err != nil {
					panic(err)
				}
				n.Children = append(n.Children, newCLIArg(argTypeText, "txt"))
				return n
			}()},
			{"text to sub command", []string{"cmd1", "sub1cmd1", "txt"}, Normal, func() *programTree {
				tree := setupOpt().programTree
				n, err := getNode(tree, "cmd1", "sub1cmd1")
				if err != nil {
					panic(err)
				}
				n.Children = append(n.Children, newCLIArg(argTypeText, "txt"))
				return n
			}()},
			// {"option", []string{"--rootopt1"}, Normal, func() *programTree {
			// 	n, _ := getNode(tree)
			// 	// n.option.Args = []string{"--rootopt1"}
			// 	return n
			// }()},
			// {"terminator", []string{"--", "--opt1"}, Normal, &programTree{
			// 	Type:   argTypeProgname,
			// 	Name:   os.Args[0],
			// 	option: option{Args: []string{"--", "--opt1"}},
			// 	Children: []*programTree{
			// 		{
			// 			Type:     argTypeText,
			// 			Name:     "--opt1",
			// 			option:   option{Args: []string{}},
			// 			Children: []*programTree{},
			// 		},
			// 	},
			// }},
			// {"command", []string{"--opt1", "cmd1", "--cmd1opt1"}, Normal, &programTree{
			// 	Type:   argTypeProgname,
			// 	Name:   os.Args[0],
			// 	option: option{Args: []string{"--opt1", "cmd1", "--cmd1opt1"}},
			// 	Children: []*programTree{
			// 		{
			// 			Type:     argTypeOption,
			// 			Name:     "opt1",
			// 			option:   option{Args: []string{}},
			// 			Children: []*programTree{},
			// 		},
			// 		{
			// 			Type:   argTypeCommand,
			// 			Name:   "cmd1",
			// 			option: option{Args: []string{}},
			// 			Children: []*programTree{
			// 				{
			// 					Type:     argTypeOption,
			// 					Name:     "cmd1opt1",
			// 					option:   option{Args: []string{}},
			// 					Children: []*programTree{},
			// 				},
			// 			},
			// 		},
			// 	},
			// }},
			// {"subcommand", []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"}, Normal, &programTree{
			// 	Type:   argTypeProgname,
			// 	Name:   os.Args[0],
			// 	option: option{Args: []string{"--opt1", "cmd1", "--cmd1opt1", "sub1cmd1", "--sub1cmd1opt1"}},
			// 	Children: []*programTree{
			// 		{
			// 			Type:     argTypeOption,
			// 			Name:     "opt1",
			// 			option:   option{Args: []string{}},
			// 			Children: []*programTree{},
			// 		},
			// 		{
			// 			Type:   argTypeCommand,
			// 			Name:   "cmd1",
			// 			option: option{Args: []string{}},
			// 			Children: []*programTree{
			// 				{
			// 					Type:     argTypeOption,
			// 					Name:     "cmd1opt1",
			// 					option:   option{Args: []string{}},
			// 					Children: []*programTree{},
			// 				},
			// 				{
			// 					Type:   argTypeCommand,
			// 					Name:   "sub1cmd1",
			// 					option: option{Args: []string{}},
			// 					Children: []*programTree{
			// 						{
			// 							Type:     argTypeOption,
			// 							Name:     "sub1cmd1opt1",
			// 							option:   option{Args: []string{}},
			// 							Children: []*programTree{},
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// }},
			// {"arg", []string{"hello", "world"}, Normal, &programTree{
			// 	Type:   argTypeProgname,
			// 	Name:   os.Args[0],
			// 	option: option{Args: []string{"hello", "world"}},
			// 	Children: []*programTree{
			// 		{
			// 			Type:     argTypeText,
			// 			Name:     "hello",
			// 			option:   option{Args: []string{}},
			// 			Children: []*programTree{},
			// 		},
			// 		{
			// 			Type:     argTypeText,
			// 			Name:     "world",
			// 			option:   option{Args: []string{}},
			// 			Children: []*programTree{},
			// 		},
			// 	},
			// }},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				tree := setupOpt().programTree
				argTree := parseCLIArgs(false, tree, test.args, test.mode)
				if !reflect.DeepEqual(test.expected, argTree) {
					t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(test.expected), spew.Sdump(argTree))
					// t.Errorf("expected tree: \n%s\n got: \n%s\n", test.expected, argTree)
				}
			})
		}
	})

	t.Cleanup(func() { t.Log(buf.String()) })
}
