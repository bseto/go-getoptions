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

func TestBuildTree(t *testing.T) {
	buf := setupLogging()
	opt := New()
	opt.String("opt1", "")
	cmd := opt.NewCommand("cmd", "")
	cmd.String("opt2", "")
	tree := opt.cliTree
	expectedTree := &CLITree{
		Type:     argTypeProgname,
		Name:     os.Args[0],
		Children: []*CLITree{},
	}
	expectedTreeCmd := &CLITree{
		Type:     argTypeCommand,
		Name:     "cmd",
		Parent:   expectedTree,
		Children: []*CLITree{},
	}
	expectedTreeCmd.Children = append(expectedTreeCmd.Children, &CLITree{
		Type:     argTypeOption,
		Name:     "opt2",
		Parent:   expectedTreeCmd,
		Children: []*CLITree{},
	})
	expectedTree.Children = append(expectedTree.Children, []*CLITree{
		{
			Type:     argTypeOption,
			Name:     "opt1",
			Parent:   expectedTree,
			Children: []*CLITree{},
		},
		expectedTreeCmd,
	}...)

	if !reflect.DeepEqual(expectedTree, tree) {
		t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(expectedTree), spew.Sdump(tree))
	}
	t.Cleanup(func() { t.Log(buf.String()) })
}
