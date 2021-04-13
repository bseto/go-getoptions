package getoptions

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var Logger = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

type argType int

const (
	argTypeProgname argType = iota
	argTypeCommand
	argTypeOption
	argTypeText
	argTypeTerminator // --
)

type CLIArg struct {
	Type     argType
	Name     string
	Args     []string
	Children []*CLIArg
}

func NewCLIArg(t argType, name string, args ...string) *CLIArg {
	arg := &CLIArg{
		Type:     t,
		Name:     name,
		Args:     []string{},
		Children: []*CLIArg{},
	}
	if len(args) > 0 {
		arg.Args = args
	}
	return arg
}

type CLITree struct {
	Type     argType
	Name     string
	Children []*CLITree
	Parent   *CLITree
}

func parseCLIArgs(tree *CLITree, args []string, mode Mode) *CLIArg {
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

	currentOpt := root
	for i, arg := range args {

		// handle terminator
		if arg == "--" {
			if len(args) > i+1 {
				for _, arg := range args[i+1:] {
					currentOpt.Children = append(currentOpt.Children, NewCLIArg(argTypeText, arg))
				}
			}
			break
		}

		// check for option
		cliArg, is := isOption(arg, mode)
		if is {
			currentOpt.Children = append(currentOpt.Children, cliArg...)
			continue
		}

		// handle command or text
	}
	return root
}

type GetOpt struct {
	cliTree *CLITree
}

type ModifyFn func(string)

func New() *GetOpt {
	gopt := &GetOpt{}
	gopt.cliTree = &CLITree{
		Type:     argTypeProgname,
		Name:     os.Args[0],
		Children: []*CLITree{},
	}
	return gopt
}

func (gopt *GetOpt) NewCommand(name string, description string) *GetOpt {
	cmd := &GetOpt{}
	tree := &CLITree{
		Type:     argTypeCommand,
		Name:     name,
		Children: []*CLITree{},
		Parent:   gopt.cliTree,
	}
	cmd.cliTree = tree
	gopt.cliTree.Children = append(gopt.cliTree.Children, tree)
	return cmd
}

func (gopt *GetOpt) String(name, def string, fns ...ModifyFn) *string {
	gopt.cliTree.Children = append(gopt.cliTree.Children, &CLITree{
		Type:     argTypeOption,
		Name:     name,
		Children: []*CLITree{},
		Parent:   gopt.cliTree,
	})
	return nil
}

func (gopt *GetOpt) StringVar(p *string, name, def string, fns ...ModifyFn) {
}

// Mode - Operation mode for short options
type Mode int

// Operation modes
const (
	Normal Mode = iota
	Bundling
	SingleDash
)

// 1: leading dashes
// 2: option
// 3: =arg
var isOptionRegex = regexp.MustCompile(`^(--?)([^=]+)(.*?)$`)

/*
isOption - Enhanced version of isOption, this one returns pairs of options and arguments
At this level we don't agregate results in case we have -- and then other options, basically we can parse one option at a time.
This makes the caller have to agregate multiple calls to the same option.
TODO: Here is where we should handle windows /option types.
*/
func isOption(s string, mode Mode) ([]*CLIArg, bool) {
	// Handle especial cases
	if s == "--" {
		return []*CLIArg{NewCLIArg(argTypeTerminator, "--")}, false
	} else if s == "-" {
		return []*CLIArg{NewCLIArg(argTypeOption, "-")}, true
	}

	match := isOptionRegex.FindStringSubmatch(s)
	if len(match) > 0 {
		// check long option
		if match[1] == "--" {
			opt := NewCLIArg(argTypeOption, match[2])
			args := strings.TrimPrefix(match[3], "=")
			if args != "" {
				// TODO: Here is where we could split on comma
				opt.Args = []string{args}
			}
			return []*CLIArg{opt}, true
		}
		// check short option
		switch mode {
		case Bundling:
			opts := []*CLIArg{}
			for _, option := range strings.Split(match[2], "") {
				opt := NewCLIArg(argTypeOption, option)
				opts = append(opts, opt)
			}
			if len(opts) > 0 {
				args := strings.TrimPrefix(match[3], "=")
				if args != "" {
					opts[len(opts)-1].Args = []string{args}
				}
			}
			return opts, true
		case SingleDash:
			opts := []*CLIArg{}
			for _, option := range []string{strings.Split(match[2], "")[0]} {
				opt := NewCLIArg(argTypeOption, option)
				opts = append(opts, opt)
			}
			if len(opts) > 0 {
				args := strings.Join(strings.Split(match[2], "")[1:], "") + match[3]
				opts[len(opts)-1].Args = []string{args}
			}
			return opts, true
		default:
			opt := NewCLIArg(argTypeOption, match[2])
			args := strings.TrimPrefix(match[3], "=")
			if args != "" {
				opt.Args = []string{args}
			}
			return []*CLIArg{opt}, true
		}
	}
	return []*CLIArg{}, false
}
