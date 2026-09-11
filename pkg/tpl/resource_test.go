package tpl

import (
	"bytes"
	"testing"
)

func TestGitIgnoreResource(t *testing.T) {
	if bytes.HasPrefix(GitIgnore, []byte("FROM ")) {
		t.Fatal(".gitignore resource contains the Dockerfile template")
	}
	if !bytes.Contains(GitIgnore, []byte("*.test")) {
		t.Fatal(".gitignore resource does not contain expected Go patterns")
	}
}
