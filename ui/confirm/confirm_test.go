package confirm

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunYes(t *testing.T) {
	ok, err := run("Force push?", strings.NewReader("y"), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !ok {
		t.Fatal("expected yes")
	}
}

func TestRunNo(t *testing.T) {
	ok, err := run("Force push?", strings.NewReader("n"), bytes.NewBuffer(nil))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if ok {
		t.Fatal("expected no")
	}
}
