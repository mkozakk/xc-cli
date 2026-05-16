# xc

A tiny clipboard manager for Linux. It pipes your command output directly to the system clipboard, with option of prepending which command produced that output.

## usage

The tool supports four modes depending on what you want to capture:

```bash
command | xc              # copy just the output to clipboard
command | xc -c           # copy "$ command" + output together
xc -l                     # copy the text of the last command you ran
xc -lc                    # copy the last command + its full output
```

**Examples:**
```bash
# Copy git log to share in a PR comment
git log --oneline -20 | xc

# Copy kubectl output with the exact command that ran it
docker images | xc -c

# Quickly copy your previous command
xc -lc
```

## install

The installer script does everything for you — downloads the right binary for your architecture, installs it, and patches your shell config:

```bash
curl -fsSL https://raw.githubusercontent.com/mkozakk/xc/main/scripts/install.sh | bash
```

The script will automatically detect your architecture, download the latest binary from github verify hecksum, and set up the shell hook for bash or zsh. Just restart your terminal and you're ready to go.

If you prefer to build from source:
```bash
git clone https://github.com/YOUR_GITHUB_USER/xc
cd xc && go build -o xc . && install -m755 xc ~/.local/bin/xc
eval "$(./xc init bash)"  # or zsh
```

## how it works

When you run the installer, it patches your shell rc file to install a hook that runs before and after every command you execute. The hook captures two things: the command text itself (like `git log --oneline`) and everything it outputs to stdout and stderr.

**Output capture mechanism:** The shell hook saves the original file descriptors, redirects stdout and stderr through the `tee` command to capture them, and then restores the original descriptors after the command finishes. This happens transparently - your interactive shell environment sees no difference.

```bash
exec 3>&1 4>&2                    # save originals
exec 1> >(tee -a capture) 2>&1   # redirect through tee
# command is ran here
exec 1>&3 2>&4                    # restore originals
```

The captured command and output are written to a temporary session file at `/tmp/xc_${USER}_$$.tmp`. When you run `xc -l` or `xc -lc`, those flags just read this file and copy the content to your clipboard.

**Clipboard backends:** The tool auto-detects which clipboard command is available on your system. On Wayland, it prefers `wl-copy`. On X11, it tries `xclip` first.

## requirements

- Linux with Wayland or X11 display server
- At least one clipboard tool installed: `wl-copy`, `xclip` or `xsel`
- Bash 4.0+ or Zsh 5.0+
- Go 1.21+ for building if building from source

## license

MIT
