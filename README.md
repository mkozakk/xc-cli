# xc

A tiny clipboard manager for Linux. It pipes your command output directly to the system clipboard, with option of prepending command which produced that output.

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

The installer script does everything for you. It downloads the right binary for your architecture, installs it, and patches your shell config:

```bash
curl -fsSL https://raw.githubusercontent.com/mkozakk/xc-cli/main/scripts/install.sh | bash
```

The script will automatically detect your architecture, download the latest binary from github verify checksum, and set up the shell hook for bash or zsh. Just restart your terminal and you're good to go.

If you prefer to build from source:
```bash
git clone https://github.com/mkozakk/xc-cli
cd xc-cli && go build -o xc . && install -m755 xc ~/.local/bin/xc
eval "$(./xc init bash)"  # or zsh
```

## how it works

When you run the installer, it patches your shell rc file to install a hook that runs before and after every command you execute. The hook captures two things: the command text itself (like `git log --oneline`) and everything it outputs to stdout and stderr.

The shell hook saves the original file descriptors, redirects stdout and stderr through the `tee` command to capture them, and then restores the original descriptors after the command finishes. Your shell environment sees no difference.

```bash
exec 3>&1 4>&2                    # save originals
exec 1> >(tee -a capture) 2>&1   # redirect through tee
# command is ran here
exec 1>&3 2>&4                    # restore originals
```

The captured command and output are written to a temporary session file at `${XDG_RUNTIME_DIR:-/tmp}/xc_$(id -u)_$$.tmp`. The `$$` in the filename is the shell's own PID, so each shell instance gets a completely isolated session file. Running `xc -l` in one terminal won't break with another.

The file uses a simple line-based format with two markers:

```
CMD:git log --oneline -20
OUTPUT:
abc1234 fix login bug
def5678 add dark mode
...
```

When you run `xc -l` or `xc -lc`, the Go binary reads this file, extracts the command and output fields, and passes the result to the clipboard tool. The file is created with `0600` permissions so only your user can read it.

The two shell hooks use different mechanisms to hook commands. Bash uses a `DEBUG` trap (fires before each command) combined with `PROMPT_COMMAND` (fires after). Zsh has first-class `preexec` and `precmd` hook arrays for exactly this purpose.

Both hooks check whether the command being run starts with `xc` and skip capture in that case. This prevents `xc -l` itself from overwriting the session file it is about to read.

When you add `-c`, `xc` reads the session file to find the command name, then prepends `$ <command>` to the clipboard text either from piped stdin or from the stored output.

The tool detects which clipboard command is available on your system. On Wayland, it prefers `wl-copy`. On X11, it tries `xclip`, then `xsel`.

## requirements

- Linux with Wayland or X11
- At least one clipboard tool installed: `wl-copy`, `xclip` or `xsel`
- Bash 4.0+ or Zsh 5.0+
- Go 1.21+ for building if building from source

## license

MIT
