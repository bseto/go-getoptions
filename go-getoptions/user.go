package getoptions

import (
	"context"
	"io/ioutil"
	"log"
	"os"
)

var Logger = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

type GetOpt struct {
	programTree *programTree
}

// Mode - Operation mode for short options
type Mode int

// Operation modes
const (
	Normal Mode = iota
	Bundling
	SingleDash
)

// CommandFn - Function signature for commands
type CommandFn func(context.Context, *GetOpt, []string) error

type ModifyFn func(string)

func New() *GetOpt {
	gopt := &GetOpt{}
	gopt.programTree = &programTree{
		Type:     argTypeProgname,
		Name:     os.Args[0],
		Children: []*programTree{},
		Level:    0,
	}
	return gopt
}

func (gopt *GetOpt) NewCommand(name string, description string) *GetOpt {
	cmd := &GetOpt{}
	command := &programTree{
		Type:     argTypeCommand,
		Name:     name,
		Children: []*programTree{},
		Parent:   gopt.programTree,
		Level:    gopt.programTree.Level + 1,
	}
	for _, child := range gopt.programTree.Children {
		// The option parent doesn't match properly here.
		// I should in a way create a copy of the option but I still want a pointer to the data.

		if child.Type == argTypeOption {
			c := child.Copy() // copy that maintains a pointer to the underlying data
			c.SetParent(command)
			command.Children = append(command.Children, c)
		}
	}
	cmd.programTree = command
	gopt.programTree.Children = append(gopt.programTree.Children, command)
	return cmd
}

func (gopt *GetOpt) String(name, def string, fns ...ModifyFn) *string {
	n := newCLIArg(gopt.programTree, argTypeOption, name)
	gopt.programTree.Children = append(gopt.programTree.Children, n)
	return nil
}

func (gopt *GetOpt) StringVar(p *string, name, def string, fns ...ModifyFn) {
}
