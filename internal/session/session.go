package session

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Session struct {
	Command string
	Output  string
}

var ErrNoSession = errors.New("no session file found; is the shell hook installed? Run: eval \"$(xc init bash)\"")

func FilePath(user, shellPID string) string {
	return fmt.Sprintf("/tmp/xc_%s_%s.tmp", user, shellPID)
}

func CurrentFilePath() string {
	user := os.Getenv("USER")
	shellPID := fmt.Sprint(os.Getppid())
	return FilePath(user, shellPID)
}

func Read(path string) (*Session, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoSession
		}
		return nil, fmt.Errorf("read session file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read session file: %w", err)
	}

	content := string(data)
	if content == "" {
		return nil, fmt.Errorf("session file empty")
	}

	lines := strings.SplitN(content, "\n", 3)
	if len(lines) < 1 {
		return nil, fmt.Errorf("session file malformed")
	}

	cmdLine := lines[0]
	if !strings.HasPrefix(cmdLine, "CMD:") {
		return nil, fmt.Errorf("expected CMD: marker, got: %s", cmdLine)
	}
	cmd := strings.TrimSpace(strings.TrimPrefix(cmdLine, "CMD:"))

	if len(lines) < 2 {
		return &Session{Command: cmd, Output: ""}, nil
	}

	outputMarker := lines[1]
	if !strings.HasPrefix(outputMarker, "OUTPUT:") {
		return nil, fmt.Errorf("expected OUTPUT: marker, got: %s", outputMarker)
	}

	var output string
	if len(lines) > 2 {
		output = lines[2]
	}

	return &Session{
		Command: cmd,
		Output:  output,
	}, nil
}

func Write(path string, s *Session) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open session file: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "CMD:%s\n", s.Command)
	fmt.Fprint(file, "OUTPUT:\n")
	fmt.Fprint(file, s.Output)

	return nil
}
