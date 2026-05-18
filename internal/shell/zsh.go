package shell

const zshHookTemplate = `
__xc_session_file() {
    local dir="${XDG_RUNTIME_DIR:-/tmp}"
    echo "${dir}/xc_$(id -u)_$$.tmp"
}

__xc_cmd=""
__xc_outfile=""
__xc_capturing=0

__xc_skip_cmd() {
    [[ ! -t 0 || ! -t 1 ]] && return 0

    local cmd="$1"
    local base_cmd="${cmd%% *}"
    base_cmd="${base_cmd##*/}"

    local skip_cmds=" vi vim nvim nano emacs micro joe kak helix hx vis \
less more most view rview \
top htop btop bpytop glances iotop iftop bmon atop nmon vtop gtop bashtop \
screen tmux mosh ssh telnet rlogin minicom picocom tio \
gdb lldb strace ltrace perf \
python python3 ipython irb node deno bun lua luajit ghci julia \
ranger lf nnn yazi mc tig gitui \
lazygit lazydocker k9s fzf watch ncdu duf \
mysql psql sqlite3 redis-cli mongosh mongo mycli pgcli litecli \
sudo sudoedit su doas pkexec passwd htpasswd chpasswd openssl \
claude calude aichat aider mods gh amazon-q copilot aws awslocal gptme open-interpreter interpreter ollama gum clear"

    [[ " $skip_cmds " == *" $base_cmd "* ]] && return 0
    
    local custom_skip_cmds="${XC_SKIP_CMDS:-}"
    [[ -n "$custom_skip_cmds" && " $custom_skip_cmds " == *" $base_cmd "* ]] && return 0

    return 1
}

__xc_preexec() {
    local cmd="$1"
    [[ "$cmd" == xc* ]] && return
    [[ "$cmd" == __xc_* ]] && return
    __xc_skip_cmd "$cmd" && return

    local display_cmd
    display_cmd="$(print -r -- "$cmd" | sed 's/[[:space:]]*|[[:space:]]*xc\b.*$//')"
    __xc_cmd="${display_cmd:-$cmd}"

    local session_file tmp_session
    session_file="$(__xc_session_file)"
    tmp_session="$(mktemp "${XDG_RUNTIME_DIR:-/tmp}/xc_sess_XXXXXX")"
    chmod 600 "$tmp_session"
    {
        print "CMD:${__xc_cmd}"
        print "OUTPUT:"
    } > "$tmp_session"
    mv -f "$tmp_session" "$session_file"

    __xc_outfile="$(mktemp "${XDG_RUNTIME_DIR:-/tmp}/xc_out_XXXXXX")"
    __xc_capturing=1

    exec 3>&1 4>&2
    exec 1> >(tee -a "$__xc_outfile") 2>&1
}

__xc_precmd() {
    if [[ "$__xc_capturing" -eq 1 ]]; then
        __xc_capturing=0
        exec 1>&3 2>&4
        exec 3>&- 4>&-

        local session_file tmp_session
        session_file="$(__xc_session_file)"
        tmp_session="$(mktemp "${XDG_RUNTIME_DIR:-/tmp}/xc_sess_XXXXXX")"
        chmod 600 "$tmp_session"
        {
            print "CMD:${__xc_cmd}"
            print "OUTPUT:"
            if [[ -f "$__xc_outfile" ]]; then
                sed $'s/\x1b\[[0-9;]*[mGKHFABCDsuhjJK]//g; s/\x1b][^\x07]*\x07//g; s/\r//g' "$__xc_outfile"
                rm -f "$__xc_outfile"
            fi
        } > "$tmp_session"
        mv -f "$tmp_session" "$session_file"
    fi
}

command_not_found_handler() {
    print "zsh: command not found: $1" >&2
    return 127
}

autoload -Uz add-zsh-hook
add-zsh-hook preexec __xc_preexec
precmd_functions=(__xc_precmd "${precmd_functions[@]}")
`

func ZshHook() string {
	return zshHookTemplate
}
