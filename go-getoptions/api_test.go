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

func TestCLITree(t *testing.T) {
	buf := setupLogging()
	opt := New()
	opt.String("opt1", "")
	cmd := opt.NewCommand("cmd1", "")
	cmd.String("cmd1opt1", "")
	cmd2 := opt.NewCommand("cmd2", "")
	cmd2.String("cmd2opt1", "")
	tree := opt.cliTree
	expectedTree := &CLITree{
		Type:     argTypeProgname,
		Name:     os.Args[0],
		Children: []*CLITree{},
	}
	expectedTreeCmd1 := &CLITree{
		Type:     argTypeCommand,
		Name:     "cmd1",
		Parent:   expectedTree,
		Children: []*CLITree{},
	}
	expectedTreeCmd1.Children = append(expectedTreeCmd1.Children, &CLITree{
		Type:     argTypeOption,
		Name:     "cmd1opt1",
		Parent:   expectedTreeCmd1,
		Children: []*CLITree{},
	})
	expectedTreeCmd2 := &CLITree{
		Type:     argTypeCommand,
		Name:     "cmd2",
		Parent:   expectedTree,
		Children: []*CLITree{},
	}
	expectedTreeCmd2.Children = append(expectedTreeCmd2.Children, &CLITree{
		Type:     argTypeOption,
		Name:     "cmd2opt1",
		Parent:   expectedTreeCmd2,
		Children: []*CLITree{},
	})
	expectedTree.Children = append(expectedTree.Children, []*CLITree{
		{
			Type:     argTypeOption,
			Name:     "opt1",
			Parent:   expectedTree,
			Children: []*CLITree{},
		},
		expectedTreeCmd1,
		expectedTreeCmd2,
	}...)

	if !reflect.DeepEqual(expectedTree, tree) {
		t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(expectedTree), spew.Sdump(tree))
	}
	t.Cleanup(func() { t.Log(buf.String()) })
}
