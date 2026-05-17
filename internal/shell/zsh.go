package shell

const zshHookTemplate = `
__xc_session_file() {
    local dir="${XDG_RUNTIME_DIR:-/tmp}"
    echo "${dir}/xc_$(id -u)_$$.tmp"
}

__xc_cmd=""
__xc_outfile=""
__xc_capturing=0

__xc_preexec() {
    local cmd="$1"
    [[ "$cmd" == xc* ]] && return
    [[ "$cmd" == __xc_* ]] && return

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
