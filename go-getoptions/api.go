package getoptions

import (
	"io/ioutil"
	"log"
	"os"
)

var Logger = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

type argType int

const (
	argTypeProgname argType = iota
	argTypeText
	argTypeOption
	argTypeCommand
	argTypeTerminator // --
)

type CLIArg struct {
	Type     argType
	Name     string
	Args     []string
	Children []CLIArg
}

func parseCLIArgs(args []string) *CLIArg {
	// Design: This function could return an array or CLIargs as a parse result
	// or I could do one level up and have a root CLIarg type with the name of
	// the program.  Having the root level might be helpful with help generation.

	// The current implementation expects os.Args[1:] as an argument so this
	// can't expect to receive the os.Args[0] as the first argument.

	root := &CLIArg{
		Type:     argTypeProgname,
		Name:     os.Args[0],
		Args:     args, // Copy of the original args
		Children: []CLIArg{},
	}
	return root
}

type GetOpt struct {
}

type ModifyFn func(string)

func New() *GetOpt {
	gopt := &GetOpt{}
	return gopt
}

func (gopt *GetOpt) NewCommand(name string, description string) *GetOpt {
	cmd := &GetOpt{}
	return cmd
}

func (gopt *GetOpt) String(name, def string, fns ...ModifyFn) *string {
	return nil
}

func (gopt *GetOpt) StringVar(p *string, name, def string, fns ...ModifyFn) {
}
