package getoptions

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// User facing tree construction tests.

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

func setupOpt() *GetOpt {
	opt := New()
	opt.String("rootopt1", "")

	cmd1 := opt.NewCommand("cmd1", "")
	cmd1.String("cmd1opt1", "")
	cmd2 := opt.NewCommand("cmd2", "")
	cmd2.String("cmd2opt1", "")

	sub1cmd1 := cmd1.NewCommand("sub1cmd1", "")
	sub1cmd1.String("sub1cmd1opt1", "")
	return opt
}

func TestTrees(t *testing.T) {
	buf := setupLogging()

	t.Run("programTree", func(t *testing.T) {
		tree := setupOpt().programTree
		root := &programTree{
			Type:     argTypeProgname,
			Name:     os.Args[0],
			Children: []*programTree{},
		}
		rootopt1 := &programTree{
			Type:     argTypeOption,
			Name:     "rootopt1",
			Parent:   root,
			Children: []*programTree{},
		}
		cmd1 := &programTree{
			Type:     argTypeCommand,
			Name:     "cmd1",
			Parent:   root,
			Children: []*programTree{},
			Level:    1,
		}
		cmd1opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "cmd1opt1",
			Parent:   cmd1,
			Children: []*programTree{},
		}
		sub1cmd1 := &programTree{
			Type:     argTypeCommand,
			Name:     "sub1cmd1",
			Parent:   cmd1,
			Children: []*programTree{},
			Level:    2,
		}
		sub1cmd1opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "sub1cmd1opt1",
			Parent:   sub1cmd1,
			Children: []*programTree{},
		}
		cmd2 := &programTree{
			Type:     argTypeCommand,
			Name:     "cmd2",
			Parent:   root,
			Children: []*programTree{},
			Level:    1,
		}
		cmd2opt1 := &programTree{
			Type:     argTypeOption,
			Name:     "cmd2opt1",
			Parent:   cmd2,
			Children: []*programTree{},
		}
		root.Children = append(root.Children, []*programTree{
			rootopt1,
			cmd1,
			cmd2,
		}...)
		cmd1.Children = append(cmd1.Children, []*programTree{
			rootopt1,
			cmd1opt1,
			sub1cmd1,
		}...)
		sub1cmd1.Children = append(sub1cmd1.Children, []*programTree{
			rootopt1,
			cmd1opt1,
			sub1cmd1opt1,
		}...)
		cmd2.Children = append(cmd2.Children, []*programTree{
			rootopt1,
			cmd2opt1,
		}...)

		if !reflect.DeepEqual(root, tree) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err := getNode(tree)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(root, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, []string{}...)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(root, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, "cmd1")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(cmd1, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, "cmd1", "sub1cmd1")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(sub1cmd1, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

		n, err = getNode(tree, "cmd2")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(cmd2, n) {
			t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
		}

	})

	t.Cleanup(func() { t.Log(buf.String()) })
}
