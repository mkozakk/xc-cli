package session

import (
	"os"
	"testing"
	"path/filepath"
)

func TestRead(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "session.tmp")

	sess := &Session{
		Command: "echo 'hello world'",
		Output:  "hello world\n",
	}

	if err := Write(testFile, sess); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	read, err := Read(testFile)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if read.Command != sess.Command {
		t.Errorf("Command mismatch: got %q, want %q", read.Command, sess.Command)
	}

	if read.Output != sess.Output {
		t.Errorf("Output mismatch: got %q, want %q", read.Output, sess.Output)
	}
}

func TestReadMissing(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "nonexistent.tmp")

	_, err := Read(testFile)
	if err != ErrNoSession {
		t.Errorf("Expected ErrNoSession, got %v", err)
	}
}

func TestFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "session.tmp")

	sess := &Session{Command: "test", Output: ""}
	if err := Write(testFile, sess); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	stat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if stat.Mode()&0o777 != 0o600 {
		t.Errorf("Wrong file mode: got %o, want 0o600", stat.Mode()&0o777)
	}
}
