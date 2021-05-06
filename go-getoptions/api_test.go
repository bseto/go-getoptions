package getoptions

import (
	"reflect"
	"testing"
)

func TestAPI(t *testing.T) {
	buf := setupLogging()

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
				tree.Children = append(tree.Children, newCLIArg(tree, argTypeText, "txt"))
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
				n, err := getNode(setupOpt().programTree, "cmd1")
				if err != nil {
					panic(err)
				}
				n.Children = append(n.Children, newCLIArg(n, argTypeText, "txt"))
				return n
			}()},
			{"text to sub command", []string{"cmd1", "sub1cmd1", "txt"}, Normal, func() *programTree {
				n, err := getNode(setupOpt().programTree, "cmd1", "sub1cmd1")
				if err != nil {
					panic(err)
				}
				n.Children = append(n.Children, newCLIArg(n, argTypeText, "txt"))
				return n
			}()},
			{"option", []string{"--rootopt1"}, Normal, func() *programTree {
				tree := setupOpt().programTree
				return tree
			}()},
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
				argTree, _, err := parseCLIArgs(false, tree, test.args, test.mode)
				if err != nil {
					t.Errorf("unexpected error")
				}
				if !reflect.DeepEqual(test.expected, argTree) {
					t.Errorf("expected tree: %s, got: %s\n", SpewToFile(t, test.expected, "expected"), SpewToFile(t, argTree, "got"))
					t.Fatalf("expected tree: \n%s\n got: \n%s\n", test.expected.Str(), argTree.Str())
				}
			})
		}
	})

	t.Cleanup(func() { t.Log(buf.String()) })
}
