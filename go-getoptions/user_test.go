package getoptions

import (
	"bytes"
	"io/ioutil"
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

	sub2cmd1 := cmd1.NewCommand("sub2cmd1", "")
	sub2cmd1.String("-", "")
	return opt
}

func SpewToFile(t *testing.T, e interface{}, label string) string {
	f, err := ioutil.TempFile("/tmp/", "spew-")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, _ = f.WriteString(label + "\n")
	spew.Fdump(f, e)
	return f.Name()
}

func TestTrees(t *testing.T) {
	buf := setupLogging()

	t.Run("programTree", func(t *testing.T) {
		tree := setupOpt().programTree
		root := &programTree{
			Type:          argTypeProgname,
			Name:          os.Args[0],
			ChildCommands: map[string]*programTree{},
			ChildOptions:  map[string]*option{},
		}
		rootopt1 := newCLIOption(root, "rootopt1")
		cmd1 := &programTree{
			Type:          argTypeCommand,
			Name:          "cmd1",
			Parent:        root,
			ChildCommands: map[string]*programTree{},
			ChildOptions:  map[string]*option{},
			Level:         1,
		}
		cmd1opt1 := newCLIOption(cmd1, "cmd1opt1")
		sub1cmd1 := &programTree{
			Type:          argTypeCommand,
			Name:          "sub1cmd1",
			Parent:        cmd1,
			ChildCommands: map[string]*programTree{},
			ChildOptions:  map[string]*option{},
			Level:         2,
		}
		sub1cmd1opt1 := newCLIOption(sub1cmd1, "sub1cmd1opt1")
		sub2cmd1 := &programTree{
			Type:          argTypeCommand,
			Name:          "sub2cmd1",
			Parent:        cmd1,
			ChildCommands: map[string]*programTree{},
			ChildOptions:  map[string]*option{},
			Level:         2,
		}
		sub2cmd1opt1 := newCLIOption(sub2cmd1, "-")
		cmd2 := &programTree{
			Type:          argTypeCommand,
			Name:          "cmd2",
			Parent:        root,
			ChildCommands: map[string]*programTree{},
			ChildOptions:  map[string]*option{},
			Level:         1,
		}
		cmd2opt1 := newCLIOption(cmd2, "cmd2opt1")

		root.ChildOptions["rootopt1"] = rootopt1
		root.ChildCommands["cmd1"] = cmd1
		root.ChildCommands["cmd2"] = cmd2

		rootopt1Copycmd1 := rootopt1.Copy().SetParent(cmd1)
		cmd1.ChildOptions["rootopt1"] = rootopt1Copycmd1
		cmd1.ChildOptions["cmd1opt1"] = cmd1opt1
		cmd1.ChildCommands["sub1cmd1"] = sub1cmd1
		cmd1.ChildCommands["sub2cmd1"] = sub2cmd1

		rootopt1Copycmd2 := rootopt1.Copy().SetParent(cmd2)
		cmd2.ChildOptions["rootopt1"] = rootopt1Copycmd2
		cmd2.ChildOptions["cmd2opt1"] = cmd2opt1

		rootopt1Copysub1cmd1 := rootopt1.Copy().SetParent(sub1cmd1)
		cmd1opt1Copysub1cmd1 := cmd1opt1.Copy().SetParent(sub1cmd1)

		sub1cmd1.ChildOptions["rootopt1"] = rootopt1Copysub1cmd1
		sub1cmd1.ChildOptions["cmd1opt1"] = cmd1opt1Copysub1cmd1
		sub1cmd1.ChildOptions["sub1cmd1opt1"] = sub1cmd1opt1

		rootopt1Copysub2cmd1 := rootopt1.Copy().SetParent(sub2cmd1)
		cmd1opt1Copysub2cmd1 := cmd1opt1.Copy().SetParent(sub2cmd1)
		sub2cmd1.ChildOptions["rootopt1"] = rootopt1Copysub2cmd1
		sub2cmd1.ChildOptions["cmd1opt1"] = cmd1opt1Copysub2cmd1
		sub2cmd1.ChildOptions["-"] = sub2cmd1opt1

		if !reflect.DeepEqual(root, tree) {
			t.Errorf("expected tree: %s, got: %s\n", SpewToFile(t, root, "expected"), SpewToFile(t, tree, "got"))
			t.Fatalf("expected tree: \n%s\n got: \n%s\n", root.Str(), tree.Str())
		}

		n, err := getNode(tree)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(root, n) {
			t.Errorf("expected tree: %s, got: %s\n", SpewToFile(t, root, "expected"), SpewToFile(t, tree, "got"))
			t.Fatalf("expected tree: \n%s\n got: \n%s\n", root.Str(), tree.Str())
		}

		n, err = getNode(tree, []string{}...)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(root, n) {
			// t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
			t.Errorf("expected tree: \n%s\n got: \n%s\n", root.Str(), tree.Str())
		}

		n, err = getNode(tree, "cmd1")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(cmd1, n) {
			// t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
			t.Errorf("expected tree: \n%s\n got: \n%s\n", root.Str(), tree.Str())
		}

		n, err = getNode(tree, "cmd1", "sub1cmd1")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(sub1cmd1, n) {
			// t.Errorf("expected tree: %s\n got: %s\n", spew.Sdump(root), spew.Sdump(tree))
			t.Errorf("expected tree: \n%s\n got: \n%s\n", root.Str(), tree.Str())
		}

		n, err = getNode(tree, "cmd2")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !reflect.DeepEqual(cmd2, n) {
			t.Errorf("expected tree: %s, got: %s\n", SpewToFile(t, cmd2, "expected"), SpewToFile(t, n, "got"))
			t.Errorf("expected tree: \n%s\n got: \n%s\n", cmd2.Str(), n.Str())
		}

	})

	t.Cleanup(func() { t.Log(buf.String()) })
}
