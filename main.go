package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mkozakk/xc-cli/internal/clipboard"
	"github.com/mkozakk/xc-cli/internal/session"
	"github.com/mkozakk/xc-cli/internal/shell"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			handleInit()
			return
		case "version":
			handleVersion()
			return
		case "-h", "--help", "help":
			printHelp()
			return
		}
	}

	flagContext := flag.Bool("c", false, "prefix output with '$ <command>'")
	flagLast := flag.Bool("l", false, "copy the last command from session")
	flag.Parse()

	stat, _ := os.Stdin.Stat()
	isPiped := (stat.Mode() & os.ModeCharDevice) == 0

	if !isPiped && !*flagLast {
		printHelp()
		os.Exit(0)
	}

	var err error
	switch {
	case isPiped && !*flagLast && !*flagContext:
		err = modeStandardPipe()
	case isPiped && !*flagLast && *flagContext:
		err = modeContextPipe()
	case !isPiped && *flagLast && !*flagContext:
		err = modeLastCommand()
	case !isPiped && *flagLast && *flagContext:
		err = modeLastCommandWithOutput()
	default:
		printHelp()
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "xc: %v\n", err)
		os.Exit(1)
	}
}

const maxStdinBytes = 50 * 1024 * 1024 // 50 MB

func modeStandardPipe() error {
	stdin, err := io.ReadAll(io.LimitReader(os.Stdin, maxStdinBytes))
	if err != nil {
		return fmt.Errorf("read stdin: %w", err)
	}
	os.Stdout.Write(stdin)
	return clipboard.Write(string(stdin))
}

func modeContextPipe() error {
	sess, err := session.Read(session.CurrentFilePath())
	if err != nil {
		return err
	}

	output, err := io.ReadAll(io.LimitReader(os.Stdin, maxStdinBytes))
	if err != nil {
		return fmt.Errorf("read stdin: %w", err)
	}

	os.Stdout.Write(output)
	text := fmt.Sprintf("$ %s\n%s", sess.Command, string(output))
	return clipboard.Write(text)
}

func modeLastCommand() error {
	sess, err := session.Read(session.CurrentFilePath())
	if err != nil {
		return err
	}
	return clipboard.Write(sess.Command)
}

func modeLastCommandWithOutput() error {
	sess, err := session.Read(session.CurrentFilePath())
	if err != nil {
		return err
	}

	text := fmt.Sprintf("$ %s\n%s", sess.Command, sess.Output)
	return clipboard.Write(text)
}

func handleInit() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: xc init [bash|zsh]\n")
		os.Exit(1)
	}

	shellType := os.Args[2]
	switch shellType {
	case "bash":
		fmt.Print(shell.BashHook())
	case "zsh":
		fmt.Print(shell.ZshHook())
	default:
		fmt.Fprintf(os.Stderr, "unknown shell: %s\n", shellType)
		os.Exit(1)
	}
}

func handleVersion() {
	fmt.Printf("xc version %s (commit %s, built %s)\n", version, commit, date)
}

func printHelp() {
	help := `xc — Smart Linux Clipboard Manager

Usage:
  command | xc              Copy output to clipboard
  command | xc -c           Copy "$ command\noutput" to clipboard
  xc -l                     Copy the last command text
  xc -l -c                  Copy the last command + its output
  xc init [bash|zsh]        Print shell hook to add to ~/.bashrc or ~/.zshrc
  xc version                Show version information
  xc --help                 Show this help

Examples:
  git log --oneline | xc
  kubectl get pods | xc -c
  xc -l
  xc -lc

For more information, visit: https://github.com/mkozakk/xc-cli
`
	fmt.Print(help)
}
