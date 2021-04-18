package getoptions

import (
	"os"
)

type argType int

const (
	argTypeProgname   argType = iota // The root node type
	argTypeCommand                   // The node type used for commands and subcommands
	argTypeOption                    // The node type used for options
	argTypeText                      // The node type used for regular cli arguments
	argTypeTerminator                // --
)

func NewCLIArg(t argType, name string, args ...string) *ProgramTree {
	arg := &ProgramTree{
		Type:     t,
		Name:     name,
		Children: []*ProgramTree{},
		Option:   Option{Args: []string{}},
	}
	if len(args) > 0 {
		arg.Args = args
	}
	return arg
}

type ProgramTree struct {
	Type     argType
	Name     string
	Children []*ProgramTree
	Parent   *ProgramTree
	Option
	Command
}

// Option - Fields that only make sense for an Option
type Option struct {
	Aliases  []string
	Args     []string
	Called   bool
	CalledAs string
	Min, Max int // Minimum and Maximun amount of fields to pass to option in one call.
}

// Command - Fields that only make sense for a Command
type Command struct {
	CommandFn CommandFn
}

func parseCLIArgs(tree *ProgramTree, args []string, mode Mode) *ProgramTree {
	// Design: This function could return an array or CLIargs as a parse result
	// or I could do one level up and have a root CLIarg type with the name of
	// the program.  Having the root level might be helpful with help generation.

	// The current implementation expects os.Args[1:] as an argument so this
	// can't expect to receive the os.Args[0] as the first argument.

	// CLI arguments are split by spaces by the shell and passed as individual
	// strings.  In most cases, a cli argument (one string) represents one option
	// or one argument, however, in the case of bundling mode a single string can
	// represent multiple options.

	// When parsing the cli args, there is no way to tell apart a command vs just
	// text input to the program, one argument to this parser needs to be the
	// command tree.

	// TODO: Question: How is text input before a command handled? Is it allowed?

	// Ensure consistent response for empty and nil slices
	if args == nil {
		args = []string{}
	}

	root := NewCLIArg(argTypeProgname, os.Args[0], args...)

	currentCLINode := root
	currentProgramNode := tree

ARGS_LOOP:
	for i, arg := range args {

		// handle terminator
		if arg == "--" {
			if len(args) > i+1 {
				for _, arg := range args[i+1:] {
					// TODO: I am not checking the option against the tree here.
					currentCLINode.Children = append(currentCLINode.Children, NewCLIArg(argTypeText, arg))
				}
			}
			break
		}

		// TODO: Handle lonesome dash

		// TODO: Handle case where option has an argument
		// check for option
		cliArg, is := isOption(arg, mode)
		if is {
			currentCLINode.Children = append(currentCLINode.Children, cliArg...)
			continue
		}

		// handle commands and subcommands
		for _, child := range currentProgramNode.Children {
			// Only check commands
			if child.Type != argTypeCommand {
				continue
			}
			if child.Name == arg {
				cmd := NewCLIArg(argTypeCommand, arg)
				currentCLINode.Children = append(currentCLINode.Children, cmd)
				currentCLINode = cmd
				currentProgramNode = child
				continue ARGS_LOOP
			}
		}

		// handle text
		currentCLINode.Children = append(currentCLINode.Children, NewCLIArg(argTypeText, arg))
	}
	return root
}

// TODO:
// suggestCompletions -
func suggestCompletions(tree *ProgramTree, args []string, mode Mode) {}
