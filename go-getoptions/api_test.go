package getoptions

import (
	"bytes"
	"testing"
)

func setupLogging() *bytes.Buffer {
	s := ""
	buf := bytes.NewBufferString(s)
	Logger.SetOutput(buf)
	return buf
}

func TestBuildTree(t *testing.T) {
	buf := setupLogging()
	opt := New()
	t.Cleanup(func() { t.Log(buf.String()) })
}
