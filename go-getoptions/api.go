package getoptions

import (
	"fmt"
	"strings"

	"github.com/DavidGamba/go-getoptions/sliceiterator"
)

type programTree struct {
	Type     argType
	Name     string
	Children []*programTree
	Parent   *programTree
	Level    int
	command

	// The option node is passed around as a copy (so the parent can be redefined), however, the data is a pointer so it is modified at all levels.
	Option *option
}

// Str - not string so it doesn't get called automatically by Spew.
func (n *programTree) Str() string {
	level := n.Level
	if n.Type == argTypeOption {
		if n.Parent != nil {
			level = n.Parent.Level + 1
		}
	}
	out := strings.Repeat("  ", level) + fmt.Sprintf("Name: %v, Type: %v", n.Name, n.Type)
	if n.Parent != nil {
		out += fmt.Sprintf(", Parent: %v", n.Parent.Name)
	}
	if len(n.Children) > 0 {
		out += ", children: [\n"
		for _, child := range n.Children {
			out += child.Str()
		}
		out += strings.Repeat("  ", level) + "]\n"
	} else {
		out += ", children: []\n"
	}
	return out
}

// Copy - Returns a copy of programTree that maintains a pointer to the underlying data
func (n *programTree) Copy() *programTree {
	// a := *n
	// c := &a
	parent := *n.Parent
	c := &programTree{
		Type:     n.Type,
		Name:     n.Name,
		Children: n.Children,
		Parent:   &parent,
		Option:   n.Option,
	}
	return c
}

func (n *programTree) SetParent(p *programTree) *programTree {
	n.Parent = p
	return n
}

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

	// TODO: with unknown at one level it might become known at another level, what to do then?
	Unknown bool // Marks this option as one that wasn't declared or expected at this level

	Min, Max int // Minimum and Maximun amount of fields to pass to option in one call.

	// option data section
	Value   interface{}
	PString *string
}

// command - Fields that only make sense for a command
type command struct {
	CommandFn CommandFn
}

// TODO: Make this a method of tree so we can add parent information
func newCLIArg(parent *programTree, t argType, name string, args ...string) *programTree {
	arg := &programTree{
		Type:     t,
		Name:     name,
		Parent:   parent,
		Children: []*programTree{},
		Option:   &option{Args: []string{}},
	}
	if len(args) > 0 {
		arg.Option.Args = args
	}
	return arg
}

type completions *[]string

// parseCLIArgs - Given the root node tree and the cli args it returns a populated tree of the node that was called.
// For example, if a command is called, then the returned node is that of the command with the options that were set updated with their values.
// Additionally, when in completion mode, it returns the possible completions
func parseCLIArgs(completionMode bool, tree *programTree, args []string, mode Mode) (*programTree, completions, error) {
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

		// We only generate completions when we reached the end of the provided args
		if completionMode && iterator.IsLast() {
			// TODO: check what was the behaviour when you have a space and hit the tab completion.

			// TODO: Handle completions
			// We check to see if this is the last arg and act on that one.
			if iterator.Value() == "-" || iterator.Value() == "--" {
				// Provide option completions
			}
			if strings.HasPrefix(iterator.Value(), "-") {
				// Provide option completions
			}
			// Iterate over commands and check prefix to see if we offer command completion

			// Provide other kinds of completions, like file completions.
		}

		// handle terminator
		if iterator.Value() == "--" {
			for iterator.Next() {
				currentProgramNode.Children = append(currentProgramNode.Children, newCLIArg(currentProgramNode, argTypeText, iterator.Value()))
			}
			break
		}

		// TODO: Handle lonesome dash

		// TODO: Handle case where option has an argument
		// check for option

		// isOption should check if a cli argument starts with -.
		// If it does, we validate that it matches an option.
		// If it does we update the option with the values that might have been provided on the CLI.
		//
		// We almost need to build a separate option tree which allows unknown options and then update the main tree when we are done parsing cli args.
		//
		// Currently go-getoptions has no knowledge of command options at the
		// parents so it marks them as an unkonw option that needs to be used at a
		// different level. It is as if it was ignoring getoptions.Pass.
		if optPair, is := isOption(iterator.Value(), mode); is {
			// iterate over the possible cli args and try matching against expectations
			for _, a := range optPair {
				matches := 0
				for _, c := range currentProgramNode.Children {
					if c.Type != argTypeOption {
						continue
					}
					// handle full option match
					// TODO: handle partial matches
					if _, ok := stringSliceIndex(append([]string{c.Name}, c.Option.Aliases...), a.Option); ok {
						c.Option.Called = true
						c.Option.CalledAs = a.Option
						c.Option.Args = append(c.Option.Args, a.Args...)
						matches++
						// TODO: Handle option having a minimum bigger than 1
					}
				}
				if matches > 1 {
					// TODO: handle ambiguous option call error
					continue
				}
				if matches == 0 {
					// TODO: This is a new option, add it as a children and mark it as unknown
					// TODO: This shouldn't append new children but update existing ones and isOption needs to be able to check if the option expects a follow up argument.
					// Check min, check max and keep ingesting until something starts with `-` or matches a command.

					opt := newCLIArg(currentProgramNode, argTypeOption, a.Option, a.Args...)
					opt.Option.Unknown = true
					currentProgramNode.Children = append(currentProgramNode.Children, opt)
				}
			}
			continue
		}

		// When handling options out of order, iterate over all possible options for all the children and set them if they match.
		// That means that the option has to match the alias and aliases need to be non ambiguous with the parent.
		// partial options can only be applied if they match a single possible option in the tree.
		// Since at the end we return the programTree node, we will only care about handling the options at one single level.

		// handle commands and subcommands
		for _, child := range currentProgramNode.Children {
			// Only check commands
			if child.Type == argTypeCommand && child.Name == iterator.Value() {
				currentProgramNode = child
				continue ARGS_LOOP
			}
		}

		// handle text
		currentProgramNode.Children = append(currentProgramNode.Children, newCLIArg(currentProgramNode, argTypeText, iterator.Value()))
	}

	// TODO: Before returning the current node, parse EnvVars and update the values.

	// TODO: After being done parsing everything validate for errors
	// Errors can be unknown options, options without values, etc

	return currentProgramNode, &[]string{}, nil
}

// TODO:
// suggestCompletions -
func suggestCompletions(tree *programTree, args []string, mode Mode) {}
