package getoptions

import (
	"context"
	"io/ioutil"
	"log"
	"os"
)


var Logger = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

type GetOpt struct {
	programTree *ProgramTree
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
	gopt.programTree = &ProgramTree{
		Type:     argTypeProgname,
		Name:     os.Args[0],
		Children: []*ProgramTree{},
	}
	return gopt
}

func (gopt *GetOpt) NewCommand(name string, description string) *GetOpt {
	cmd := &GetOpt{}
	tree := &ProgramTree{
		Type:     argTypeCommand,
		Name:     name,
		Children: []*ProgramTree{},
		Parent:   gopt.programTree,
	}
	for _, child := range gopt.programTree.Children {
		if child.Type == argTypeOption {
			tree.Children = append(tree.Children, child)
		}
	}
	cmd.programTree = tree
	gopt.programTree.Children = append(gopt.programTree.Children, tree)
	return cmd
}

func (gopt *GetOpt) String(name, def string, fns ...ModifyFn) *string {
	gopt.programTree.Children = append(gopt.programTree.Children, &ProgramTree{
		Type:     argTypeOption,
		Name:     name,
		Children: []*ProgramTree{},
		Parent:   gopt.programTree,
	})
	return nil
}

func (gopt *GetOpt) StringVar(p *string, name, def string, fns ...ModifyFn) {
}
