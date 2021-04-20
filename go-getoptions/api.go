package getoptions

import "github.com/DavidGamba/go-getoptions/sliceiterator"

type programTree struct {
	Type     argType
	Name     string
	Children []*programTree
	Parent   *programTree
	Level    int
	option
	command
}

// func (t *programTree) String() string {
// 	out := strings.Repeat("  ", t.Level) + fmt.Sprintf("Name: %v, Type: %v", t.Name, t.Type)
// 	if len(t.Children) > 0 {
// 		out += ", children: [\n"
// 		for _, child := range t.Children {
// 			out += child.String()
// 		}
// 		out += strings.Repeat("  ", t.Level) + "]\n"
// 	} else {
// 		out += ", children: []\n"
// 	}
// 	return out
// }

func getNode(tree *programTree, element ...string) (*programTree, error) {
	if len(element) == 0 {
		return tree, nil
	}
	currentNode := tree
	for _, child := range currentNode.Children {
		if child.Name == element[0] {
			return getNode(child, element[1:]...)
		}
	}
	return tree, nil
}

type argType int

const (
	argTypeProgname   argType = iota // The root node type
	argTypeCommand                   // The node type used for commands and subcommands
	argTypeOption                    // The node type used for options
	argTypeText                      // The node type used for regular cli arguments
	argTypeTerminator                // --
)

// option - Fields that only make sense for an option
type option struct {
	Aliases  []string
	Args     []string
	Called   bool
	CalledAs string
	Min, Max int // Minimum and Maximun amount of fields to pass to option in one call.
}

// command - Fields that only make sense for a command
type command struct {
	CommandFn CommandFn
}

func newCLIArg(t argType, name string, args ...string) *programTree {
	arg := &programTree{
		Type:     t,
		Name:     name,
		Children: []*programTree{},
		option:   option{Args: []string{}},
	}
	if len(args) > 0 {
		arg.Args = args
	}
	return arg
}

func parseCLIArgs(tree *programTree, args []string, mode Mode) *programTree {
	// Design: This function could return an array or CLIargs as a parse result
	// or I could do one level up and have a root CLIarg type with the name of
	// the program.  Having the root level might be helpful with help generation.

	// The current implementation expects os.Args[1:] as an argument so this
	// can't expect to receive the os.Args[0] as the first argument.

	// CLI arguments are split by spaces by the shell and passed as individual
	// strings.  In most cases, a cli argument (one string) represents one option
	// or one argument, however, in the case of bundling mode a single string can
	// represent multiple options.

	// Ensure consistent response for empty and nil slices
	if args == nil {
		args = []string{}
	}

	currentProgramNode := tree

	iterator := sliceiterator.New(&args)

ARGS_LOOP:
	for iterator.Next() {

		// handle terminator
		if iterator.Value() == "--" {
			for iterator.Next() {
				currentProgramNode.Children = append(currentProgramNode.Children, newCLIArg(argTypeText, iterator.Value()))
			}
			break
		}

		// TODO: Handle lonesome dash

		// TODO: Handle case where option has an argument
		// check for option
		if cliArg, is := isOption(iterator.Value(), mode); is {
			currentProgramNode.Children = append(currentProgramNode.Children, cliArg...)
			continue
		}

		// handle commands and subcommands
		for _, child := range currentProgramNode.Children {
			// Only check commands
			if child.Type == argTypeCommand && child.Name == iterator.Value() {
				currentProgramNode = child
				continue ARGS_LOOP
			}
		}

		// handle text
		currentProgramNode.Children = append(currentProgramNode.Children, newCLIArg(argTypeText, iterator.Value()))
	}
	return currentProgramNode
}

// TODO:
// suggestCompletions -
func suggestCompletions(tree *programTree, args []string, mode Mode) {}
