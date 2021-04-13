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

type CLITree struct {
	Type     argType
	Name     string
	Children []*CLITree
	Parent   *CLITree
}

func parseCLIArgs(tree *CLITree, args []string) *CLIArg {
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

	root := &CLIArg{
		Type:     argTypeProgname,
		Name:     os.Args[0],
		Args:     args, // Copy of the original args
		Children: []*CLIArg{},
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

type optionPair struct {
	Option string
	// We allow multiple args in case of splitting on comma.
	Args []string
}

/*
isOption - Enhanced version of isOption, this one returns pairs of options and arguments
At this level we don't agregate results in case we have -- and then other options, basically we can parse one option at a time.
This makes the caller have to agregate multiple calls to the same option.
TODO: Here is where we should handle windows /option types.
*/
func isOption(s string, mode Mode) ([]optionPair, bool) {
	// Handle especial cases
	if s == "--" {
		return []optionPair{{Option: "--"}}, false
	} else if s == "-" {
		return []optionPair{{Option: "-"}}, true
	}

	match := isOptionRegex.FindStringSubmatch(s)
	if len(match) > 0 {
		// check long option
		if match[1] == "--" {
			opt := optionPair{}
			opt.Option = match[2]
			args := strings.TrimPrefix(match[3], "=")
			if args != "" {
				// TODO: Here is where we could split on comma
				opt.Args = []string{args}
			}
			return []optionPair{opt}, true
		}
		// check short option
		switch mode {
		case Bundling:
			opts := []optionPair{}
			for _, option := range strings.Split(match[2], "") {
				opt := optionPair{}
				opt.Option = option
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
			opts := []optionPair{}
			for _, option := range []string{strings.Split(match[2], "")[0]} {
				opt := optionPair{}
				opt.Option = option
				opts = append(opts, opt)
			}
			if len(opts) > 0 {
				args := strings.Join(strings.Split(match[2], "")[1:], "") + match[3]
				opts[len(opts)-1].Args = []string{args}
			}
			return opts, true
		default:
			opt := optionPair{}
			opt.Option = match[2]
			args := strings.TrimPrefix(match[3], "=")
			if args != "" {
				opt.Args = []string{args}
			}
			return []optionPair{opt}, true
		}
	}
	return []optionPair{}, false
}
