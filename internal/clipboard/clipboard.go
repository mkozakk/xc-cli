package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Backend int

const (
	BackendWlCopy Backend = iota
	BackendXclip
	BackendXsel
)

func (b Backend) String() string {
	switch b {
	case BackendWlCopy:
		return "wl-copy"
	case BackendXclip:
		return "xclip"
	case BackendXsel:
		return "xsel"
	default:
		return "unknown"
	}
}

func Detect() (Backend, error) {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if _, err := exec.LookPath("wl-copy"); err == nil {
			return BackendWlCopy, nil
		}
	}

	if _, err := exec.LookPath("wl-copy"); err == nil {
		return BackendWlCopy, nil
	}

	if _, err := exec.LookPath("xclip"); err == nil {
		return BackendXclip, nil
	}

	if _, err := exec.LookPath("xsel"); err == nil {
		return BackendXsel, nil
	}

	return 0, fmt.Errorf("no clipboard tool found; install wl-clipboard, xclip, or xsel")
}

func Write(text string) error {
	backend, err := Detect()
	if err != nil {
		return err
	}
	return WriteWithBackend(text, backend)
}

func WriteWithBackend(text string, b Backend) error {
	var cmd *exec.Cmd

	switch b {
	case BackendWlCopy:
		cmd = exec.Command("wl-copy")

	case BackendXclip:
		cmd = exec.Command("xclip", "-selection", "clipboard")

	case BackendXsel:
		cmd = exec.Command("xsel", "--clipboard", "--input")

	default:
		return fmt.Errorf("unknown backend: %v", b)
	}

	cmd.Stdin = strings.NewReader(text)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard write failed (%s): %w", b.String(), err)
	}

	return nil
}
